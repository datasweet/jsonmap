package jsonmap

import (
	"errors"
	"strings"
)

// tabify implementation
type tabify struct {
	excluder  TabifyKeyExcluderFunc
	formatter TabifyKeyFormatterFunc
	nodes     chan *nodeValue
}

type TabifyKeyFormatterFunc func([]string) string

type TabifyKeyExcluderFunc func([]string) bool

// nodeValue is a helper type to store a final value of json
// defined by the keys[] -  ie the json path to access to this value.
type nodeValue struct {
	eventType nodeEventType
	key       string
	value     interface{}
}

type nodeEventType int

const (
	readValue nodeEventType = iota
	startRow
	endRow
)

// newNodeValue to create a new node value
// a node value is a node with a primitive type
func newNodeValue(value interface{}, key string) *nodeValue {
	return &nodeValue{
		key:       key,
		eventType: readValue,
		value:     value,
	}
}

// startRow to signal our collector we need to create a  row
func newStartRow(key string) *nodeValue {
	return &nodeValue{
		key:       key,
		eventType: startRow,
	}
}

// endRow to signal our collector we need to end the row
func newEndRow(key string) *nodeValue {
	return &nodeValue{
		key:       key,
		eventType: endRow,
	}
}

// TabifyOption is an option setter to tabify
type TabifyOption func(t *tabify)

// TabifyKeyFormatter sets the key formatter
// default : func (keys []string) => strings.Join(keys, "#")
func TabifyKeyFormatter(v TabifyKeyFormatterFunc) TabifyOption {
	return func(t *tabify) {
		t.formatter = v
	}
}

// TabifyKeyExcluder sets the key excluder
// default :  func (keys []string) => false
func TabifyKeyExcluder(v TabifyKeyExcluderFunc) TabifyOption {
	return func(t *tabify) {
		t.excluder = v
	}
}

func defaultFormatter(keys []string) string {
	return strings.Join(keys, "#")
}

func defaultExcluder(keys []string) bool {
	return false
}

// Tabify a json
func Tabify(json *Json, opt ...TabifyOption) ([]map[string]interface{}, error) {
	t := &tabify{
		formatter: defaultFormatter,
		excluder:  nil,
	}

	for _, o := range opt {
		o(t)
	}

	// Be sure
	if t.formatter == nil {
		t.formatter = defaultFormatter
	}

	return t.Compute(json)
}

func (t *tabify) Compute(json *Json) ([]map[string]interface{}, error) {
	if IsNil(json) {
		return nil, errors.New("no json provided")
	}

	t.nodes = make(chan *nodeValue)
	table := newTableWriter()

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
			table.openRow()

		case endRow:
			// fmt.Println("END ROW", node.key)
			table.closeRow()
		default:
			// fmt.Println("VALUE", node.key, "\t\t = ", node.value)
			table.cell(node)
		}
	}

	return table.write(), nil
}

// collects node in json
func (t *tabify) collect(node *Json, keys ...string) {
	if IsNil(node) {
		return
	}
	if node.IsObject() {
		for _, key := range node.Keys() {
			t.collect(node.Get(key), append(keys, key)...)
		}
	} else if node.IsArray() {
		for _, item := range node.Values() {
			t.nodes <- newStartRow(t.formatter(keys))
			t.collect(item, keys...)
			t.nodes <- newEndRow(t.formatter(keys))
		}
	} else if t.excluder == nil || !t.excluder(keys) {
		t.nodes <- newNodeValue(node.Data(), t.formatter(keys))
	}
}
