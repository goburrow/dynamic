# Polytype [![Build Status](https://travis-ci.org/goburrow/polytype.svg?branch=master)](https://travis-ci.org/goburrow/polytype) [![GoDoc](https://godoc.org/github.com/goburrow/polytype?status.svg)](https://godoc.org/github.com/goburrow/polytype)

Unmarshal JSON with multiple types.

## Examples
See [godoc](https://godoc.org/github.com/goburrow/polytype#pkg-examples)

```
go get github.com/goburrow/polytype
```
```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/goburrow/polytype"
)

type Info struct {
	Type    string
	Message string
}

type Error struct {
	Type string
	Code int
}

func init() {
	polytype.Register("info", func() interface{} { return &Info{} })
	polytype.Register("error", func() interface{} { return &Error{} })
}

func main() {
	json1 := `{"Type": "info", "Message": "hello"}`
	json2 := `{"Type": "error", "Code": -213}`
	var obj1, obj2 polytype.Type
	json.Unmarshal([]byte(json1), &obj1)
	json.Unmarshal([]byte(json2), &obj2)
	fmt.Printf("1: %#v\n2: %#v\n", obj1.Value(), obj2.Value())
}
```
Output:
```
1: &main.Info{Type:"info", Message:"hello"}
2: &main.Error{Type:"error", Code:-213}
```



## TODO
- Support sub-types.
