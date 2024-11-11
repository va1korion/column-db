package engines

import (
	"fmt"
	"sync"
)

var InternalError = fmt.Errorf("Some internal error")

type InMemory struct {
	tables map[string]Table
}

func (db InMemory) MakeTable(columns map[string]columnType, tableName, dbDir, walPath string) (Table, error) {

	cols := make(map[string]*Column)
	for key, value := range columns {
		var ds int
		switch value{
			case typeNumeric: ds = 8
			case typeTextual: ds = 256
			default : ds = 256
		}

		cols[key] = &Column{
			Mutex: &sync.Mutex{},
			Type: value,
			PrimaryIndex: false,
			DataSize: ds,
		}
	}

	table := Table{
		TableName: tableName,
		Columns: cols,
		Schema: columns,
	}
	return table, ErrNotImplemented
}


func (db InMemory) DeleteTable(tableName, dbDir string) error{
	delete(db.tables, tableName)
	return nil
}

func (db InMemory) GetTable(tableName string) (Table, error) {
	return db.tables[tableName], nil
}


func insertRow(row Row, t Table) error {
	for key, value := range row.Schema {
		t.Columns[key].Mutex.Lock()
		defer t.Columns[key].Mutex.Unlock()
		t.Columns[key].Data = append(t.Columns[key].Data, value) // basically delegerate shuffling pointers around to golang
	}
	return nil
}

func (t Table) Create(data []Row) error {
	// aka insert
	t.mutex.Lock()
	for row := range data {
		
		go insertRow(data[row], t) // probably need table-wide mutex
		
	}
	return nil
}


func (t Table) findIndices(key_column string, key interface{}) ([]uint32, error){
	col := t.Columns[key_column]
	idx := make([]uint32, 0)
	
	for k, v := range (*col).Data {
		if v == key {
			idx = append(idx, uint32(k))
		}
	}
	return idx, nil
}

func (t Table) Read(key_column string, key interface{}) ([]Row, error) {
	idx, _ := t.findIndices(key_column, key)
	
	rows := make([]Row, 0)
	for _, i := range idx {
		// todo add write logs
		// t.Mutex.Lock()
		// defer t.Mutex.Unlock()
		ans := make(map[string]*interface{})
		for k, v := range t.Columns {
			ans[k] = &v.Data[i]
		}

		rows = append(rows, Row{
			Columns: ans,
			Schema: t.Schema,
		},	
		)
	}

	return rows, nil
}

func (t Table) Update(key_column string, key interface{}, value Row) error {
	idx, _ := t.findIndices(key_column, key)
	for _, i := range idx {
		for k, v := range value.Columns {
			t.Columns[k].Data[i] = *v
		}
	}

	return nil
}

func (t Table) Delete(key_column string, key interface{}) error {
	idx, _ := t.findIndices(key_column, key)
	for _, i := range idx {
		for k := range t.Columns {
			// O(N) deletion is what I get for not doing a tree
			t.Columns[k].Data = append(t.Columns[k].Data[:i], t.Columns[k].Data[i+1:]...) 
		}
	}
	return nil
}

func (t Table) ReadRange(Column string,  lower_bound int, upper_bound int) ([]Row, error) {
	rows := make([]Row, 0)

	for k, _ := range t.Columns[Column].Data {
		if k >= lower_bound && k <= upper_bound {
			rows = append(rows, Row{
				Columns: map[string]*interface{}{
					Column: &t.Columns[Column].Data[k],
				},
				Schema: t.Schema,
			})
		}
	}
	return rows, nil
}


func (t Table) GetAvg(Column string) (float64, error) {
	return 0.0, nil
}

func (t Table) GetMode(Column string) (float64, error) {
	return 0.0, nil
}

func (t Table) Count(Column string, Value []byte) (int, error) {
	return 0, nil
}

