package tabify

import "github.com/datasweet/jsonmap"

// JSONTableWriter to writes a json array
type JSONTableWriter struct {
	row   map[string]interface{}
	table []map[string]interface{}
}

// OpenRow implements TableWriter interface
func (w *JSONTableWriter) OpenRow() {
	w.row = make(map[string]interface{})
}

// Cell implements TableWriter interface
func (w *JSONTableWriter) Cell(k string, v interface{}) {
	if w.row == nil {
		w.OpenRow()
	}
	w.row[k] = v
}

// CloseRow implements TableWriter interface
func (w *JSONTableWriter) CloseRow() {
	w.table = append(w.table, w.row)
}

// Table to gets the output table
func (w *JSONTableWriter) JSON() *jsonmap.Json {
	j := jsonmap.New()
	j.Set("", w.table)
	return j
}
