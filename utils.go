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

// EscapePath to escape a path
// Example
// - By default  jsonmap.Set("message.raw", "hello world !")
//   =>  { "message": { "raw": "hello world !" }}
// - With escape jsonmap.Set(jsonmap.EscapePath("message.raw", "hello world !"))
//   => { "message.raw": "hello world !" }}
func EscapePath(path string) string {
	return strings.Replace(path, ".", "\\.", -1)
}
