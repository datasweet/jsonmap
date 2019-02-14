package jsonmap

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricAgg(t *testing.T) {
	testfile(t, "metric_agg", "aggregations")
}

func TestHistogramAgg(t *testing.T) {
	testfile(t, "histogram_terms_agg", "aggregations")
}

func TestNested(t *testing.T) {
	testfile(t, "nested_agg", "aggregations")
}

func testfile(t *testing.T, filename string, path string) {
	src, err := readJSON("./tests/" + filename + ".json")
	if err != nil {
		t.Fatal("Unable to read source file ", filename, err)
	}

	if len(path) > 0 {
		src = src.Get(path)
	}

	expected, err := readJSON("./tests/" + filename + "_expected.json")
	if err != nil {
		t.Fatal("Unable to read expected file ", filename, err)
	}

	tabified, err := Tabify(src,
		TabifyKeyExcluder(func(keys []string) bool {
			last := keys[len(keys)-1]
			return last == "doc_count_error_upper_bound" || last == "sum_other_doc_count"
		}),
		TabifyKeyFormatter(func(keys []string) string {
			var nk []string
			for _, k := range keys {
				if k == "buckets" || k == "value" || k == "key" {
					continue
				}
				nk = append(nk, k)
			}

			return strings.Join(nk, "#")
		}),
	)
	if err != nil {
		t.Fatal("Error during tabify", err)
	}

	actual, err := json.Marshal(tabified)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, expected.Stringify(), string(actual[:]))
}

// readJSON to read a json file and store to a jsonmap
func readJSON(filename string) (*Json, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	j := FromBytes(data)
	return j, nil
}
