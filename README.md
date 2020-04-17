
# jsonmap
[![Circle CI](https://circleci.com/gh/datasweet/jsonmap.svg?style=svg)](https://circleci.com/gh/datasweet/jsonmap) [![Go Report Card](https://goreportcard.com/badge/github.com/datasweet/jsonmap)](https://goreportcard.com/report/github.com/datasweet/jsonmap) [![GoDoc](https://godoc.org/github.com/datasweet/jsonmap?status.png)](https://godoc.org/github.com/datasweet/jsonmap) [![GitHub stars](https://img.shields.io/github/stars/datasweet/jsonmap.svg)](https://github.com/datasweet/jsonmap/stargazers)
[![GitHub license](https://img.shields.io/github/license/datasweet/jsonmap.svg)](https://github.com/datasweet/jsonmap/blob/master/LICENSE)

[![datasweet-logo](https://www.datasweet.fr/wp-content/uploads/2019/02/datasweet-black.png)](http://www.datasweet.fr)

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

## Who are we ?
We are Datasweet, a french startup providing full service (big) data solutions.

## Questions ? problems ? suggestions ?
If you find a bug or want to request a feature, please create a [GitHub Issue](https://github.com/datasweet/jsonmap/issues/new).

## License
```
This software is licensed under the Apache License, version 2 ("ALv2"), quoted below.

Copyright 2017-2020 Datasweet <http://www.datasweet.fr>

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy of
the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations under
the License.
```
