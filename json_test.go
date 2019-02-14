package jsonmap

import (
	"bytes"
	"encoding/json"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

const jsonTest = `
{
	"string": "hello",
	"bool": true,
	"number": 123,
	"array": [1,2,3,4,5],
	"object": {
		"test": "world",
		"sub": [
			{"a": 4, "1": "a" },
			{"a": 5, "1": "b" }
		]
	}
}
`

func TestNewJson(t *testing.T) {
	j := New()
	assert.False(t, j.IsNil())
	assert.Equal(t, "{}", j.Stringify())
}

func TestFromWrongString(t *testing.T) {
	j := FromString("hello")
	assert.True(t, j.IsNil())
	assert.Equal(t, "{}", j.Stringify())
}

func TestFromString(t *testing.T) {
	j := FromString(jsonTest)
	assert.False(t, j.IsNil())
	assert.False(t, j.IsValue())
	assert.False(t, j.IsArray())
	assert.True(t, j.IsObject())
}

func TestAsString(t *testing.T) {
	j := FromString(jsonTest)
	assert.Equal(t, "hello", j.Get("string").AsString())
	assert.Equal(t, "", j.Get("bool").AsString())
	assert.Equal(t, "", j.Get("number").AsString())
	assert.Equal(t, "", j.Get("array").AsString())
	assert.Equal(t, "", j.Get("object").AsString())
	assert.Equal(t, "", j.Get("unknown").AsString())
}

func TestAsBool(t *testing.T) {
	j := FromString(jsonTest)
	assert.False(t, j.Get("string").AsBool())
	assert.True(t, j.Get("bool").AsBool())
	assert.False(t, j.Get("number").AsBool())
	assert.False(t, j.Get("array").AsBool())
	assert.False(t, j.Get("object").AsBool())
	assert.False(t, j.Get("unknown").AsBool())
}

func TestAsInt(t *testing.T) {
	j := FromString(jsonTest)
	assert.Equal(t, int64(0), j.Get("string").AsInt())
	assert.Equal(t, int64(0), j.Get("bool").AsInt())
	assert.Equal(t, int64(123), j.Get("number").AsInt())
	assert.Equal(t, int64(0), j.Get("array").AsInt())
	assert.Equal(t, int64(0), j.Get("object").AsInt())
	assert.Equal(t, int64(0), j.Get("unknown").AsInt())
}

func TestAsUint(t *testing.T) {
	j := FromString(jsonTest)
	assert.Equal(t, uint64(0), j.Get("string").AsUint())
	assert.Equal(t, uint64(0), j.Get("bool").AsUint())
	assert.Equal(t, uint64(123), j.Get("number").AsUint())
	assert.Equal(t, uint64(0), j.Get("array").AsUint())
	assert.Equal(t, uint64(0), j.Get("object").AsUint())
	assert.Equal(t, uint64(0), j.Get("unknown").AsUint())
}

func TestAsFloat64(t *testing.T) {
	j := FromString(jsonTest)
	assert.Equal(t, float64(0), j.Get("string").AsFloat())
	assert.Equal(t, float64(0), j.Get("bool").AsFloat())
	assert.Equal(t, float64(123), j.Get("number").AsFloat())
	assert.Equal(t, float64(0), j.Get("array").AsFloat())
	assert.Equal(t, float64(0), j.Get("object").AsFloat())
	assert.Equal(t, float64(0), j.Get("unknown").AsFloat())
}

func TestObjectKeys(t *testing.T) {
	j := FromString(jsonTest)

	expected := []string{"string", "bool", "number", "array", "object"}
	sort.Strings(expected)

	actual := j.Keys()
	sort.Strings(actual)

	assert.EqualValues(t, expected, actual)
}

func TestGetPath(t *testing.T) {
	j := FromString(jsonTest)
	assert.Equal(t, "a", j.Get("object.sub[0].1").AsString())
	assert.Equal(t, "b", j.Get("object.sub[1].1").AsString())
	assert.True(t, j.Get("object.test.sub2").IsNil())
	assert.True(t, j.Get("object.test.sub[2].1").IsNil())
}

func TestSet(t *testing.T) {
	j := New()
	assert.True(t, j.Set("", 3.14))
	assert.Equal(t, 3.14, j.AsFloat())

	// Collision test => try to sets object on value
	assert.True(t, j.Set("hello", "world"))
	assert.Equal(t, "world", j.Get("hello").AsString())

	// Can reinject new value
	assert.True(t, j.Set("", "3.14"))
	assert.Equal(t, "3.14", j.AsString())

	// Create auto map
	j = New()
	assert.True(t, j.Set("hello", "world"))
	assert.Equal(t, "world", j.Get("hello").AsString())
	assert.True(t, j.Set("the.number.pi.is", 3.14))
	assert.Equal(t, 3.14, j.Get("the.number.pi.is").AsFloat())

	// Can sets array
	j = FromString(jsonTest)
	assert.Equal(t, int64(2), j.Get("array[1]").AsInt())
	assert.True(t, j.Set("array[1]", 3.14))
	assert.Equal(t, 3.14, j.Get("array[1]").AsFloat())
	assert.Equal(t, "b", j.Get("object.sub[1].1").AsString())
	assert.True(t, j.Set("object.sub[1].1", 3.14))
	assert.Equal(t, 3.14, j.Get("object.sub[1].1").AsFloat())

	// can create auto array
	assert.True(t, j.Set("hello", 1, 2, 3, 4, 5))
	assert.Subset(t, []int{1, 2, 3, 4, 5}, j.Get("hello").AsArray())
}

func TestWrap(t *testing.T) {
	j := New()
	assert.True(t, j.Set("pi", 3.14))

	wrap := j.Wrap("the.best")
	assert.False(t, j.IsNil())
	assert.JSONEq(t, `{ "the": { "best": { "pi": 3.14 } } }`, wrap.Stringify())

}

func TestForEachArray(t *testing.T) {
	j := FromString("[1, 2]")
	j.ForEach(func(k string, v *Json) bool {

		switch k {
		case "0":
			assert.Equal(t, int64(1), v.AsInt())
		case "1":
			assert.Equal(t, int64(2), v.AsInt())
		default:
			assert.Fail(t, "unknown key '%s'", k)

		}
		return true
	})
}

func TestForEachObject(t *testing.T) {
	j := FromString(`{ "a": 1, "b": 2 }`)
	j.ForEach(func(k string, v *Json) bool {
		switch k {
		case "a":
			assert.Equal(t, int64(1), v.AsInt())
		case "b":
			assert.Equal(t, int64(2), v.AsInt())
		default:
			assert.Fail(t, "unknown key '%s'", k)

		}
		return true
	})
}

func TestKeys(t *testing.T) {
	keys := FromString("[1, 2]").Keys()
	assert.Equal(t, 2, len(keys))
	assert.Subset(t, []string{"0", "1"}, keys)

	keys = FromString(`{ "a": 1, "b": 2 }`).Keys()
	assert.Equal(t, 2, len(keys))
	assert.Subset(t, []string{"a", "b"}, keys)
}

func TestValues(t *testing.T) {
	values := FromString("[1, 2]").Values()
	assert.Equal(t, 2, len(values))
	var d1 []int64
	for _, v := range values {
		d1 = append(d1, v.AsInt())
	}
	assert.Subset(t, []int64{1, 2}, d1)

	values = FromString(`{ "a": 1, "b": 2 }`).Values()
	assert.Equal(t, 2, len(values))
	var d2 []int64
	for _, v := range values {
		d2 = append(d2, v.AsInt())
	}
	assert.Subset(t, []int64{1, 2}, d2)
}

func TestUnset(t *testing.T) {
	j := FromString(jsonTest)
	assert.False(t, j.Get("object.sub[0].a").IsNil())
	assert.True(t, j.Unset("object.sub[0].a"))
	assert.True(t, j.Get("object.sub[0].a").IsNil())

	assert.False(t, j.Get("object.sub[1]").IsNil())
	assert.True(t, j.Unset("object.sub[1]"))
	assert.True(t, j.Get("object.sub[1]").IsNil())
}

func TestRewrite(t *testing.T) {
	j := FromString(jsonTest)
	assert.False(t, j.Get("object.sub[0].a").IsNil())
	assert.True(t, j.Rewrite("object.sub[0].a", "new"))
	assert.True(t, j.Get("object.sub[0].a").IsNil())
	assert.False(t, j.Get("new").IsNil())
}

type Dummy struct {
	Raw *Json `json:"raw"`
}

func TestEncodeDecode(t *testing.T) {
	v := &Dummy{
		Raw: FromString(jsonTest),
	}

	// encode
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	assert.NoError(t, enc.Encode(v))
	assert.JSONEq(t, FromString(jsonTest).Wrap("raw").Stringify(), buf.String())

	var dummy Dummy
	assert.NoError(t, json.Unmarshal(buf.Bytes(), &dummy))
	assert.JSONEq(t, FromString(jsonTest).Stringify(), dummy.Raw.Stringify())
}
