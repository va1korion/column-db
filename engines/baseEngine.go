package engines

import "bytes"

type Column struct{
	Name string;
	Type string;
	Primary bool;
	Data bytes.Buffer;
}

type Table struct {
	Col Column;
	DiskAddress int;
}
	
type Engine interface {
    MakeTable() Table
    Insert() error
	Delete() error
	Update() error
	Read() *bytes.Buffer
	Save() error
}


