package godantic

import (
	"fmt"
	"reflect"
	"strings"
)

// extractConditionalFields scans the struct and stores "when" conditions in a map
func extractEnumValues(rootVal reflect.Value, parentPath string) map[string]string {
	enumMap := make(map[string]string)
	queue := []struct {
		val  reflect.Value
		path string
	}{{rootVal, parentPath}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		v := current.val
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			continue
		}

		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)

			// Convert JSON tag to JSON path
			jsonKey := field.Tag.Get("json")
			if jsonKey == "" {
				jsonKey = field.Name // Fallback to field name
			}

			fullPath := jsonKey
			if current.path != "" {
				fullPath = current.path + "." + jsonKey // Construct full path
			}

			// Extract inputted enum value
			if _, ok := field.Tag.Lookup("enum"); ok {
				if fieldValue.Kind() == reflect.Ptr {
					if fieldValue.IsNil() {
						continue // Skip nil pointers
					}
					fieldValue = fieldValue.Elem()
				}

				if fieldValue.Kind() == reflect.String {
					enumMap[fullPath] = fieldValue.String()
				}
			}

			// Add nested structs to the queue
			if fieldValue.Kind() == reflect.Struct || (fieldValue.Kind() == reflect.Ptr && fieldValue.Elem().Kind() == reflect.Struct) {
				queue = append(queue, struct {
					val  reflect.Value
					path string
				}{fieldValue, fullPath})
			}
		}
	}
	return enumMap
}

// parseCondition extracts conditions and bindings from the "when" tag.
// parseCondition extracts conditions and bindings from the "when" tag.
func parseCondition(conditionTag string) (map[string]string, map[string]string) {
	conditions := make(map[string]string)
	bindings := make(map[string]string)

	if conditionTag == "" {
		return conditions, bindings // No condition present
	}

	// Split multiple conditions separated by ";"
	parts := strings.Split(conditionTag, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue // Ignore empty segments
		}

		// Handle "binding=" separately
		if strings.HasPrefix(part, "binding=") {
			bindings["binding"] = strings.TrimPrefix(part, "binding=")
			continue
		}

		// Find the first occurrence of an operator (=, >, <, >=, <=, !=)
		operatorIndex := strings.IndexAny(part, "=><!")
		if operatorIndex == -1 {
			continue // Invalid condition format
		}

		// Extract key, operator, and value separately
		key := strings.TrimSpace(part[:operatorIndex])
		//operator := string(part[operatorIndex]) // Capture operator
		value := strings.TrimSpace(part[operatorIndex+1:])

		// Store condition with the correct value
		conditions[key] = value // No more `=` prefix
	}

	return conditions, bindings
}

// validateCondition checks if a field's condition is met and applies validation rules accordingly.
func (g *Validate) validateCondition(f reflect.StructField, valField reflect.Value, fullPath string, enumMap map[string]string) error {
	// Get the "when" tag
	conditionTag, hasCondition := f.Tag.Lookup("when")
	if !hasCondition {
		return nil // No condition, proceed with normal validation
	}

	// Parse the condition and binding rules
	conditions, bindings := parseCondition(conditionTag)

	// ✅ Step 1: Check if all conditions are met
	conditionMet := true
	var conditionKey, expectedValue string
	for conditionKey, expectedValue = range conditions {
		actualValue, exists := enumMap[conditionKey]
		if !exists || actualValue != expectedValue {
			conditionMet = false
			break // Condition not met, skip validation
		}
	}
	fName := fieldName(f, fullPath)
	// ✅ Step 2: Apply binding rule only if condition is met
	if conditionMet {
		if bindingType, hasBinding := bindings["binding"]; hasBinding {
			switch bindingType {
			case "required":
				if valField.IsZero() {
					return &Error{
						ErrType: "REQUIRED_FIELD_ERR",
						Path:    fName,
						Message: fmt.Sprintf("The field '%s' is required when '%s' is '%s'", fName, conditionKey, expectedValue),
					}
				}
			}
		}
	}

	return nil
}
