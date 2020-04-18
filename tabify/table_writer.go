package tabify

import (
	"sort"

	"github.com/datasweet/jsonmap"
)

// TableWriter is an interface to define a table writer
type TableWriter interface {
	OpenRow()
	Cell(k string, v interface{}, deep int)
	CloseRow()
}

// jsonTableWriter to write a json array
type jsonTableWriter struct {
	row   *jsonmap.Json
	table []*jsonmap.Json
}

func (w *jsonTableWriter) OpenRow() {
	w.row = jsonmap.New()
}

func (w *jsonTableWriter) Cell(k string, v interface{}, deep int) {
	if w.row == nil {
		w.OpenRow()
	}
	w.row.Set(k, v)
}

func (w *jsonTableWriter) CloseRow() {
	w.table = append(w.table, w.row)
}

func (w *jsonTableWriter) JSON() []*jsonmap.Json {
	return w.table
}

// mapTableWriter to write to a map
type mapTableWriter struct {
	row   map[string]interface{}
	table []map[string]interface{}
}

func (w *mapTableWriter) OpenRow() {
	w.row = make(map[string]interface{})
}

func (w *mapTableWriter) Cell(k string, v interface{}, deep int) {
	w.row[k] = v
}

func (w *mapTableWriter) CloseRow() {
	w.table = append(w.table, w.row)
}

func (w *mapTableWriter) Table() []map[string]interface{} {
	return w.table
}

// sliceTableWriter to writes into a slice
type sliceTableWriter struct {
	cols  []*sliceCol
	row   map[string]interface{}
	table [][]interface{}
	len   int
}

type sliceCol struct {
	name string
	deep int
}

func (w *sliceTableWriter) OpenRow() {
	w.row = make(map[string]interface{})
}

func (w *sliceTableWriter) Cell(k string, v interface{}, deep int) {
	if w.len == 0 {
		w.cols = append(w.cols, &sliceCol{k, deep})
	}
	w.row[k] = v
}

func (w *sliceTableWriter) CloseRow() {
	if w.len == 0 {
		// compute cols
		sort.Slice(w.cols, func(i, j int) bool {
			if w.cols[i].deep == w.cols[j].deep {
				return w.cols[i].name < w.cols[j].name
			}
			return w.cols[i].deep < w.cols[j].deep
		})

		c := make([]interface{}, len(w.cols))
		for i, col := range w.cols {
			c[i] = col.name
		}
		w.table = append(w.table, c)
	}

	r := make([]interface{}, len(w.cols))
	for i, col := range w.cols {
		if v, ok := w.row[col.name]; ok {
			r[i] = v
		}
	}

	w.table = append(w.table, r)
	w.len++
}

func (w *sliceTableWriter) Table() [][]interface{} {
	return w.table
}
