package polytype

import (
	"encoding/json"
	"fmt"
	"testing"
)

type x struct {
	X string
}

type y struct {
	Y int
}

func createX() interface{} {
	return &x{}
}

func createY() interface{} {
	return &y{}
}

func init() {
	AddType("x", createX)
	AddType("y", createY)
}

func BenchmarkUnmarshal(b *testing.B) {
	data := `
{
  "object": {
    "type": "x",
    "x": "x"
  },
  "list": [
    {
      "type": "x",
      "x": "0"
    },
    {
      "type": "y",
      "y": 1
    },
    {
      "type": "x",
      "x": "0",
      "y": 1
    },
    {
      "type": "y",
      "x": "0",
      "y": 1
    }
  ]
}
`
	var test struct {
		Object Polytype
		List   []Polytype
	}
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal([]byte(data), &test)
		if err != nil {
			fmt.Printf("%v\n%#v\n", err, test)
			break
		}
	}
}
