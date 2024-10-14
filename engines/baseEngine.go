package engines

import "bytes"

type Column struct{
	Name string;
	Type string;
	PrimaryIndex bool;
	Data bytes.Buffer;
}

type Table struct {
	Col Column;
	DiskAddress int;
	PrimaryIndex []int;
	SecondaryIndex map[string][]int;
}
	
type Engine interface {
    MakeTable() (Table, error)
    Insert() error
	Delete() error
	Update() error 
	Read() (*bytes.Buffer, error) // read data from table
	Save() (int, error) // save table as a file on disk
	Load() (Table, error) // load table from disk
}


type WriteAheadLogEngine interface {
	MakeTable() (Table, error) // a lot of mmaps? need to realloc?
	WriteTransaction() error // not really a transaction, I just can't think of a better name
	LoadWAL() (Table, error) // load wal from disk
}