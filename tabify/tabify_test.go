package tabify

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/datasweet/jsonmap"
	"github.com/stretchr/testify/assert"
)

func TestMetricAgg(t *testing.T) {
	testfile(t, "metric_agg")
}

func TestHistogramAgg(t *testing.T) {
	testfile(t, "histogram_terms_agg")
}

func TestNested(t *testing.T) {
	testfile(t, "nested_agg")
}

func testfile(t *testing.T, filename string) {
	src := readJSON(t, "./tests/"+filename+".json").Get("aggregations")

	// Options excluder / formatter
	excluder := func(keys []string) bool {
		last := keys[len(keys)-1]
		return last == "doc_count_error_upper_bound" || last == "sum_other_doc_count"
	}

	formatter := func(keys []string) string {
		var nk []string
		for _, k := range keys {
			if k == "buckets" || k == "value" || k == "key" {
				continue
			}
			nk = append(nk, k)
		}

		return strings.Join(nk, "#")
	}

	// JSONTableWriter
	jt, err := JSON(src, KeyExcluder(excluder), KeyFormatter(formatter))
	if err != nil {
		t.Fatal("json tabify", err)
	}
	jtw := readJSON(t, "./tests/"+filename+"_expected.json")
	assert.JSONEq(t, jtw.Stringify(), jt.Stringify(), "json table writer")

	// SliceTableWriter
	st, err := Slice(src, KeyExcluder(excluder), KeyFormatter(formatter))
	if err != nil {
		t.Fatal("slice tabify", err)
	}
	mst, err := json.Marshal(st)
	if err != nil {
		t.Fatal("slice json marshal")
	}
	stw := readJSON(t, "./tests/"+filename+"_slice.json")
	assert.JSONEq(t, stw.Stringify(), string(mst[:]), "slice table writer")

	// MapTableWriter
	mt, err := Map(src, KeyExcluder(excluder), KeyFormatter(formatter))
	if err != nil {
		t.Fatal("map tabify", err)
	}
	mmt, err := json.Marshal(mt)
	if err != nil {
		t.Fatal("slice json marshal")
	}
	mtw := readJSON(t, "./tests/"+filename+"_expected.json")
	assert.JSONEq(t, mtw.Stringify(), string(mmt[:]), "map table writer")
}

// readJSON to read a json file and store to a jsonmap
func readJSON(t *testing.T, filename string) *jsonmap.Json {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("unable to read %s", filename)
	}

	j := jsonmap.FromBytes(data)
	return j
}
