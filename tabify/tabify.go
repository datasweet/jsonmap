package tabify

import (
	"errors"

	"github.com/datasweet/jsonmap"
)

// TableWriter is an interface to define a table writer
type TableWriter interface {
	OpenRow()
	Cell(k string, v interface{})
	CloseRow()
}

// JSON to flatten a json
func JSON(j *jsonmap.Json, opt ...Option) ([]*jsonmap.Json, error) {
	t := newTabify(opt...)
	writer := &JSONTableWriter{}

	if err := t.Compute(j, writer); err != nil {
		return nil, err
	}

	return writer.JSON(), nil
}

// Slice to tabify into a slice array.
// Note : first row contains headers
func Slice(j *jsonmap.Json, opt ...Option) ([][]interface{}, error) {
	t := newTabify(opt...)
	writer := &SliceTableWriter{}

	if err := t.Compute(j, writer); err != nil {
		return nil, err
	}

	return writer.Table(), nil
}

// Map to tabify into a map array
func Map(j *jsonmap.Json, opt ...Option) ([]map[string]interface{}, error) {
	t := newTabify(opt...)
	writer := &MapTableWriter{}

	if err := t.Compute(j, writer); err != nil {
		return nil, err
	}

	return writer.Table(), nil
}

// tabify is our main implementation
type tabify struct {
	opts  Options
	nodes chan *nodeValue
}

func newTabify(opt ...Option) *tabify {
	opts := newOptions(opt...)

	return &tabify{
		opts: opts,
	}
}

func (t *tabify) Options() Options {
	return t.opts
}

func (t *tabify) Compute(json *jsonmap.Json, tw TableWriter) error {
	if jsonmap.IsNil(json) {
		return errors.New("no json provided")
	}

	t.nodes = make(chan *nodeValue)
	tb := newTableBuffer()

	go func() {
		defer close(t.nodes)
		t.collect(json)
	}()

	// Listen new node entry
	for node := range t.nodes {

		// We need a new row !
		switch node.eventType {
		case startRow:
			// fmt.Println("START ROW", node.key)
			tb.openRow()

		case endRow:
			// fmt.Println("END ROW", node.key)
			tb.closeRow()
		default:
			// fmt.Println("VALUE", node.key, "\t\t = ", node.value)
			tb.cell(node)
		}
	}

	// Write table
	tb.write(tw)

	return nil
}

// collects node in json
func (t *tabify) collect(node *jsonmap.Json, keys ...string) {
	if jsonmap.IsNil(node) {
		return
	}
	if node.IsObject() {
		for _, key := range node.Keys() {
			t.collect(node.Get(key), append(keys, key)...)
		}
	} else if node.IsArray() {
		for _, item := range node.Values() {
			t.nodes <- newStartRow(t.opts.KeyFormatter(keys))
			t.collect(item, keys...)
			t.nodes <- newEndRow(t.opts.KeyFormatter(keys))
		}
	} else if t.opts.KeyExcluder == nil || !t.opts.KeyExcluder(keys) {
		t.nodes <- newNodeValue(t.opts.KeyFormatter(keys), node.Data())
	}
}
