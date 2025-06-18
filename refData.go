package godantic

import (
	"reflect"
	"strings"
)

func buildRefData(v any) map[string]any {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	result := make(map[string]any)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		fieldType := field.Type

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		fieldName := strings.Split(jsonTag, ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		switch {
		case fieldType == reflect.TypeOf(Object{}):
			result[fieldName] = Object{}

		case fieldType.Kind() == reflect.Struct:
			result[fieldName] = buildRefData(fieldVal.Interface())

		case fieldType.Kind() == reflect.Slice:
			slice := []any{}
			if fieldType.Elem().Kind() == reflect.Struct {
				slice = append(slice, buildRefData(reflect.New(fieldType.Elem()).Interface()))
			}
			result[fieldName] = slice

		case fieldType.Kind() == reflect.Map:
			result[fieldName] = map[string]any{}

		default:
			// Tipos primitivos ou interface{}
			result[fieldName] = reflect.Zero(fieldType).Interface()
		}
	}

	return result
}
