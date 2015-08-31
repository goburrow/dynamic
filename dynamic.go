/*
Package dynamic provides support for unmarshal dynamic JSON objects.
*/
package dynamic

import (
	"encoding/json"
	"errors"
)

// Supported types
var types = make(map[string]func() interface{})

// Register associates factory method for a type name.
// The factory method must return a pointer to a struct it's going to create.
func Register(name string, f func() interface{}) {
	if _, existed := types[name]; existed {
		panic("polytype: type \"" + name + "\" has already been added")
	}

	types[name] = f
}

// Type represents objects that have their properties at top level along with
// `Type' property. The `Type' property is not compulsory when unmarshalling
// but needed when marshalling to JSON. For example:
//
// 	{
// 		"Type": "Point",
//	 	"X": 21,
// 		"Y": 3
// 	}
//
// is the JSON representation of:
//
// 	type Point struct {
//		X int
// 		Y int
// 	}
//
type Type [1]interface{}

// MarshalJSON marshals the t.Value.
func (t *Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t[0])
}

// UnmarshalJSON first reads Type property in JSON object, then unmarshals JSON
// to the instance created by respective factory method.
func (t *Type) UnmarshalJSON(data []byte) error {
	var envelope struct {
		Type string
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return err
	}
	if envelope.Type == "" {
		return errors.New("type must be specified")
	}

	f, ok := types[envelope.Type]
	if !ok {
		return errors.New("type \"" + envelope.Type + "\" is not supported")
	}
	value := f()
	if err := json.Unmarshal(data, value); err != nil {
		return err
	}
	t[0] = value
	return nil
}

// Value returns the value marshalled in t.
func (t *Type) Value() interface{} {
	return t[0]
}

// SetValue sets value to t.
func (t *Type) SetValue(v interface{}) {
	t[0] = v
}

// Data represents JSON objects that have a `Type' and `Data' property for
// underlining data structure. For example:
//
// 	{
// 		"Type": "Point",
//	 	"Data": {
// 			"X": 21,
// 			"Y": 3
// 		}
// 	}
//
// is the JSON representation of:
//
// 	type Point struct {
//		X int
// 		Y int
// 	}
//
type Data struct {
	Type string
	Data interface{}
}

// UnmarshalJSON first reads Type property in JSON object, then unmarshals JSON
// to the instance created by respective factory method.
func (d *Data) UnmarshalJSON(data []byte) error {
	var envelope struct {
		Type string
		Data json.RawMessage
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return err
	}
	if envelope.Type == "" {
		return errors.New("type must be specified")
	}

	f, ok := types[envelope.Type]
	if !ok {
		return errors.New("type \"" + envelope.Type + "\" is not supported")
	}
	value := f()
	if err := json.Unmarshal(envelope.Data, value); err != nil {
		return err
	}
	d.Type = envelope.Type
	d.Data = value
	return nil
}

func (d *Data) Value() interface{} {
	return d.Data
}
