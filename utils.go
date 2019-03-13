package jsonmap

import "strings"

func createPath(path string) []string {
	var keys []string
	infos := strings.Split(path, ".")
	var tmp strings.Builder
	for _, k := range infos {
		if len(k) > 0 && k[len(k)-1] == '\\' {
			tmp.WriteString(k[:len(k)-1])
			tmp.WriteString(".")
			continue
		} else {
			tmp.WriteString(k)
		}
		parts := strings.Split(tmp.String(), "[")
		for _, p := range parts {
			t := strings.TrimSpace(strings.TrimSuffix(p, "]"))
			if len(t) > 0 {
				keys = append(keys, t)
			}
		}
		tmp.Reset()
	}
	return keys
}
