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

func (g *Validate) validateInterfaceHooks(val any, path string) *Error {
	rv := reflect.ValueOf(val)

	// ValidationPlugin hook
	if cv, ok := resolveInterface[ValidationPlugin](rv); ok {
		if err := cv.Validate(); err != nil {
			if err.Path == "" {
				err.Path = path
			}
			return &Error{
				ErrType: err.ErrType,
				Message: err.Message,
				Path:    err.Path,
				err:     err,
			}
		}
	}

	// DynamicFieldsValidator hook
	if df, ok := resolveInterface[DynamicFieldsValidator](rv); ok {
		if err := validateDynamicFields(df.GetValue(), df.GetAttribute(), df.GetValueType(), path); err != nil {
			return err
		}
	}

	return nil
}
