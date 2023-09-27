// Package godantic is written by Ismael GraHms
package godantic

import (
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
	switch {
	case isPtr(v):
		return g.inspect(v.Elem().Interface(), tree)
	case isStruct(v):
		return g.checkStruct(v, tree)
	case isString(v):
		return g.checkString(v, tree)
	case isTime(v):
		return g.checkTime(v, tree)
	case isList(v):
		return g.checkList(v, tree)
	default:
		return nil

	}
}

func (g *Validate) checkTime(v reflect.Value, tree string) error {
	timeValue := v.Interface().(time.Time)
	if timeValue.IsZero() {
		return &Error{
			ErrType: "INVALID_TIME_ERR",
			Path:    tree,
			Message: "The field <" + tree + "> cannot have an invalid time value",
		}
	}
	return nil
}

func (g *Validate) checkString(v reflect.Value, tree string) error {
	s := strings.TrimSpace(v.String())
	if len(s) < 1 {
		return &Error{
			ErrType: "EMPTY_STRING_ERR",
			Path:    tree,
			Message: "The field <" + tree + "> cannot be an empty string",
		}
	}
	return nil
}

func (g *Validate) checkList(v reflect.Value, tree string) error {
	min := 1
	if g.IgnoreMinLen == true {
		min = 0
	}
	isLesserThanMinLength := v.Len() < min
	if isLesserThanMinLength {

		return &Error{
			ErrType: "EMPTY_LIST_ERR",
			Path:    tree,
			Message: fmt.Sprintf("Field <%s> must be with at least one value.", tree),
		}
	}
	for i := 0; i < v.Len(); i++ {
		err := g.inspect(v.Index(i).Interface(), tree)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Validate) checkStruct(v reflect.Value, tree string) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if isTime(v.Field(i)) {
			// lets ignore time.Time fields, they are already checked in bindJSON
			continue
		}
		err := g.checkField(v, t, tree, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Validate) checkField(v reflect.Value, t reflect.Type, tree string, i int) error {
	f := t.Field(i)

	enums := f.Tag.Get("enum")
	if len(enums) == 0 {
		f.Tag.Get("enums")
	}
	valField := v.Field(i)
	switch {
	case (f.Type.Kind() == reflect.Struct || f.Type.Kind() == reflect.Ptr) && !valField.IsNil():
		err := g.inspect(valField.Interface(), fieldName(f, tree))
		if err != nil {
			return err
		}
	case g.IgnoreRequired != true:
		if isFieldRequired(f) && valField.IsNil() {
			return RequiredFieldError(f, tree)

		}
	case isPtr(v):
		if !v.IsValid() || v.IsNil() {
			return nil // nil pointer is valid
		}
		return g.inspect(v.Elem().Interface(), tree)
	}
	// Check for enum validation tags.
	if len(enums) > 0 {
		enumValues := strings.Split(strings.TrimSpace(enums), ",")
		err := g.strEnums(f, valField, tree, enumValues)

		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Validate) strEnums(f reflect.StructField, val reflect.Value, tree string, allowedValues []string) error {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	fieldValue := val.String()

	// Convert the allowedValues slice to a map for faster searching.
	allowedMap := make(map[string]bool, len(allowedValues))
	for _, allowedValue := range allowedValues {
		allowedMap[allowedValue] = true
	}

	if _, ok := allowedMap[fieldValue]; !ok {
		return &Error{
			ErrType: "INVALID_ENUM_ERR",
			Path:    fieldName(f, tree),
			Message: fmt.Sprintf("The field <%s> must have one of the following values: %s, '%s' was given",
				fieldName(f, tree), strings.Join(allowedValues, ", "), val.String()),
		}
	}

	return nil
}

func fieldName(f reflect.StructField, tree string) string {
	var name string
	if f.Tag != "" {
		jsonTag := strings.Split(f.Tag.Get("json"), ",")[0]
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
