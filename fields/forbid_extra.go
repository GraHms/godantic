// Package fields is written by Ismael GraHms
package fields

import (
	"reflect"
	"strings"
)

type Fields struct {
	ErrorPathSeparator string
}

func NewFields(ErrorPathSeparator string) *Fields {
	return &Fields{ErrorPathSeparator: ErrorPathSeparator}
}

// ForbidExtra checks if the given map (requestData) has any fields that are not
// present in the reference map (referenceData). If any such fields are found, it returns
// an error with the field names joined together with ".". If referenceData contains a
// field that is an object or a list, this function will recursively call itself to check
// the corresponding field in requestData against the object or list in referenceData.
func (fr *Fields) ForbidExtra(requestData map[string]interface{}, referenceData map[string]interface{}, currentPath string) error {
	for requestField := range requestData {
		// Check if requestField is a valid field in referenceData
		fieldType, isValidField := referenceData[requestField]
		if !isValidField {
			// If requestField is not a valid field in referenceData, return an error with the
			// joined field names
			err := &GodanticError{
				ErrType: "extraFieldErr",
				Path:    strings.Join([]string{currentPath, requestField}, fr.ErrorPathSeparator),
				Message: "Invalid field <" + strings.Join([]string{currentPath, requestField}, fr.ErrorPathSeparator) + ">",
			}
			return err
		}
		// Check if the type of the field in requestData matches the type in referenceData
		if reflect.TypeOf(requestData[requestField]).String() != reflect.TypeOf(fieldType).String() {
			err := &GodanticError{
				ErrType: "typeMismatchErr",
				Path:    strings.Join([]string{currentPath, requestField}, fr.ErrorPathSeparator),
				Message: "The field <" + strings.Join([]string{currentPath, requestField}, ".") + "> has an invalid type",
			}
			return err
		}

		// Check if field is an object
		fieldObjectType, isFieldTypeObject := fieldType.(map[string]interface{})
		bodyObjectValue, _ := requestData[requestField].(map[string]interface{})
		if isFieldTypeObject {
			// If field is an object, recursively call ForbidExtraFields on the object value
			// in requestData
			result := fr.ForbidExtra(bodyObjectValue, fieldObjectType, requestField)
			if result != nil {
				return result
			}
		}

		// Check if field is a list
		fieldListType, isFieldList := fieldType.([]interface{})
		bodyListValue, isBodyListOk := requestData[requestField].([]interface{})
		if isFieldList && isBodyListOk {
			// If field is a list, recursively call ForbidExtraFields on each object in the
			// list
			listObjectType := fieldListType[0].(map[string]interface{})
			for _, obj := range bodyListValue {
				listData := obj.(map[string]interface{})
				result := fr.ForbidExtra(listData, listObjectType, requestField)
				if result != nil {
					return result
				}
			}
		}
	}
	return nil
}
