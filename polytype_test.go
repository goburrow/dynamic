package polytype

import (
	"encoding/json"
	"fmt"
	"testing"
)

type a struct {
	A string
}

type b struct {
	B string
}

func createA() interface{} {
	return &a{}
}

func createB() interface{} {
	return &b{}
}

func init() {
	AddType("a", createA)
	AddType("b", createB)
}

func ExamplePolytype_unmarshal() {
	var test Polytype
	data := `
{
  "type": "a",
  "A": "This is A"
}`

	if err := json.Unmarshal([]byte(data), &test); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%#v\n", test.Value)
	// Output:
	// &polytype.a{A:"This is A"}
}

func ExamplePolytype_unmarshalStruct() {
	var test struct {
		X Polytype
	}
	data := `
{
  "X": {
    "type": "a",
    "A": "This is A"
  }
}`
	if err := json.Unmarshal([]byte(data), &test); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%#v\n", test.X.Value)
	// Output:
	// &polytype.a{A:"This is A"}
}

func ExamplePolytype_unmarshalList() {
	var test []Polytype
	data := `
[
  {
    "type": "a",
    "A": "This is A",
    "C": "This is C"
  },
  {
    "Type": "b",
    "A": "This is A",
    "B": "This is B"
  }
]`

	if err := json.Unmarshal([]byte(data), &test); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	for _, t := range test {
		fmt.Printf("%#v\n", t.Value)
	}
	// Output:
	// &polytype.a{A:"This is A"}
	// &polytype.b{B:"This is B"}
}

func ExamplePolytype_unmarshalListInStruct() {
	var test struct {
		X []Polytype
	}
	data := `
{
  "X": [
    {
      "type": "a",
      "A": "This is A"
    },
    {
      "Type": "b",
      "A": "This is A",
      "B": "This is B"
    }
  ]
}`
	if err := json.Unmarshal([]byte(data), &test); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	for _, t := range test.X {
		fmt.Printf("%#v\n", t.Value)
	}
	// Output:
	// &polytype.a{A:"This is A"}
	// &polytype.b{B:"This is B"}
}

func TestUnmarshalUnsupportedType(t *testing.T) {
	var test Polytype
	data := `
{
  "type": "c",
  "C": "This is C"
}`

	err := json.Unmarshal([]byte(data), &test)
	if err == nil {
		t.Fatal("error expected")
	}
	if "type \"c\" is not supported" != err.Error() {
		t.Fatal(err)
	}
}

func TestUnmarshalEmptyType(t *testing.T) {
	var test Polytype
	data := `
{
  "type": "",
  "A": "This is A"
}`

	err := json.Unmarshal([]byte(data), &test)
	if err == nil {
		t.Fatal("error expected")
	}
	if "type must be specified" != err.Error() {
		t.Fatal(err)
	}
}

func TestUnmarshalEmbeddedStruct(t *testing.T) {
	type inner struct {
		Polytype
	}
	type iinner struct {
		inner
	}
	var test struct {
		X []iinner
	}
	data := `
{
  "X": [
    {
      "type": "a",
      "A": "This is A"
    },
    {
      "type": "a",
      "A": "This is A"
    }
  ]
}`

	if err := json.Unmarshal([]byte(data), &test); err != nil {
		t.Fatal(err)
	}
	for _, i := range test.X {
		switch i.Value.(type) {
		case *a:
		default:
			t.Fatalf("%#v is not *a", i.Value)
		}
	}
}

func ExamplePolytype_marshal() {
	var test Polytype

	test.Value = &a{"This is A"}
	data, err := json.Marshal(&test)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%s\n", data)
	// Output:
	// {"A":"This is A"}
}

func ExamplePolytype_marshalStruct() {
	var test struct {
		X Polytype
		Y Polytype
	}

	test.X.Value = &a{"This is A"}
	test.Y.Value = &b{"This is B"}
	data, err := json.Marshal(&test)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%s\n", data)
	// Output:
	// {"X":{"A":"This is A"},"Y":{"B":"This is B"}}
}

func ExamplePolytype_marshalListInStruct() {
	var test struct {
		X []Polytype
	}

	test.X = []Polytype{
		Polytype{Value: &a{"This is A"}},
		Polytype{Value: &b{"This is B"}},
	}
	data, err := json.Marshal(&test)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%s\n", data)
	// Output:
	// {"X":[{"A":"This is A"},{"B":"This is B"}]}
}
