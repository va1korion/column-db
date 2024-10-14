package engines

type LogEngine struct {
	
}


func (l LogEngine) MakeTable() (Table, error) {
	return Table{}, ErrNotImplemented
}

func (l LogEngine) Insert() error {
	return ErrNotImplemented
}

func (l LogEngine) Delete() error {
	return ErrNotImplemented
}

func (l LogEngine) Update() error {
	return ErrNotImplemented

}

func (l LogEngine) Read() (*bytes.Buffer, error) {
	return nil, ErrNotImplemented

}

func (l LogEngine) Save() (int, error) {
	return 0, ErrNotImplemented
}

