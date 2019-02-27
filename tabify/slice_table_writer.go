package tabify

import (
	"sort"
)

// SliceTableWriter to writes a json array
type SliceTableWriter struct {
	cols  []string
	row   map[string]interface{}
	table [][]interface{}
	len   int
}

// OpenRow implements TableWriter interface
func (w *SliceTableWriter) OpenRow() {
	w.row = make(map[string]interface{})
}

// Cell implements TableWriter interface
func (w *SliceTableWriter) Cell(k string, v interface{}) {
	if w.len == 0 {
		w.cols = append(w.cols, k)
	}
	w.row[k] = v
}

// CloseRow implements TableWriter interface
func (w *SliceTableWriter) CloseRow() {
	if w.len == 0 {
		sort.Strings(w.cols)
		c := make([]interface{}, len(w.cols))
		for i, col := range w.cols {
			c[i] = col
		}
		w.table = append(w.table, c)
	}

	r := make([]interface{}, len(w.cols))
	for i, col := range w.cols {
		if v, ok := w.row[col]; ok {
			r[i] = v
		}
	}

	w.table = append(w.table, r)
	w.len++
}

// Table to gets the output table
func (w *SliceTableWriter) Table() [][]interface{} {
	return w.table
}
