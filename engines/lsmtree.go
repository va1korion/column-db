package engines

import (
	"bytes";
 	"sync";
	"fmt";
	"os"
)

var ErrNotImplemented = fmt.Errorf("not implemented error")

func (l LsmTree) MakeTable() (LsmTree, error) {
	wal, _ := os.OpenFile("/mnt/d/projects/github.com/va1korion/column-db/data/wal/wal", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return LsmTree{
		Mutex: &sync.Mutex{},
		DbDir: "/mnt/d/projects/github.com/va1korion/column-db/data/lsm", 
		Wal: wal,
		MaxDiskTableIndex: 4096, 
		DiskTableNum: 8, 
		MemTable: &memTable{
			data: Tree{},
			b: 0,
		}, 
		MemTableThreshold: 4096, 
		DiskTableNumThreshold: 65536, 
	}, ErrNotImplemented
}

func (l LsmTree) Insert() error {
	return ErrNotImplemented
}

func (l LsmTree) Delete() error {
	return ErrNotImplemented

}

func (l LsmTree) Update() error {
	return ErrNotImplemented

}

func (l LsmTree) Read() (*bytes.Buffer, error) {
	return nil, ErrNotImplemented

}

func (l LsmTree) Flush() (int, error) {
	return 0, ErrNotImplemented
}

