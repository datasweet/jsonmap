package jsonmap_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/datasweet/jsonmap"
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

type Person struct {
	FirstName string
	Name      string
	Age       int
}

func (p *Person) JSON() *jsonmap.Json {
	j := jsonmap.New()
	j.Set("name", fmt.Sprintf("%s %s", p.FirstName, p.Name))
	j.Set("age", p.Age)
	return j
}

func TestNewJson(t *testing.T) {
	j := jsonmap.New()
	assert.False(t, j.IsNil())
	assert.Equal(t, "{}", j.Stringify())
}

func TestFromWrongString(t *testing.T) {
	j := jsonmap.FromString("hello")
	assert.True(t, j.IsNil())
	assert.Equal(t, "null", j.Stringify())
}

func TestFromString(t *testing.T) {
	j := jsonmap.FromString(jsonTest)
	assert.False(t, j.IsNil())
	assert.False(t, j.IsValue())
	assert.False(t, j.IsArray())
	assert.True(t, j.IsObject())
}

func TestAsString(t *testing.T) {
	j := jsonmap.FromString(jsonTest)
	assert.Equal(t, "hello", j.Get("string").AsString())
	assert.Equal(t, "", j.Get("bool").AsString())
	assert.Equal(t, "", j.Get("number").AsString())
	assert.Equal(t, "", j.Get("array").AsString())
	assert.Equal(t, "", j.Get("object").AsString())
	assert.Equal(t, "", j.Get("unknown").AsString())
}

func TestAsBool(t *testing.T) {
	j := jsonmap.FromString(jsonTest)
	assert.False(t, j.Get("string").AsBool())
	assert.True(t, j.Get("bool").AsBool())
	assert.False(t, j.Get("number").AsBool())
	assert.False(t, j.Get("array").AsBool())
	assert.False(t, j.Get("object").AsBool())
	assert.False(t, j.Get("unknown").AsBool())
}

func TestAsInt(t *testing.T) {
	j := jsonmap.FromString(jsonTest)
	assert.Equal(t, int64(0), j.Get("string").AsInt())
	assert.Equal(t, int64(0), j.Get("bool").AsInt())
	assert.Equal(t, int64(123), j.Get("number").AsInt())
	assert.Equal(t, int64(0), j.Get("array").AsInt())
	assert.Equal(t, int64(0), j.Get("object").AsInt())
	assert.Equal(t, int64(0), j.Get("unknown").AsInt())
}

func TestAsUint(t *testing.T) {
	j := jsonmap.FromString(jsonTest)
	assert.Equal(t, uint64(0), j.Get("string").AsUint())
	assert.Equal(t, uint64(0), j.Get("bool").AsUint())
	assert.Equal(t, uint64(123), j.Get("number").AsUint())
	assert.Equal(t, uint64(0), j.Get("array").AsUint())
	assert.Equal(t, uint64(0), j.Get("object").AsUint())
	assert.Equal(t, uint64(0), j.Get("unknown").AsUint())
}

func TestAsFloat64(t *testing.T) {
	j := jsonmap.FromString(jsonTest)
	assert.Equal(t, float64(0), j.Get("string").AsFloat())
	assert.Equal(t, float64(0), j.Get("bool").AsFloat())
	assert.Equal(t, float64(123), j.Get("number").AsFloat())
	assert.Equal(t, float64(0), j.Get("array").AsFloat())
	assert.Equal(t, float64(0), j.Get("object").AsFloat())
	assert.Equal(t, float64(0), j.Get("unknown").AsFloat())
}

func TestObjectKeys(t *testing.T) {
	j := jsonmap.FromString(jsonTest)

	expected := []string{"string", "bool", "number", "array", "object"}
	sort.Strings(expected)

	actual := j.Keys()
	sort.Strings(actual)

	assert.EqualValues(t, expected, actual)
}

func TestGetPath(t *testing.T) {
	j := jsonmap.FromString(jsonTest)
	assert.Equal(t, "a", j.Get("object.sub[0].1").AsString())
	assert.Equal(t, "b", j.Get("object.sub[1].1").AsString())
	assert.True(t, j.Get("object.test.sub2").IsNil())
	assert.True(t, j.Get("object.test.sub[2].1").IsNil())
}

func TestSet(t *testing.T) {
	t.Run("can set root path", func(t *testing.T) {
		j := jsonmap.New()
		assert.True(t, j.Set("", 3.14))
		assert.JSONEq(t, "3.14", j.Stringify())
		assert.Equal(t, 3.14, j.AsFloat())
	})

	t.Run("can replace a value by an object", func(t *testing.T) {
		j := jsonmap.New()
		assert.True(t, j.Set("", 3.14))
		assert.True(t, j.Set("hello", "world"))
		assert.JSONEq(t, `{ "hello": "world" }`, j.Stringify())
	})

	t.Run("can replace an object by a value", func(t *testing.T) {
		j := jsonmap.New()
		assert.True(t, j.Set("hello", "world"))
		assert.True(t, j.Set("", 3.14))
		assert.JSONEq(t, "3.14", j.Stringify())
	})

	t.Run("can create sub path", func(t *testing.T) {
		j := jsonmap.New()
		assert.True(t, j.Set("hello", "world"))
		assert.True(t, j.Set("the.number.pi.is", 3.14))
		assert.JSONEq(t, `{ "hello": "world", "the": { "number": { "pi": { "is": 3.14 }}}}`, j.Stringify())
	})

	t.Run("can create path with escape '.'", func(t *testing.T) {
		j := jsonmap.New()
		assert.True(t, j.Set("test\\.machin", 45))
		assert.True(t, j.Set("choux\\.machin.truc", "bidule"))
		assert.True(t, j.Set("choux.machin\\.truc", "bidule"))
		assert.JSONEq(t, `{ "test.machin": 45, "choux.machin": { "truc": "bidule" }, "choux": { "machin.truc": "bidule" }}`, j.Stringify())
	})

	t.Run("can set array", func(t *testing.T) {
		j := jsonmap.New()
		assert.True(t, j.Set("items", []int{1, 2, 3, 4, 5}))
		assert.JSONEq(t, `{ "items": [1, 2, 3, 4, 5] }`, j.Stringify())

		assert.True(t, j.Set("items[2]", 3.14))
		assert.JSONEq(t, `{ "items": [1, 2, 3.14, 4, 5] }`, j.Stringify())

		assert.True(t, j.Set("items[3]", &Person{FirstName: "Thomas", Name: "CHARLOT", Age: 36}))
		assert.JSONEq(t, `{ "items": [1, 2, 3.14, {"name": "Thomas CHARLOT", "age": 36 }, 5] }`, j.Stringify())

		assert.True(t, j.Set("items[3].age", 37))
		assert.JSONEq(t, `{ "items": [1, 2, 3.14, {"name": "Thomas CHARLOT", "age": 37 }, 5] }`, j.Stringify())

		assert.False(t, j.Set("items[7]", 11))
		assert.JSONEq(t, `{ "items": [1, 2, 3.14, {"name": "Thomas CHARLOT", "age": 37 }, 5] }`, j.Stringify())
	})

	t.Run("can set nil", func(t *testing.T) {
		j := jsonmap.New()
		assert.True(t, j.Set("nil", nil))
		assert.JSONEq(t, `{ "nil": null }`, j.Stringify())
	})

	t.Run("can set a sub json", func(t *testing.T) {
		j := jsonmap.New()
		sub := jsonmap.New()
		assert.True(t, sub.Set("pi", 3.14))
		assert.True(t, j.Set("wrapped", sub))
		assert.JSONEq(t, `{ "wrapped": { "pi": 3.14 }}`, j.Stringify())
	})

	t.Run("can set a json array", func(t *testing.T) {
		j := jsonmap.New()
		items := []*jsonmap.Json{
			jsonmap.FromString(`{ "string": "hello" }`),
			jsonmap.FromString(`{ "bool": true }`),
			jsonmap.FromString(`{ "number": 3.14 }`),
			jsonmap.FromString(`{ "array": [1,2,3,4,5] }`),
			jsonmap.FromString(`{ "object": { "a": 4, "1": "a" }}`),
			jsonmap.Nil(),
		}

		assert.True(t, j.Set("items", items))
		assert.JSONEq(t,
			`{ "items": [{ "string": "hello" }, { "bool": true }, { "number": 3.14 }, { "array": [1,2,3,4,5] }, { "object": { "a": 4, "1": "a" }}, null ]}`,
			j.Stringify(),
		)
	})

	t.Run("can set a jsonizer", func(t *testing.T) {
		j := jsonmap.New()
		person := &Person{FirstName: "Thomas", Name: "CHARLOT", Age: 36}
		assert.True(t, j.Set("person", person))
		assert.JSONEq(t,
			`{ "person": { "name": "Thomas CHARLOT", "age": 36 }}`,
			j.Stringify(),
		)
	})

	t.Run("can set a nil jsonizer", func(t *testing.T) {
		j := jsonmap.New()
		var person *Person
		assert.Nil(t, person)
		assert.True(t, j.Set("person", person))
		assert.JSONEq(t,
			`{ "person": null }`,
			j.Stringify(),
		)
	})

	t.Run("can set an array of jsonizer", func(t *testing.T) {
		j := jsonmap.New()
		peoples := []*Person{
			&Person{FirstName: "Thomas", Name: "CHARLOT", Age: 36},
			nil,
			&Person{FirstName: "Lionel", Name: "FROMENT", Age: 46},
		}
		assert.True(t, j.Set("peoples", peoples))
		assert.JSONEq(t,
			`{ "peoples": [{ "name": "Thomas CHARLOT", "age": 36 }, null, { "name": "Lionel FROMENT", "age": 46 }]}`,
			j.Stringify(),
		)
	})
}

func TestWrap(t *testing.T) {
	j := jsonmap.New()
	assert.True(t, j.Set("pi", 3.14))

	wrap := j.Wrap("the.best")
	assert.False(t, j.IsNil())
	assert.JSONEq(t, `{ "the": { "best": { "pi": 3.14 } } }`, wrap.Stringify())

}

func TestForEachArray(t *testing.T) {
	j := jsonmap.FromString("[1, 2]")
	j.ForEach(func(k string, v *jsonmap.Json) bool {

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
	j := jsonmap.FromString(`{ "a": 1, "b": 2 }`)
	j.ForEach(func(k string, v *jsonmap.Json) bool {
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
	keys := jsonmap.FromString("[1, 2]").Keys()
	assert.Equal(t, 2, len(keys))
	assert.Subset(t, []string{"0", "1"}, keys)

	keys = jsonmap.FromString(`{ "a": 1, "b": 2 }`).Keys()
	assert.Equal(t, 2, len(keys))
	assert.Subset(t, []string{"a", "b"}, keys)
}

func TestValues(t *testing.T) {
	values := jsonmap.FromString("[1, 2]").Values()
	assert.Equal(t, 2, len(values))
	var d1 []int64
	for _, v := range values {
		d1 = append(d1, v.AsInt())
	}
	assert.Subset(t, []int64{1, 2}, d1)

	values = jsonmap.FromString(`{ "a": 1, "b": 2 }`).Values()
	assert.Equal(t, 2, len(values))
	var d2 []int64
	for _, v := range values {
		d2 = append(d2, v.AsInt())
	}
	assert.Subset(t, []int64{1, 2}, d2)
}

func TestUnset(t *testing.T) {
	j := jsonmap.FromString(jsonTest)
	assert.False(t, j.Get("object.sub[0].a").IsNil())
	assert.True(t, j.Unset("object.sub[0].a"))
	assert.True(t, j.Get("object.sub[0].a").IsNil())

	assert.False(t, j.Get("object.sub[1]").IsNil())
	assert.True(t, j.Unset("object.sub[1]"))
	assert.True(t, j.Get("object.sub[1]").IsNil())
}

func TestRewrite(t *testing.T) {
	j := jsonmap.FromString(jsonTest)
	assert.False(t, j.Get("object.sub[0].a").IsNil())
	assert.True(t, j.Rewrite("object.sub[0].a", "new"))
	assert.True(t, j.Get("object.sub[0].a").IsNil())
	assert.False(t, j.Get("new").IsNil())
}

type Dummy struct {
	Raw *jsonmap.Json `json:"raw"`
}

func TestEncodeDecode(t *testing.T) {
	v := &Dummy{
		Raw: jsonmap.FromString(jsonTest),
	}

	bytes, err := json.Marshal(v)
	assert.NoError(t, err)
	assert.NotNil(t, bytes)

	// // encode
	// buf := &bytes.Buffer{}
	// enc := json.NewEncoder(buf)
	// enc.SetEscapeHTML(true)	data, err := json.Marshal(v)
	// assert.NoError(t, enc.Encode(v))
	// assert.JSONEq(t, FromString(jsonTest).Wrap("raw").Stringify(), buf.String())

	var dummy Dummy
	assert.NoError(t, json.Unmarshal(bytes, &dummy))
	assert.JSONEq(t, jsonmap.FromString(jsonTest).Stringify(), dummy.Raw.Stringify())
}

func TestClone(t *testing.T) {
	j := jsonmap.FromString(jsonTest)

	clone := j.Clone()

	assert.False(t, jsonmap.IsNil(clone))

	j.Unset("string")
	assert.Equal(t, "", j.Get("string").AsString())

	clone.Set("test", 12345)

	assert.Equal(t, "hello", clone.Get("string").AsString())
	assert.Equal(t, true, clone.Get("bool").AsBool())
	assert.Equal(t, float64(123), clone.Get("number").AsFloat())
	assert.Equal(t, int64(4), clone.Get("object.sub[0].a").AsInt())
	assert.Equal(t, int64(12345), clone.Get("test").AsInt())
	assert.Equal(t, int64(0), j.Get("test").AsInt())
}

func TestMerge(t *testing.T) {
	j := jsonmap.FromString(jsonTest)
	j2 := jsonmap.FromString(`{
			"new": "value",
			"name": "john"
		}`)

	merge := jsonmap.Merge(j, j2)
	assert.NotNil(t, merge)
	assert.Equal(t, "hello", merge.Get("string").AsString())
	assert.Equal(t, true, merge.Get("bool").AsBool())
	assert.Equal(t, float64(123), merge.Get("number").AsFloat())
	assert.Equal(t, int64(4), merge.Get("object.sub[0].a").AsInt())
	assert.Equal(t, "value", merge.Get("new").AsString())
	assert.Equal(t, "john", merge.Get("name").AsString())
	assert.Len(t, merge.Keys(), 7)
}
