package tabify

// nodeValue is a helper type to store a final value of json
// defined by the keys[] -  ie the json path to access to this value.
type nodeValue struct {
	eventType nodeEventType
	key       string
	value     interface{}
	deep      int
}

type nodeEventType uint8

const (
	readValue nodeEventType = iota
	startRow
	endRow
)
