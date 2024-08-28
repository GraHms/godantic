package godantic

import "fmt"

type ValidationPlugin interface {
	Validate() *CustomErr
}

type DynamicFieldsValidator interface {
	GetValue() any
	GetValueType() string
	GetAttribute() string
}

func validateDynamicFields(fieldValue interface{}, attr, valueType, tree string) *Error {
	// Build the full path to the field in the data structure
	fullPath := fmt.Sprintf("%s.%s", tree, attr)

	// Perform validation based on the value type
	switch valueType {
	case "numeric", "number", "integer", "int":
		// Check if the value is an integer
		_, ok := fieldValue.(int)
		if !ok {
			// If not, check if it can be converted to an integer
			floatVal, fok := fieldValue.(float64)
			if !fok || floatVal != float64(int(floatVal)) {
				return &Error{
					ErrType: "INVALID_VALUE_TYPE_ERR",
					Path:    fullPath,
					Message: fmt.Sprintf("Invalid value type for field '%s' at path '%s'. Expected numeric value.", attr, fullPath),
				}
			}
		}
	case "string":
		if _, ok := fieldValue.(string); !ok {
			return &Error{
				ErrType: "INVALID_VALUE_TYPE_ERR",
				Path:    fullPath,
				Message: fmt.Sprintf("Invalid value type for field '%s' at path '%s'. Expected string value.", attr, fullPath),
			}
		}
	case "float":
		if _, ok := fieldValue.(float64); !ok {
			return &Error{
				ErrType: "INVALID_VALUE_TYPE_ERR",
				Path:    fullPath,
				Message: fmt.Sprintf("Invalid value type for field '%s' at path '%s'. Expected float value.", attr, fullPath),
			}
		}
	case "boolean":
		if _, ok := fieldValue.(bool); !ok {
			return &Error{
				ErrType: "INVALID_VALUE_TYPE_ERR",
				Path:    fullPath,
				Message: fmt.Sprintf("Invalid value type for field '%s' at path '%s'. Expected boolean value.", attr, fullPath),
			}
		}
	default:
		return &Error{
			ErrType: "INVALID_VALUE_TYPE_ERR",
			Path:    fullPath,
			Message: fmt.Sprintf("Invalid value type '%s' for field '%s' at path '%s'.", valueType, attr, fullPath),
		}
	}

	return nil
}
