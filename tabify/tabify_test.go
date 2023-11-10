package tabify_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/datasweet/jsonmap"
	"github.com/datasweet/jsonmap/tabify"
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

func TestDoubleTerms(t *testing.T) {
	testfile(t, "terms_terms_sum_agg")
}

func TestTripleAgg(t *testing.T) {
	testfile(t, "triple_agg")
}

func testfile(t *testing.T, filename string) {
	jsonf := readJSON(t, "./tests/"+filename+".json")
	src := jsonmap.FromString(jsonf).Get("aggregations")

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

	// SliceTableWriter
	t.Run("can write to slice "+filename, func(t *testing.T) {
		st, err := tabify.Slice(src, tabify.KeyExcluder(excluder), tabify.KeyFormatter(formatter))
		assert.NoError(t, err, "slice table writer")
		stw := readJSON(t, "./tests/"+filename+"_slice.json")
		var arr []interface{}
		assert.NoError(t, json.Unmarshal([]byte(stw), &arr))
		assert.Equal(t, len(arr), len(st))
		for i, row := range arr {
			assert.Equalf(t, row, st[i], "at index %d", i)
		}

		//assert.JSONEq(t, stw, mst.String(), "slice table writer")
	})

	// MapTableWriter
	t.Run("can write to map "+filename, func(t *testing.T) {
		mt, err := tabify.Map(src, tabify.KeyExcluder(excluder), tabify.KeyFormatter(formatter))
		assert.NoError(t, err, "map table writer")
		mmt, err := json.Marshal(mt)
		if err != nil {
			t.Fatal("slice json marshal")
		}
		mtw := readJSON(t, "./tests/"+filename+"_expected.json")
		assert.JSONEq(t, mtw, string(mmt[:]), "map table writer")
	})

	// JSONTableWriter
	t.Run("can write to json "+filename, func(t *testing.T) {
		jt, err := tabify.JSON(src, tabify.KeyExcluder(excluder), tabify.KeyFormatter(formatter))
		assert.NoError(t, err, "json table writer")
		jtw := readJSON(t, "./tests/"+filename+"_expected.json")
		jo := jsonmap.New()
		jo.Set("", jt)
		assert.JSONEq(t, jtw, jo.Stringify(), "json table writer")
	})
}

// readJSON to read a json file and store to a jsonmap
func readJSON(t *testing.T, filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("unable to read %s", filename)
	}
	return string(data[:])
}
