# dynamic [![Build Status](https://travis-ci.org/goburrow/dynamic.svg?branch=master)](https://travis-ci.org/goburrow/dynamic) [![GoDoc](https://godoc.org/github.com/goburrow/dynamic?status.svg)](https://godoc.org/github.com/goburrow/dynamic)

Unmarshal JSON with dynamic (multiple) types.

## Examples
See [godoc](https://godoc.org/github.com/goburrow/dynamic#pkg-examples)

```
go get github.com/goburrow/dynamic
```
```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/goburrow/dynamic"
)

type Info struct {
	Type    string
	Message string
}

type Error struct {
	Errors []string
}

type Status struct {
	Code int
}

func init() {
	dynamic.Register("info", func() interface{} { return &Info{} })
	dynamic.Register("error", func() interface{} { return &Error{} })
	dynamic.Register("status", func() interface{} { return &Status{} })
}

func main() {
	json1 := `{"Type": "info", "Message": "hello"}`
	json2 := `{"Type": "error", "Errors": ["unknown"]}`
	json3 := `{"Type": "status", "Data": {"Code": 213}}`

	var obj1, obj2 dynamic.Type
	var obj3 dynamic.Data
	json.Unmarshal([]byte(json1), &obj1)
	json.Unmarshal([]byte(json2), &obj2)
	json.Unmarshal([]byte(json3), &obj3)

	fmt.Printf("1: %#v\n", obj1.Value())
	fmt.Printf("2: %#v\n", obj2.Value())
	fmt.Printf("3: %#v\n", obj3.Value())
}
```
Output:
```
1: &main.Info{Type:"info", Message:"hello"}
2: &main.Error{Errors:[]string{"unknown"}}
3: &main.Status{Code:213}
```

## TODO
- Support sub-types.
