package godantic

import "reflect"

func resolveInterface[T any](val reflect.Value) (T, bool) {
	var zero T
	if v, ok := val.Interface().(T); ok {
		return v, true
	}

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		if v, ok := val.Elem().Interface().(T); ok {
			return v, true
		}
	}

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		elem := val.Elem()
		if inst, ok := elem.Interface().(T); ok {
			return inst, true
		}
	}

	return zero, false
}
