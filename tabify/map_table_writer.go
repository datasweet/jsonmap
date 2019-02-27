package tabify

// MapTableWriter to writes to a map array
type MapTableWriter struct {
	row   map[string]interface{}
	table []map[string]interface{}
}

// OpenRow implements TableWriter interface
func (w *MapTableWriter) OpenRow() {
	w.row = make(map[string]interface{})
}

// Cell implements TableWriter interface
func (w *MapTableWriter) Cell(k string, v interface{}) {
	w.row[k] = v
}

// CloseRow implements TableWriter interface
func (w *MapTableWriter) CloseRow() {
	w.table = append(w.table, w.row)
}

// Table to gets the output table
func (w *MapTableWriter) Table() []map[string]interface{} {
	return w.table
}
