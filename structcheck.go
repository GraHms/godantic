// Package godantic is written by Ismael GraHms
package godantic

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func getValueOf(val interface{}) reflect.Value {
	return reflect.ValueOf(val)
}
func isPtr(value reflect.Value) bool {
	return value.Kind() == reflect.Ptr
}
func isStruct(value reflect.Value) bool {
	kind := value.Kind()
	return kind == reflect.Struct || (kind == reflect.Ptr && value.Elem().Kind() == reflect.Struct)
}

func isString(value reflect.Value) bool {
	return value.Kind() == reflect.String
}

func Validate(val interface{}) error {
	g := Godentic{}
	return g.InspectStruct(val, "")
}

func (g *Godentic) InspectStruct(val interface{}, tree string) error {
	// Get the reflection of the value.
	v := getValueOf(val)

	// Get the type of the value.
	t := v.Type()

	if isPtr(v) {
		return g.InspectStruct(v.Elem().Interface(), tree)
	}
	if isStruct(v) {
		for i := 0; i < t.NumField(); i++ {
			// Get the field at the current index.
			err := g.checkField(v, t, tree, i)
			if err != nil {
				return err
			}
		}
	}
	if isString(v) {
		s := strings.TrimSpace(v.String())
		if len(s) < 1 {
			return &GodanticError{
				ErrType: "emptyStrFieldErr",
				Path:    tree,
				Message: "The field <" + tree + "> cannot be an empty string",
			}
		}
	}

	return nil
}

func (g *Godentic) checkField(v reflect.Value, t reflect.Type, tree string, i int) error {
	f := t.Field(i)

	// Get the value of the field.
	valField := v.Field(i)

	// If the field is a struct or pointer to a struct and its value is not nil,
	// call the function recursively on the field.
	if (f.Type.Kind() == reflect.Struct || f.Type.Kind() == reflect.Ptr) && !valField.IsNil() {

		err := g.InspectStruct(valField.Interface(), fieldName(f, tree))
		if err != nil {
			return err
		}
	}

	// If the field is required and its value is nil, print an error message.
	if isFieldRequired(f) && valField.IsNil() {
		return RequiredFieldError(f, tree)

	}
	return nil
}

// fieldName takes a struct field and a tree string, and returns the name of
// the field. If the field has a JSON tag, the name is the value of the tag.
// Otherwise, the name is the field name. If the tree string is not empty,
// the name is the tree string concatenated with the field name.
func fieldName(f reflect.StructField, tree string) string {
	var name string
	if f.Tag != "" {
		jsonTag := f.Tag.Get("json")
		name = jsonTag
		if len(tree) > 0 {
			name = tree + "." + jsonTag
		}
	}
	return name
}

func isFieldRequired(f reflect.StructField) bool {
	if f.Tag.Get("binding") == "required" {
		return true
	}
	return false
}

// RequiredFieldError returns an error indicating that a required field is missing.
func RequiredFieldError(field reflect.StructField, tree string) error {
	return &GodanticError{
		ErrType: "requiredFieldErr",
		Path:    fieldName(field, tree),
		Message: fmt.Sprintf("The field <%s> is required", fieldName(field, tree)),
	}
}

func (e *GodanticError) Error() string {
	e.err = errors.New(e.Message)
	return e.err.Error()
}

type GodanticError struct {
	ErrType string
	Message string
	Path    string
	err     error
}
