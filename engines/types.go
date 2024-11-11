package engines

type columnType uint8

const (
	typeGeneric = columnType(0)      // Generic column, every column should support this
	typeNumeric = columnType(1 << 0) // Numeric column supporting float64, int64 
	typeTextual = columnType(1 << 1) // Textual column supporting strings
)

// typeOf resolves all supported types of the column
func typeOf(column IColumn) (typ columnType) {
	if _, ok := column.(Numeric); ok {
		typ = typ | typeNumeric
	}
	if _, ok := column.(Textual); ok {
		typ = typ | typeTextual
	}
	return
}