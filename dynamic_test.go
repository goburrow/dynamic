package dynamic

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
	Register("a", createA)
	Register("b", createB)
}

func TestNameDuplicated(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("panic expected")
		}
	}()

	Register("a", createB)
}

func ExampleType_unmarshal() {
	var test Type
	data := `
{
  "type": "a",
  "A": "This is A"
}`

	if err := json.Unmarshal([]byte(data), &test); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%#v\n", test.Value())
	// Output:
	// &dynamic.a{A:"This is A"}
}

func ExampleType_unmarshalStruct() {
	var test struct {
		X Type
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
	fmt.Printf("%#v\n", test.X.Value())
	// Output:
	// &dynamic.a{A:"This is A"}
}

func ExampleType_unmarshalList() {
	var test []Type
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
		fmt.Printf("%#v\n", t.Value())
	}
	// Output:
	// &dynamic.a{A:"This is A"}
	// &dynamic.b{B:"This is B"}
}

func ExampleType_unmarshalListInStruct() {
	var test struct {
		X []Type
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
		fmt.Printf("%#v\n", t.Value())
	}
	// Output:
	// &dynamic.a{A:"This is A"}
	// &dynamic.b{B:"This is B"}
}

func TestUnmarshalUnsupportedType(t *testing.T) {
	var test Type
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
	var test Type
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
		Type
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
		switch i.Value().(type) {
		case *a:
		default:
			t.Fatalf("%#v is not *a", i.Value())
		}
	}
}

func ExampleType_marshal() {
	var test Type

	test.SetValue(&a{"This is A"})
	data, err := json.Marshal(&test)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%s\n", data)
	// Output:
	// {"A":"This is A"}
}

func ExampleType_marshalStruct() {
	var test struct {
		X Type
		Y Type
	}

	test.X.SetValue(&a{"This is A"})
	test.Y.SetValue(&b{"This is B"})
	data, err := json.Marshal(&test)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%s\n", data)
	// Output:
	// {"X":{"A":"This is A"},"Y":{"B":"This is B"}}
}

func ExampleType_marshalListInStruct() {
	var test struct {
		X []Type
	}

	test.X = []Type{
		Type{&a{"This is A"}},
		Type{&b{"This is B"}},
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
