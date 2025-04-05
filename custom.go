package godantic

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type customValidatorFunc func(value any, path string) *Error

var (
	customValidators   = make(map[reflect.Type]map[string]customValidatorFunc)
	customValidatorMux sync.RWMutex
)

func RegisterCustom[T any](tag string, fn func(T, string) *Error) {
	customValidatorMux.Lock()
	defer customValidatorMux.Unlock()

	var zero T
	t := reflect.TypeOf(zero)

	if customValidators[t] == nil {
		customValidators[t] = make(map[string]customValidatorFunc)
	}

	customValidators[t][tag] = func(value any, path string) *Error {
		v, ok := value.(T)
		if !ok {
			return &Error{
				ErrType: "INVALID_TYPE_ERR",
				Path:    path,
				Message: fmt.Sprintf("Expected type %T but got %T", zero, value),
			}
		}
		return fn(v, path)
	}
}

func getCustomValidator(t reflect.Type, tag string) (customValidatorFunc, bool) {
	customValidatorMux.RLock()
	defer customValidatorMux.RUnlock()

	if m, ok := customValidators[t]; ok {
		fn, exists := m[tag]
		return fn, exists
	}
	return nil, false
}

func (g *Validate) validateWithCustomTag(val any, f reflect.StructField, path string) *Error {
	tag := f.Tag.Get("validate")
	if tag == "" {
		return nil
	}

	t := f.Type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	tags := strings.Split(tag, ",")
	for _, singleTag := range tags {
		singleTag = strings.TrimSpace(singleTag)
		if singleTag == "" {
			continue
		}

		if fn, ok := getCustomValidator(t, singleTag); ok {
			err := fn(val, path)
			if err != nil {
				if err.Path == "" {
					err.Path = path
				}
				return err
			}
		}
	}

	return nil
}
