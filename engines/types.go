package engines

type ColumnType uint8

const (
	TypeGeneric = ColumnType(0)      // Generic column, every column should support this
	TypeNumeric = ColumnType(1 << 0) // Numeric column supporting float64, int64
	TypeTextual = ColumnType(1 << 1) // Textual column supporting strings
)

// typeOf resolves all supported types of the column
func TypeOf(column IColumn) (typ ColumnType) {
	if _, ok := column.(Numeric); ok {
		typ = typ | TypeNumeric
	}
	if _, ok := column.(Textual); ok {
		typ = typ | TypeTextual
	}
	return
}
