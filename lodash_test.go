package jsonmap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	values := Filter(FromString("[1, 2, 3, 4, 5]").Values(), func(v *Json) bool {
		return v.AsInt()%2 == 1
	})
	assert.Equal(t, 3, len(values))
	for i, v := range values {
		assert.Equal(t, int64(2*i+1), v.AsInt())
	}

	values = Filter(FromString(`[{ "value": 1 }, { "value": 2 }, { "value": 3 }, { "value": 4 }, { "value": 5 }]`).Values(), func(v *Json) bool {
		return v.Get("value").AsInt()%2 == 1
	})
	assert.Equal(t, 3, len(values))
	for i, v := range values {
		assert.JSONEq(t, fmt.Sprintf(`{ "value": %d }`, 2*i+1), v.Stringify())
	}
}

func TestMap(t *testing.T) {
	values := Map(FromString("[1, 2, 3, 4, 5]").Values(), func(v *Json) *Json {
		if v.AsInt()%2 == 1 {
			return v
		}
		return Nil()
	})
	assert.Equal(t, 5, len(values))
	for i, v := range values {
		if i%2 == 0 {
			assert.Equal(t, int64(i+1), v.AsInt())
		} else {
			assert.True(t, v.IsNil())
		}
	}
}

func TestAssign(t *testing.T) {
	json1 := FromString(`{
		"string": "hello",
		"bool": true,
		"number": 123
	}`)

	json2 := FromString(`{
		"array": [1,2,3,4,5]
	}`)

	json3 := FromString(`{
		"object": {
			"test": "world",
			"sub": [
				{"a": 4, "1": "a" },
				{"a": 5, "1": "b" }
			]
		}
	}`)

	j := Assign(FromString(`{"hello": "world" }`), json1, json2, json3)

	expected := FromString(`{
		"hello": "world",
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
	}`)

	assert.JSONEq(t, expected.Stringify(), j.Stringify())
}
