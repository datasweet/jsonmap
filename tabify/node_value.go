package tabify

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
func newNodeValue(key string, value interface{}) *nodeValue {
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
