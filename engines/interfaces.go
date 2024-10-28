package engines

import (
    "sync";
    "os";
    "go/types";
    rbytree "github.com/emirpasic/gods/trees/redblacktree"
)

type Column struct{
    Type types.Type; // too deep?
    PrimaryIndex bool;
    DataSize int; // in bytes
    Data []interface{};
}

type Table struct {
    TableName string;
    Columns map[string]Column;
}

type Row struct {
    Columns map[string]interface{}; // types? schema? yes.
}

type WriteAheadLog interface {
    MakeTable() (Table, error) // a lot of mmaps? need to realloc?
    WriteTransaction([]Row) error // not really a transaction, I just can't think of a better name
    ReadLast(nrows int) ([]Row, error)
    Load() (Table, error) // load wal from disk
}
    

type Database interface {
    MakeTable(columns map[string]types.Type, // header: type
        tableName, dbDir, walPath string) (Engine, error)
    DeleteTable(tableName, dbDir string) error
    GetTable(tableName string) (Engine, error)
}


type Engine interface {
    // kinda want to say data is supposed to be consecutive, but that's limiting
    Create(data []Row) error // aka insert
    
    // note that updates are expected in batches 
    // read data from table in memory. Should think of better query criteria than "key equals stuff"
    Read(key_column string, keys []byte) ([]Row, error) 
    Update(key_column string, keys []byte, values []Row) error
    Delete(key_column string, keys []byte) error 
    
    // assuming the column is numerical
    ReadRange(Column string, lower_bound float64, upper_bound float64) ([]Row, error)
    
    // todo figure out usage
    // todo aggregation
    GetAvg(Column string) (float64, error)
    GetMode(Column string) (float64, error)

    // disk operations
    Flush() (int, error) 
    Load() (Table, error) 
}

type memTable struct {
    data rbytree.Tree // can't be bothered to do red-black trees right now. Can I just use GoDS?
    b int
}

type LsmTree struct {
    mutex *sync.Mutex
    dbDir string  // directory to store data
    wal *os.File // write ahead log file
    maxDiskTableIndex int // latest disk table
    diskTableNum int // total disk tables
    memTable *memTable // current mem table
    memTableThreshold int // flush after this threshold
    diskTableNumThreshold int // merge after this threshold
}

type Log struct {
    Mutex *sync.Mutex;
    Table Table
    DbDir string
    Wal *os.File // do I need wal for log-like stuff?
}

