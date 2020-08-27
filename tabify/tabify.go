package tabify

import (
	"errors"

	"github.com/datasweet/jsonmap"
)

// Tabify using a custom TableWriter
func Tabify(j *jsonmap.Json, writer TableWriter, opt ...Option) error {
	if jsonmap.IsNil(j) {
		return errors.New("nil json")
	}
	if writer == nil {
		return errors.New("nil table writer")
	}
	t := newTabify(opt...)
	if err := t.Compute(j, writer); err != nil {
		return err
	}
	return nil
}

// JSON to flatten a json
func JSON(j *jsonmap.Json, opt ...Option) ([]*jsonmap.Json, error) {
	writer := &jsonTableWriter{}
	if err := Tabify(j, writer, opt...); err != nil {
		return nil, err
	}
	return writer.JSON(), nil
}

// Slice to tabify into a slice array.
// Note : first row contains headers
func Slice(j *jsonmap.Json, opt ...Option) ([][]interface{}, error) {
	writer := &sliceTableWriter{}
	if err := Tabify(j, writer, opt...); err != nil {
		return nil, err
	}
	return writer.Table(), nil
}

// Map to tabify into a map array
func Map(j *jsonmap.Json, opt ...Option) ([]map[string]interface{}, error) {
	writer := &mapTableWriter{}
	if err := Tabify(j, writer, opt...); err != nil {
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
			//fmt.Println(strings.Repeat("   ", node.deep), "NEWR", node.key)
			tb.openRow()

		case endRow:
			//fmt.Println(strings.Repeat("   ", node.deep), "ENDR", node.key)
			tb.closeRow()
		default:
			//fmt.Println(strings.Repeat("   ", node.deep), "CELL ", node.key, " = ", node.value)
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
	if t.opts.KeyExcluder != nil && len(keys) > 0 && t.opts.KeyExcluder(keys) {
		return
	}
	if node.IsObject() {
		for _, key := range node.Keys() {
			t.collect(node.Get(key), append(keys, key)...)
		}
	} else if node.IsArray() {

		nvalues := node.Values()
		switch l := len(nvalues); l {
		case 0:
			return
		case 1:
			// treat one line array as object
			t.collect(nvalues[0], keys...)
		default:
			for _, item := range nvalues {
				t.nodes <- &nodeValue{
					eventType: startRow,
					key:       t.opts.KeyFormatter(keys),
					deep:      len(keys),
				}

				t.collect(item, keys...)

				t.nodes <- &nodeValue{
					eventType: endRow,
					key:       t.opts.KeyFormatter(keys),
					deep:      len(keys),
				}
			}
		}
	} else {
		t.nodes <- &nodeValue{
			eventType: readValue,
			key:       t.opts.KeyFormatter(keys),
			deep:      len(keys),
			value:     node.Data(),
		}
	}
}
