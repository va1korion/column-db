package engines

import (
    "bytes";
    "sync";
    "os";
    "go/types"
)
type Column struct{
    Type types.Type; // too deep?
    PrimaryIndex bool;
    DataSize int; // in bytes
    Data bytes.Buffer;
}

type Table struct {
    TableName string;
    Columns map[string]Column;
}

type Row struct {
    Columns map[string]bytes.Buffer; // types? schema?
    Types map[string]types.Type
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
    Read(key_column string, keys []bytes.Buffer) ([]Row, error) 
    Update(key_column string, keys []bytes.Buffer, values []Row) error
    Delete(key_column string, keys []bytes.Buffer) error 
    
    // assuming the column is numerical
    ReadRange(Column string, lower_bound float64, upper_bound float64) ([]Row, error)
    // todo figure out usage
    // todo aggregation
    Flush() (int, error) 
    Load() (Table, error) 
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

