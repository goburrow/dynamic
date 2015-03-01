/*
Package polytype allows unmarshal JSON property with multiple types.
*/
package polytype

import (
	"encoding/json"
	"errors"
)

// Supported types
var types = make(map[string]func() interface{})

// Add associates factory method for a type name.
// The factory method must return a pointer to a struct it's going to create.
func Add(name string, f func() interface{}) {
	if _, existed := types[name]; existed {
		panic("polytype: type \"" + name + "\" has already been added")
	}

	types[name] = f
}

// Type uses interface{} as the underlying data structure.
type Type struct {
	// Value is made public so the inner structs can be validated.
	Value interface{}
}

var _ (json.Marshaler) = (*Type)(nil)
var _ (json.Unmarshaler) = (*Type)(nil)

// MarshalJSON marshals the polyType.Value
func (t *Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

// UnmarshalJSON first read Type property in JSON object, then unmarshal JSON
// to the instance created by respective factory method.
func (t *Type) UnmarshalJSON(data []byte) error {
	var typed struct {
		Type string
	}
	if err := json.Unmarshal(data, &typed); err != nil {
		return err
	}
	if typed.Type == "" {
		return errors.New("type must be specified")
	}

	f, ok := types[typed.Type]
	if !ok {
		return errors.New("type \"" + typed.Type + "\" is not supported")
	}
	value := f()
	if err := json.Unmarshal(data, value); err != nil {
		return err
	}
	t.Value = value
	return nil
}
