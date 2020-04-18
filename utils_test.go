package jsonmap_test

import (
	"testing"

	"github.com/datasweet/jsonmap"
	"github.com/stretchr/testify/assert"
)

func TestEscapePath(t *testing.T) {
	json := jsonmap.New()
	json.Set("message.raw", "hello world !")
	assert.JSONEq(t, `{ "message": { "raw": "hello world !" }}`, json.Stringify())

	json = jsonmap.New()
	json.Set(jsonmap.EscapePath("message.raw"), "hello world !")
	assert.JSONEq(t, `{ "message.raw": "hello world !" }`, json.Stringify())
}
