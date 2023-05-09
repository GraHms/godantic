// Package godantic is written by Ismael GraHms
package godantic

import (
	"strings"
)

// Godentic is a struct that can be used to call the ForbidExtraFields function.
type Godentic struct {
}

// ForbidExtraFields checks if the given map (requestData) has any fields that are not
// present in the reference map (referenceData). If any such fields are found, it returns
// an error with the field names joined together with ".". If referenceData contains a
// field that is an object or a list, this function will recursively call itself to check
// the corresponding field in requestData against the object or list in referenceData.
func (g *Godentic) ForbidExtraFields(requestData map[string]interface{}, referenceData map[string]interface{}, currentPath string) error {
	for requestField := range requestData {

		err := g.validateField(referenceData, requestField, currentPath)
		if err != nil {
			return err
		}
		// Check if requestField is a valid field in referenceData
		fieldType, _ := referenceData[requestField]
		// Check if field is an object
		err = g.validateObject(fieldType, requestField, requestData)
		if err != nil {
			return err
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
				result := g.ForbidExtraFields(listData, listObjectType, requestField)
				if result != nil {
					return result
				}
			}
		}
	}
	return nil
}

func (g *Godentic) validateField(referenceData map[string]interface{}, requestField string, currentPath string) error {
	_, isValidField := referenceData[requestField]
	if !isValidField {
		// If requestField is not a valid field in referenceData, return an error with the
		// joined field names
		path := requestField
		if len(currentPath) > 0 {
			path = strings.Join([]string{currentPath, requestField}, ".")
		}
		err := &GodanticError{
			ErrType: "extraFieldErr",
			Path:    path,
			Message: "Invalid field <" + path + ">",
		}
		return err
	}
	return nil
}

func (g *Godentic) validateObject(fieldType interface{}, requestField string, requestData map[string]interface{}) error {
	fieldObjectType, isFieldTypeObject := fieldType.(map[string]interface{})
	bodyObjectValue, _ := requestData[requestField].(map[string]interface{})
	if isFieldTypeObject {
		// If field is an object, recursively call ForbidExtraFields on the object value
		// in requestData
		result := g.ForbidExtraFields(bodyObjectValue, fieldObjectType, requestField)
		if result != nil {
			return result
		}
	}
	return nil
}
