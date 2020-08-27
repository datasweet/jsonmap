package tabify_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	// MapTableWriter
	// mt, err := tabify.Map(src, tabify.KeyExcluder(excluder), tabify.KeyFormatter(formatter))
	// assert.NoError(t, err, "map table writer")
	// mmt, err := json.Marshal(mt)
	// if err != nil {
	// 	t.Fatal("slice json marshal")
	// }
	// mtw := readJSON(t, "./tests/"+filename+"_expected.json")
	// assert.JSONEq(t, mtw, string(mmt[:]), "map table writer")

	// // JSONTableWriter
	// jt, err := tabify.JSON(src, tabify.KeyExcluder(excluder), tabify.KeyFormatter(formatter))
	// assert.NoError(t, err, "json table writer")
	// jtw := readJSON(t, "./tests/"+filename+"_expected.json")
	// jo := jsonmap.New()
	// jo.Set("", jt)
	// assert.JSONEq(t, jtw, jo.Stringify(), "json table writer")

}

// readJSON to read a json file and store to a jsonmap
func readJSON(t *testing.T, filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("unable to read %s", filename)
	}
	return string(data[:])
}

func TestTabify(t *testing.T) {
	json := `{
		"took": 5,
		"timed_out": false,
		"_shards": {
			"total": 5,
			"successful": 5,
			"skipped": 0,
			"failed": 0
		},
		"hits": {
			"total": 8006,
			"max_score": null,
			"hits": [
				{
					"_index": "shakespeare",
					"_type": "doc",
					"_id": "36578",
					"_score": null,
					"_source": {
						"play_name": "Hamlet",
						"line_id": 36579
					},
					"fields": {
						"expr_2": [
							"ALL"
						],
						"expr_1": [
							"all"
						]
					},
					"sort": [
						"all",
						"Hamlet"
					]
				},
				{
					"_index": "shakespeare",
					"_type": "doc",
					"_id": "34655",
					"_score": null,
					"_source": {
						"play_name": "Hamlet",
						"line_id": 34656
					},
					"fields": {
						"expr_2": [
							"ALL"
						],
						"expr_1": [
							"all"
						]
					},
					"sort": [
						"all",
						"Hamlet"
					]
				},
				{
					"_index": "shakespeare",
					"_type": "doc",
					"_id": "72001",
					"_score": null,
					"_source": {
						"play_name": "Othello",
						"line_id": 72002
					},
					"fields": {
						"expr_2": [
							"BENEDICK"
						],
						"expr_1": [
							"benedick"
						]
					},
					"sort": [
						"benedick",
						"Othello"
					]
				},
				{
					"_index": "shakespeare",
					"_type": "doc",
					"_id": "72000",
					"_score": null,
					"_source": {
						"play_name": "Othello",
						"line_id": 72001
					},
					"fields": {
						"expr_2": [
							"BENEDICK"
						],
						"expr_1": [
							"benedick"
						]
					},
					"sort": [
						"benedick",
						"Othello"
					]
				},
				{
					"_index": "shakespeare",
					"_type": "doc",
					"_id": "72002",
					"_score": null,
					"_source": {
						"play_name": "Othello",
						"line_id": 72003
					},
					"fields": {
						"expr_2": [
							"BENEDICK"
						],
						"expr_1": [
							"benedick"
						]
					},
					"sort": [
						"benedick",
						"Othello"
					]
				}
			]
		}
	}`

	hits := jsonmap.FromString(json).Get("hits.hits")
	assert.False(t, jsonmap.IsNil(hits))
	assert.Greater(t, len(hits.AsArray()), 0)

	tabifyDocExcluder := func(keys []string) bool {
		last := keys[len(keys)-1]
		return last == "_score" || last == "sort"
	}

	slice, err := tabify.Map(hits, tabify.KeyExcluder(tabifyDocExcluder))
	assert.NoError(t, err)

	for _, row := range slice {
		fmt.Println(row)
	}

}

func TestTabifyNested(t *testing.T) {
	json := jsonmap.FromString(`{
		"took": 1,
		"timed_out": false,
		"_shards": {
			"total": 5,
			"successful": 5,
			"skipped": 0,
			"failed": 0
		},
		"hits": {
			"total": 5,
			"max_score": 0,
			"hits": []
		},
		"aggregations": {
			"terms_client": {
				"doc_count_error_upper_bound": 0,
				"sum_other_doc_count": 0,
				"buckets": [
					{
						"key": "Paul",
						"doc_count": 2,
						"nested_achats": {
							"doc_count": 7,
							"sum_achats_qte": {
								"value": 15
							}
						}
					},
					{
						"key": "Jacques",
						"doc_count": 1,
						"nested_achats": {
							"doc_count": 4,
							"sum_achats_qte": {
								"value": 9
							}
						}
					},
					{
						"key": "Marie",
						"doc_count": 1,
						"nested_achats": {
							"doc_count": 2,
							"sum_achats_qte": {
								"value": 9
							}
						}
					},
					{
						"key": "Pierre",
						"doc_count": 1,
						"nested_achats": {
							"doc_count": 3,
							"sum_achats_qte": {
								"value": 5
							}
						}
					}
				]
			}
		}
	}`)

	slice, err := tabify.Map(json.Get("aggregations"), tabify.KeyExcluder(func(keys []string) bool {
		last := keys[len(keys)-1]
		return last == "doc_count_error_upper_bound" || last == "sum_other_doc_count"
	}))
	assert.NoError(t, err)

	for _, row := range slice {
		fmt.Println(row)
	}

}
