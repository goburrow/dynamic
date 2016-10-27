package dynamic_test

import (
	"encoding/json"
	"fmt"

	"github.com/goburrow/dynamic"
)

type Gopher struct {
	G string
}

type Breaver struct {
	B string
}

func init() {
	dynamic.Register("go", func() interface{} { return &Gopher{} })
	dynamic.Register("br", func() interface{} { return &Breaver{} })
}

func ExampleType_unmarshal() {
	var test dynamic.Type
	data := `
{
  "type": "go",
  "G": "gopher"
}`

	if err := json.Unmarshal([]byte(data), &test); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%#v\n", test.Value())
	// Output:
	// &dynamic_test.Gopher{G:"gopher"}
}

func ExampleType_unmarshalStruct() {
	var test struct {
		X dynamic.Type
	}
	data := `
{
  "X": {
    "type": "go",
    "G": "gopher"
  }
}`
	if err := json.Unmarshal([]byte(data), &test); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%#v\n", test.X.Value())
	// Output:
	// &dynamic_test.Gopher{G:"gopher"}
}

func ExampleType_unmarshalList() {
	var test []dynamic.Type
	data := `
[
  {
    "type": "go",
    "G": "gopher"
  },
  {
    "Type": "br",
    "B": "breaver"
  }
]`

	if err := json.Unmarshal([]byte(data), &test); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	for _, t := range test {
		fmt.Printf("%#v\n", t.Value())
	}
	// Output:
	// &dynamic_test.Gopher{G:"gopher"}
	// &dynamic_test.Breaver{B:"breaver"}
}

func ExampleType_unmarshalListInStruct() {
	var test struct {
		X []dynamic.Type
	}
	data := `
{
  "X": [
    {
      "Type": "go",
      "G": "gopher"
    },
    {
      "type": "br",
      "B": "breaver"
    }
  ]
}`
	if err := json.Unmarshal([]byte(data), &test); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	for _, t := range test.X {
		fmt.Printf("%#v\n", t.Value())
	}
	// Output:
	// &dynamic_test.Gopher{G:"gopher"}
	// &dynamic_test.Breaver{B:"breaver"}
}

func ExampleType_marshal() {
	var test dynamic.Type

	test.SetValue(&Gopher{"gopher"})
	data, err := json.Marshal(&test)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%s\n", data)
	// Output:
	// {"G":"gopher"}
}

func ExampleType_marshalStruct() {
	var test struct {
		X dynamic.Type
		Y dynamic.Type
	}

	test.X.SetValue(&Gopher{"gopher"})
	test.Y.SetValue(&Breaver{"breaver"})
	data, err := json.Marshal(&test)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%s\n", data)
	// Output:
	// {"X":{"G":"gopher"},"Y":{"B":"breaver"}}
}

func ExampleType_marshalListInStruct() {
	var test struct {
		X []dynamic.Type
	}

	test.X = []dynamic.Type{
		dynamic.Type{&Gopher{"gopher"}},
		dynamic.Type{&Breaver{"breaver"}},
	}

	data, err := json.Marshal(&test)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%s\n", data)
	// Output:
	// {"X":[{"G":"gopher"},{"B":"breaver"}]}
}
