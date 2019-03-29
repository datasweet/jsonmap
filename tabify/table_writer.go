package tabify

import (
	"sort"

	"github.com/datasweet/jsonmap"
)

// TableWriter is an interface to define a table writer
type TableWriter interface {
	OpenRow()
	Cell(k string, v interface{})
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

func (w *jsonTableWriter) Cell(k string, v interface{}) {
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

func (w *mapTableWriter) Cell(k string, v interface{}) {
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
	cols        []string
	row         map[string]interface{}
	table       [][]interface{}
	len         int
	WithHeaders bool
}

func (w *sliceTableWriter) OpenRow() {
	w.row = make(map[string]interface{})
}

func (w *sliceTableWriter) Cell(k string, v interface{}) {
	if w.len == 0 {
		w.cols = append(w.cols, k)
	}
	w.row[k] = v
}

func (w *sliceTableWriter) CloseRow() {
	if w.WithHeaders && w.len == 0 {
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

func (w *sliceTableWriter) Table() [][]interface{} {
	return w.table
}
