package dynamic

import (
	"encoding/json"
	"testing"
)

var _ (json.Marshaler) = (*Type)(nil)
var _ (json.Unmarshaler) = (*Type)(nil)
var _ (json.Unmarshaler) = (*Data)(nil)

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
