package engines

import (
	"fmt"
	"sync"
)

var InternalError = fmt.Errorf("some internal error")

type InMemory struct {
	Tables map[string]Table
}

func (db InMemory) MakeTable(columns map[string]ColumnType, tableName, dbDir, walPath string) (Table, error) {

	cols := make(map[string]*Column)
	for key, value := range columns {
		var ds int
		switch value {
		case TypeNumeric:
			ds = 8
		case TypeTextual:
			ds = 256
		default:
			ds = 256
		}

		cols[key] = &Column{
			Mutex:        &sync.Mutex{},
			Type:         value,
			PrimaryIndex: false,
			DataSize:     ds,
			Avg:          0.0,
			Count:        0,
			Sum:          0.0,
			Min:          0.0,
			Max:          0.0,
		}
	}

	table := Table{
		TableName: tableName,
		Columns:   cols,
		Schema:    columns,
	}
	DB.Tables[tableName] = table
	return table, ErrNotImplemented
}

func (db InMemory) DeleteTable(tableName, dbDir string) error {
	delete(db.Tables, tableName)
	return nil
}

func (db InMemory) GetTable(tableName string) (Table, error) {
	return db.Tables[tableName], nil
}

func (t Table) InsertRow(row Row) error {
	for key, value := range row.Columns {
		t.Columns[key].Mutex.Lock()
		defer t.Columns[key].Mutex.Unlock()
		t.Columns[key].Data = append(t.Columns[key].Data, value) // basically delegerate shuffling pointers around to golang
		fmt.Println(t.Columns[key].Data)
		t.Columns[key].Count++
		if t.Columns[key].Type == TypeNumeric {
			t.Columns[key].Sum += float64(value.(float64))
			t.Columns[key].Avg = t.Columns[key].Sum / float64(t.Columns[key].Count)
			
			if t.Columns[key].Min > float64(value.(float64)) {
				t.Columns[key].Min = float64(value.(float64))
			}
			if t.Columns[key].Max < float64(value.(float64)) {
				t.Columns[key].Max = float64(value.(float64))
			}
		}
	}
	return nil
}

func (t Table) Create(data []Row) error {
	// aka insert
	t.mutex.Lock()
	for row := range data {

		go t.InsertRow(data[row]) // probably need table-wide mutex

	}
	return nil
}

func (t Table) findIndices(key_column string, key interface{}) ([]uint32, error) {
	col := t.Columns[key_column]
	
	idx := make([]uint32, 1)
	for k, v := range col.Data {
		fmt.Println(key, v, k)
		if v == key {
			idx[0] = uint32(k)
			fmt.Println(idx)
			return idx, nil
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
		ans := make(map[string]interface{})
		for k, v := range t.Columns {
			ans[k] = &v.Data[i]
		}

		rows = append(rows, Row{
			Columns: ans,
			Schema:  t.Schema,
		},
		)
	}

	return rows, nil
}

func (t Table) Update(key_column string, key interface{}, value Row) error {
	idx, _ := t.findIndices(key_column, key)
	for _, i := range idx {
		for k, v := range value.Columns {
			t.Columns[k].Data[i] = v
		}
	}

	return nil
}

func (t Table) Delete(key_column string, key interface{}) error {
	idx, _ := t.findIndices(key_column, key)
	fmt.Println(idx)
	
	for _, i := range idx {
		fmt.Println("deleting", i)
		for k := range t.Columns {
			// O(N) deletion is what I get for not doing a tree
			value := t.Columns[k].Data[i]

			t.Columns[k].Data = append(t.Columns[k].Data[:i], t.Columns[k].Data[i+1:]...)
			t.Columns[k].Count--
			if t.Columns[k].Type == TypeNumeric {
				t.Columns[k].Sum -= float64(value.(float64))
				t.Columns[k].Avg = t.Columns[k].Sum / float64(t.Columns[k].Count)
			}
		}
	}
	return nil
}

func (t Table) ReadRange(Column string, lower_bound int, upper_bound int) ([]Row, error) {
	rows := make([]Row, 0)

	for k, _ := range t.Columns[Column].Data {
		if k >= lower_bound && k <= upper_bound {
			rows = append(rows, Row{
				Columns: map[string]interface{}{
					Column: &t.Columns[Column].Data[k],
				},
				Schema: t.Schema,
			})
		}
	}
	return rows, nil
}

func (t Table) GetAvg(Column string) (float64, error) {
	r := t.Columns[Column].Avg
	return r, nil
}

func (t Table) GetMode(Column string) (float64, error) {
	return 0.0, nil
}

func (t Table) Count(Column string) (int, error) {
	r := t.Columns[Column].Count
	return r, nil
}
