package tabify

import "github.com/datasweet/jsonmap"

// JSONTableWriter to writes a json array
type JSONTableWriter struct {
	row   *jsonmap.Json
	table []*jsonmap.Json
}

// OpenRow implements TableWriter interface
func (w *JSONTableWriter) OpenRow() {
	w.row = jsonmap.New()
}

// Cell implements TableWriter interface
func (w *JSONTableWriter) Cell(k string, v interface{}) {
	if w.row == nil {
		w.OpenRow()
	}
	w.row.Set(k, v)
}

// CloseRow implements TableWriter interface
func (w *JSONTableWriter) CloseRow() {
	w.table = append(w.table, w.row)
}

// Table to gets the output table
func (w *JSONTableWriter) JSON() []*jsonmap.Json {
	return w.table
}
