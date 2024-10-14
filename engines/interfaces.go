package engines

import (
	"bytes";
 	"sync";
	"os"
)
type Column struct{
	Type string;
	PrimaryIndex bool;
	DataSize int;
	Data bytes.Buffer;
}

type Table struct {
	Columns map[string]Column;
	PrimaryIndex []int;
}

type Row struct {
	Columns map[string]bytes.Buffer; // basically a JSON
}

type WriteAheadLog interface {
	MakeTable() (Table, error) // a lot of mmaps? need to realloc?
	WriteTransaction() error // not really a transaction, I just can't think of a better name
	Load() (Table, error) // load wal from disk
}
	
type Engine interface {
    MakeTable(headers []string, dbDir, walPath string) (Engine, error)
	// kinda want to say data is supposed to be consecutive, but that's limiting
    Insert(data []Row) error 
	// note that updates are expected in batches, 
	// I'm gonna die implementing it later
	// keys and values here are supposed to correspond to this SQL 
	// SELECT * FROM table WHERE columns[0] = values[0]
	Delete(columns []string, values []bytes.Buffer) error 
	Update(columns []int, values []bytes.Buffer, data []Row) error

	Read(columns []int, values []bytes.Buffer) (Row, error) // read data from table in memory
	Flush() (int, error) // save table as a file on disk
	Load() (Table, error) // load table from disk
	DeleteTable() error
}

type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

type memTable struct {
	data Tree // can't be bothered to do red-black trees right now. Can I just use GoDS?
	b int
}

type LsmTree struct {
	Mutex *sync.Mutex
	DbDir string  // directory to store data
	Wal *os.File // write ahead log file
	MaxDiskTableIndex int // latest disk table
	DiskTableNum int // total disk tables
	MemTable *memTable // current mem table
	MemTableThreshold int // flush after this threshold
	DiskTableNumThreshold int // merge after this threshold
}

type Log struct {
	Mutex *sync.Mutex;
	Table Table
	DbDir string
	Wal *os.File // do I need wal for log-like stuff?
}

