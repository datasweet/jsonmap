package jsonmap

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

// Json is our wrapper to an unmarshalled json
type Json struct {
	data interface{}
}

// New creates an empty object Json, ie {}
func New() *Json {
	return &Json{
		data: make(map[string]interface{}),
	}
}

// Nil creates an nil Json
func Nil() *Json {
	return &Json{nil}
}

// FromBytes to creates a Json from bytes
func FromBytes(bytes []byte) *Json {
	j := new(Json)
	err := j.UnmarshalJSON(bytes)
	if err != nil {
		return Nil()
	}
	return j
}

// FromString to creates a Json from a string
func FromString(str string) *Json {
	return FromBytes([]byte(str))
}

// FromMap to creates a Json from an unmarshalled map
func FromMap(m map[string]interface{}) *Json {
	if m == nil {
		return Nil()
	}

	return &Json{
		data: m,
	}
}

// Stringify formats current node to a json string
func (j *Json) Stringify() string {
	return string(j.Bytes())
}

// Bytes return json bytes
func (j *Json) Bytes() []byte {
	bytes, _ := j.MarshalJSON()
	return bytes
}

// Marshaler interface encoding/json encode.go
func (j *Json) MarshalJSON() ([]byte, error) {
	if j.IsNil() {
		return []byte("{}"), nil
	}
	bytes, err := json.Marshal(j.data)
	if err != nil {
		return []byte("{}"), err
	}
	return bytes, nil
}

// Unmarshaler interface encoding/json decode.go
func (j *Json) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("unmarshal JSON on a nil pointer")
	}
	return json.Unmarshal(data, &j.data)
}

// Data get uncasted data
func (j *Json) Data() interface{} {
	return j.data
}

// IsNil to check if the current Json is nil
func (j *Json) IsNil() bool {
	return j.data == nil
}

// IsObject to know if the current Json is an object
func (j *Json) IsObject() bool {
	_, ok := (j.data).(map[string]interface{})
	return ok
}

// AsObject casts underlying to object (map[string]interface{})
// Returns nil if not an object
func (j *Json) AsObject() map[string]interface{} {
	if casted, ok := (j.data).(map[string]interface{}); ok {
		return casted
	}
	return nil
}

// IsArray to check if the current Json is an array
func (j *Json) IsArray() bool {
	_, ok := (j.data).([]interface{})
	return ok
}

// AsArray casts underlying to array ([]interface{})
// Returns nil if not an array
func (j *Json) AsArray() []interface{} {
	if casted, ok := (j.data).([]interface{}); ok {
		return casted
	}
	return nil
}

// IsValue to check if the current Json is a type value
func (j *Json) IsValue() bool {
	return !j.IsNil() && !j.IsObject() && !j.IsArray()
}

// AsString casts underlying to string
// Returns an empty string if not a string
func (j *Json) AsString() string {
	if casted, ok := (j.data).(string); ok {
		return casted
	}
	return ""
}

// AsBool casts underlying to boolean
// Returns false if not a boolean
func (j *Json) AsBool() bool {
	if casted, ok := (j.data).(bool); ok {
		return casted
	}
	return false
}

// AsInt casts underlying to int64
// Returns 0 if not an int
func (j *Json) AsInt() int64 {
	switch j.data.(type) {
	case json.Number:
		if i, err := (j.data).(json.Number).Int64(); err != nil {
			return i
		}
		return 0
	case float32, float64:
		return int64(reflect.ValueOf(j.data).Float())
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(j.data).Int()
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(j.data).Uint())
	default:
		return 0
	}
}

// AsUint casts underlying to uint64
// Returns 0 if not an int
func (j *Json) AsUint() uint64 {
	switch j.data.(type) {
	case json.Number:
		if u, err := strconv.ParseUint(j.data.(json.Number).String(), 10, 64); err != nil {
			return u
		}
		return 0
	case float32, float64:
		return uint64(reflect.ValueOf(j.data).Float())
	case int, int8, int16, int32, int64:
		return uint64(reflect.ValueOf(j.data).Int())
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(j.data).Uint()
	default:
		return 0
	}
}

// AsFloat casts underlying to float64
// Returns 0 if not a float
func (j *Json) AsFloat() float64 {
	switch j.data.(type) {
	case json.Number:
		if f, err := (j.data).(json.Number).Float64(); err != nil {
			return f
		}
		return 0
	case float32, float64:
		return reflect.ValueOf(j.data).Float()
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(j.data).Int())
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(j.data).Uint())
	default:
		return 0
	}
}

// Get gets the value at path of object. If not found returns Nils() value
func (j *Json) Get(path string) *Json {
	keys := createPath(path)
	curr := j
	for _, k := range keys {

		// Get  as object
		if o := curr.AsObject(); o != nil {
			val, ok := o[k]
			if !ok {
				return Nil()
			}
			curr = &Json{val}
			continue
		}

		// Get as array
		if a := curr.AsArray(); a != nil {
			// Must be an int
			idx, e := strconv.Atoi(k)
			if e != nil || idx < 0 || idx >= len(a) {
				return Nil()
			}
			curr = &Json{a[idx]}
			continue
		}

		// Not found
		return Nil()
	}

	return curr
}

// Has checks if path is a direct property of object.
func (j *Json) Has(path string) bool {
	o := j.Get(path)
	return !o.IsNil()
}

// Set sets the value at path of object. If a portion of path doesn't exist, it's created.
// Arrays are created for missing index properties while objects are created for all other missing properties
func (j *Json) Set(path string, value ...interface{}) bool {
	keys := createPath(path)
	lastIndex := len(keys) - 1

	// Pick value
	var newValue interface{}
	lv := len(value)

	if lv == 0 {
		newValue = nil
	} else if lv == 1 {
		newValue = value[0]
	} else {
		newValue = value
	}

	if lastIndex == -1 {
		j.data = newValue
		return true
	}

	curr := j
	for i, k := range keys {
		// Get as object
		if o := curr.AsObject(); o != nil {
			// Assign value
			if i == lastIndex {
				o[k] = newValue
				return true
			}

			if _, ok := o[k]; !ok {
				o[k] = make(map[string]interface{})
			}
			curr = &Json{o[k]}
			continue
		}

		// Get as array
		if a := curr.AsArray(); a != nil {
			// Must be an int
			idx, e := strconv.Atoi(k)
			if e != nil || idx < 0 {
				return false
			}

			if idx >= len(a) {
				for q := len(a); q <= idx; q++ {
					a = append(a, nil)
				}
			}

			// Assign value
			if i == lastIndex {
				a[idx] = newValue
				return true
			}

			curr = &Json{a[idx]}
			continue
		}

		// Value or nil => we force the rewrite
		m := make(map[string]interface{})
		curr.data = m

		// Assign
		if i == lastIndex {
			m[k] = newValue
			return true
		}

		m[k] = make(map[string]interface{})
		curr = &Json{m[k]}
	}

	return false
}

// SetJSON to sets a json or an array of json to path
func (j *Json) SetJSON(path string, json ...*Json) bool {
	var d []interface{}
	for _, o := range json {
		d = append(d, o.data)
	}
	return j.Set(path, d...)
}

// Unset deletes the value
func (j *Json) Unset(path string) bool {
	keys := createPath(path)
	curr := j
	lastIndex := len(keys) - 1

	for i, k := range keys {
		// Get  as object
		if o := curr.AsObject(); o != nil {
			if i == lastIndex {
				delete(o, k)
				return true
			}
			val, ok := o[k]
			if !ok {
				return false
			}
			curr = &Json{val}
			continue
		}

		// Get as array
		if a := curr.AsArray(); a != nil {
			// Must be an int
			idx, e := strconv.Atoi(k)
			if e != nil || idx < 0 || idx >= len(a) {
				return false
			}
			if i == lastIndex {
				a[idx] = nil
				return true
			}
			curr = &Json{a[idx]}
			continue
		}

		// Not found
		return false
	}
	return false
}

// Rewrite changes a path
func (j *Json) Rewrite(oldPath string, newPath string) bool {
	d := j.Get(oldPath).Data()
	return j.Unset(oldPath) && j.Set(newPath, d)
}

// Wrap the current json to a new json
// Example : { "pi": 3.14 }.Wrap("const") => { "const": { "pi": 3.14 }} }
// Returns the new parent or nilJson if error
func (j *Json) Wrap(path string) *Json {
	wrap := New()
	if !wrap.Set(path, j.data) {
		return Nil()
	}
	return wrap
}

// ForEach : Iterates over elements of collection and invokes iteratee for each element.
// Iteratee functions may exit iteration early by explicitly returning false.
func (j *Json) ForEach(iteratee func(k string, v *Json) bool) {
	if iteratee == nil {
		return
	}

	if o := j.AsObject(); o != nil {
		for k, v := range o {
			if !iteratee(k, &Json{v}) {
				break
			}
		}
	} else if a := j.AsArray(); a != nil {
		for i, v := range a {
			if !iteratee(strconv.Itoa(i), &Json{v}) {
				break
			}
		}
	}
}

// Keys : creates an array of the own property names of object.
func (j *Json) Keys() []string {
	var keys []string

	callback := func(k string, v *Json) bool {
		keys = append(keys, k)
		return true
	}

	j.ForEach(callback)
	return keys
}

// Values : creates an array of the own enumerable string keyed property values of object.
func (j *Json) Values() []*Json {
	var values []*Json

	callback := func(k string, v *Json) bool {
		values = append(values, v)
		return true
	}

	j.ForEach(callback)

	return values
}
