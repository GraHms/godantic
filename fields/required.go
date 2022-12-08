// Package fields  is written by Ismael GraHms
package fields

import (
	"reflect"
	"strings"
)

// ValidateStruct  takes an interface value and a tree string, and recursively
// examines the struct fields of the value using reflection. If the field is a
// struct or a pointer to a struct, the function calls itself recursively on that
// field.
func (fr *Fields) ValidateStruct(val interface{}, currentPath string) error {
	// Get the reflection of the value.
	v := reflect.ValueOf(val)

	// Get the type of the value.
	t := v.Type()

	// Check the kind of the value.
	switch v.Kind() {
	// If the value is a pointer, call the function recursively on the value
	// that the pointer points to.
	case reflect.Ptr:
		err := fr.ValidateStruct(v.Elem().Interface(), currentPath)
		if err != nil {
			return err
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			// Use the Index() method to get the value of the element at the current index.
			valElement := v.Index(i)

			// If the element is a struct or pointer to a struct and its value is not nil,
			// call the function recursively on the element.
			if (valElement.Kind() == reflect.Struct || valElement.Kind() == reflect.Ptr) && !valElement.IsNil() {
				err := fr.ValidateStruct(valElement.Interface(), currentPath)
				if err != nil {
					return err
				}
			}
		}

	// If the value is a struct, loop through its fields and call the function
	// recursively on each field if it is a struct or pointer to a struct.
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			// Get the field at the current index.
			f := t.Field(i)

			// Get the value of the field.
			valField := v.Field(i)

			// If the field is a struct or pointer to a struct and its value is not nil,
			// call the function recursively on the field.
			if (f.Type.Kind() == reflect.Struct || f.Type.Kind() == reflect.Ptr) && !valField.IsNil() {

				err := fr.ValidateStruct(valField.Interface(), fr.pathBuilder(f, currentPath))
				if err != nil {
					return err
				}
			}

			// If the field is required and its value is nil, print an error message.
			if isRequired(f) && valField.IsNil() {
				return &GodanticError{
					ErrType: "requiredFieldErr",
					Path:    fr.pathBuilder(f, currentPath),
					Message: "The field <" + fr.pathBuilder(f, currentPath) + "> is required",
				}

			}

		}
	case reflect.String:
		s := strings.TrimSpace(v.String())
		if len(s) < 1 {
			return &GodanticError{
				ErrType: "emptyStrFieldErr",
				Path:    currentPath,
				Message: "The field <" + currentPath + "> cannot be an empty string",
			}
		}
	}
	return nil
}

// pathBuilder takes a struct field and a tree string, and returns the name of
// the field. If the field has a JSON tag, the name is the value of the tag.
// Otherwise, the name is the field name. If the tree string is not empty,
// the name is the tree string concatenated with the field name.
func (fr *Fields) pathBuilder(f reflect.StructField, fieldPath string) string {
	var name string
	if f.Tag != "" {
		jsonTag := f.Tag.Get("json")
		name = jsonTag
		if len(fieldPath) > 0 {
			name = fieldPath + fr.ErrorPathSeparator + jsonTag
		}
	}
	return name
}

// isRequired takes a struct field and returns true if the field is required
// and false otherwise. A field is considered required if it has a "binding"
// tag with the value "required".
func isRequired(f reflect.StructField) bool {
	if f.Tag.Get("binding") == "required" {
		return true
	}
	return false
}
