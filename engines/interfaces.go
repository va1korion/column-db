package engines

import (
	"bytes"
	"os"
	"sync"
)

const walFileName = "wal.gob"

type IColumn interface {
	Grow(idx uint32)
	Apply([]byte)
	Value(idx uint32) (interface{}, bool)
	Contains(idx uint32) bool
	Index([]byte) error
	Snapshot(chunk []byte, dst *bytes.Buffer)
}

// Numeric represents a column that stores numbers.
type Numeric interface {
	IColumn
	LoadFloat64(uint32) (float64, bool)
	LoadInt64(uint32) (int64, bool)
	FilterFloat64([]byte, func(v float64) bool)
	FilterInt64([]byte, func(v int64) bool)
}

// Textual represents a column that stores strings.
type Textual interface {
	IColumn
	LoadString(uint32) (string, bool)
	FilterString([]byte, func(v string) bool)
}

type Column struct {
	Mutex        *sync.Mutex
	Type         ColumnType // too deep?
	PrimaryIndex bool
	DataSize     int // in bytes
	Data         []interface{}
	Avg          float64
	Count        int
	Sum          float64
	Min          float64
	Max 	     float64
}

type Table struct {
	mutex     *sync.Mutex
	TableName string
	Columns   map[string]*Column
	Schema    map[string]ColumnType
}

type Row struct {
	Schema  map[string]ColumnType
	Columns map[string]interface{} // types? schema? yes.
}

type WriteAheadLog interface {
	MakeTable() (Table, error)    // a lot of mmaps? need to realloc?
	WriteTransaction([]Row) error // not really a transaction, I just can't think of a better name
	ReadLast(nrows int) ([]Row, error)
	Load() (Table, error) // load wal from disk
}

type Database interface {
	MakeTable(columns map[string]ColumnType, tableName, dbDir, walPath string) (Table, error)
	DeleteTable(tableName, dbDir string) error
	GetTable(tableName string) (Table, error)
}

type Engine interface {
	// kinda want to say data is supposed to be consecutive, but that's limiting
	Create(data []Row) error // aka insert

	// note that updates are expected in batches
	// read data from table in memory. Should think of better query criteria than "key equals stuff"
	Read(key_column string, key interface{}) ([]Row, error)
	Update(key_column string, key interface{}, value Row) error
	Delete(key_column string, key interface{}) error

	// assuming the column is numerical
	ReadRange(Column string, lower_bound int, upper_bound int) ([]Row, error)

	// todo figure out usage
	// todo aggregation
	GetAvg(Column string) (float64, error)
	GetMode(Column string) (float64, error)
	Count(Column string, Value []byte) (int, error)

	// disk operations, todo later
	/*
	   Flush() (int, error)
	   Load() (Table, error)
	*/
}

// later
/*
type LsmTree struct {
    mutex *sync.Mutex
    memTable *memTable // current mem table
    dbDir string  // directory to store data
    wal *os.File // write ahead log file
    maxDiskTableIndex int // latest disk table
    diskTableNum int // total disk tables
    memTableThreshold int // flush after this threshold
    diskTableNumThreshold int // merge after this threshold
}
*/

type Log struct {
	Table Table
	DbDir string
	Wal   *os.File // do I need wal for log-like stuff?
}
