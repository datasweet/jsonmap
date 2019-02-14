package jsonmap

import "strings"

func createPath(path string) []string {
	var keys []string
	infos := strings.Split(path, ".")
	for _, k := range infos {
		parts := strings.Split(k, "[")
		for _, p := range parts {
			t := strings.TrimSpace(strings.TrimSuffix(p, "]"))
			if len(t) > 0 {
				keys = append(keys, t)
			}
		}
	}
	return keys
}
