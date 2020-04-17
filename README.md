
# jsonmap

Jsonmap is a Go package to parse and use raw JSON with silent as Javascript language.
Jsonmap was inspired by [go-simplejson](https://github.com/bitly/go-simplejson) and the [lodash](https://lodash.com/) javascript library.


## Installation
```
go get github.com/datasweet/jsonmap
```

## Usage

### Parsing a new json
```
import (
  "fmt"
  "github.com/datasweet/jsonmap"
)

func main() {
  j := jsonmap.FromString(`{ "the": { "best": { "pi": 3.14 } } }`)
  i := j.Get("the.best.pi").AsInt()
  fmt.Println(i)
}
```

### Lodash utilities
You can use some lodash function utilities : 
* Filter
* Map
* Assign
* etc.

### Tabify
Tabify was created to flatten a json into tabular datas. We created this functionality to flatten json response from elasticsearch.
