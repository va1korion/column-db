package engines

import (
	"bytes";
 	sync "sync";
	fmt "fmt"
)

var ErrNotImplemented = fmt.Errorf("not implemented error")


type LsmTreeEngine struct {
	mutex *sync.Mutex
	table Table
}

func (l LsmTreeEngine) MakeTable() (Table, error) {
	return Table{}, ErrNotImplemented
}

func (l LsmTreeEngine) Insert() error {
	return ErrNotImplemented
}

func (l LsmTreeEngine) Delete() error {
	return ErrNotImplemented

}

func (l LsmTreeEngine) Update() error {
	return ErrNotImplemented

}

func (l LsmTreeEngine) Read() (*bytes.Buffer, error) {
	return nil, ErrNotImplemented

}

func (l LsmTreeEngine) Save() (int, error) {
	return 0, ErrNotImplemented
}

