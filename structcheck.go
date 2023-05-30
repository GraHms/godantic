// Package godantic is written by Ismael GraHms
package godantic

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
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
func isTime(value reflect.Value) bool {
	return value.Type().ConvertibleTo(TimeType)
}
func isList(value reflect.Value) bool {
	return value.Kind() == reflect.Slice || value.Kind() == reflect.Array
}

var TimeType = reflect.TypeOf(time.Time{})

func (g *Validate) InspectStruct(val interface{}) error {
	return g.inspect(val, "")
}

func (g *Validate) inspect(val interface{}, tree string) error {
	v := getValueOf(val)

	t := v.Type()

	if isPtr(v) {
		return g.inspect(v.Elem().Interface(), tree)
	}
	if isStruct(v) {
		for i := 0; i < t.NumField(); i++ {
			err := g.checkField(v, t, tree, i)
			if err != nil {
				return err
			}
		}
	}
	if isString(v) {
		s := strings.TrimSpace(v.String())
		if len(s) < 1 {
			return &Error{
				ErrType: "EMPTY_STRING_ERR",
				Path:    tree,
				Message: "The field <" + tree + "> cannot be an empty string",
			}
		}
	}
	if isTime(v) {
		timeValue := v.Interface().(time.Time)
		if timeValue.IsZero() {
			return &Error{
				ErrType: "INVALID_TIME_ERR",
				Path:    tree,
				Message: "The field <" + tree + "> cannot have an invalid time value",
			}
		}
	}
	if isList(v) {

		isLesserThanMinLength := v.Len() < 1

		if isLesserThanMinLength {

			return &Error{
				ErrType: "EMPTY_LIST_ERR",
				Path:    tree,
				Message: fmt.Sprintf("Field <%s> must be with at least one value.", tree),
			}
		}
	}

	return nil
}
func (g *Validate) checkField(v reflect.Value, t reflect.Type, tree string, i int) error {
	f := t.Field(i)

	valField := v.Field(i)

	if (f.Type.Kind() == reflect.Struct || f.Type.Kind() == reflect.Ptr) && !valField.IsNil() {

		err := g.inspect(valField.Interface(), fieldName(f, tree))
		if err != nil {
			return err
		}
	}

	if isFieldRequired(f) && valField.IsNil() {
		return RequiredFieldError(f, tree)

	}
	return nil
}

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

func RequiredFieldError(field reflect.StructField, tree string) error {
	return &Error{
		ErrType: "REQUIRED_FIELD_ERR",
		Path:    fieldName(field, tree),
		Message: fmt.Sprintf("The field <%s> is required", fieldName(field, tree)),
	}
}

func (e *Error) Error() string {
	e.err = errors.New(e.Message)
	return e.err.Error()
}

type Error struct {
	ErrType string
	Message string
	Path    string
	err     error
}
