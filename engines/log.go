package engines

import (
	"bytes";
)


func (l Log) MakeTable() (Table, error) {
	return Table{}, ErrNotImplemented
}

func (l Log) Insert() error {
	return ErrNotImplemented
}

func (l Log) Delete() error {
	return ErrNotImplemented
}

func (l Log) Update() error {
	return ErrNotImplemented

}

func (l Log) Read() (*bytes.Buffer, error) {
	return nil, ErrNotImplemented

}

func (l Log) Save() (int, error) {
	return 0, ErrNotImplemented
}

