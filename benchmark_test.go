package dynamic

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
	Register("x", createX)
	Register("y", createY)
}

func BenchmarkUnmarshal(b *testing.B) {
	data := []byte(`
{
  "object": {
    "type": "x",
    "x": "x",
    "foo": "bar"
  },
  "list": [
    {
      "type": "x",
      "x": "0",
      "foo": ["bar", "bas", "baz"]
    },
    {
      "type": "y",
      "y": 1,
      "foo": {"bar": "bar", "baz": "baz"}
    },
    {
      "type": "x",
      "x": "0",
      "y": 1,
      "foo": [1, 2, 3, 4, 5]
    },
    {
      "type": "y",
      "x": "0",
      "y": 1,
      "foo": {"bar": ["b", "ba", "ar"]}
    }
  ]
}
`)
	var test struct {
		Object Type
		List   []Type
	}
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, &test)
		if err != nil {
			fmt.Printf("%v\n%#v\n", err, test)
			break
		}
	}
}
