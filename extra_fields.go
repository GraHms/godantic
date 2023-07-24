// Package godantic is written by Ismael GraHms
package godantic

import (
	"fmt"
)

type Validate struct {

	IgnoreRequired bool
	IgnoreMinLen   bool

}

func (g *Validate) CheckTypeCompatibility(requestData map[string]interface{}, referenceData map[string]interface{}) error {

	return g.typeCheck(requestData, referenceData, "")

}

func (g *Validate) typeCheck(requestData map[string]interface{}, referenceData map[string]interface{}, currentPath string) error {
	for requestField := range requestData {

		err := g.validateExtra(referenceData, requestField, currentPath)
		if err != nil {
			return err
		}
		// Check if requestField is a valid field in referenceData
		fieldType, _ := referenceData[requestField]
		// Check if field is an object
		err = g.validateObject(fieldType, requestField, requestData, requestField)
		if err != nil {
			return err
		}

		fieldListType, isFieldList := fieldType.([]interface{})
		bodyListValue, isBodyListOk := requestData[requestField].([]interface{})
    
		if isFieldList && isBodyListOk && len(fieldListType) > 0 {
			listObjectType := fieldListType[0].(map[string]interface{})
			for _, obj := range bodyListValue {
				listData := obj.(map[string]interface{})
				result := g.typeCheck(listData, listObjectType, requestField)
				if result != nil {
					return result
				}
			}
		}
	}

	return nil
}
func (g *Validate) validateExtra(referenceData map[string]interface{}, requestField string, currentPath string) error {
	_, isValidField := referenceData[requestField]
	if !isValidField {
		path := requestField
		if len(currentPath) > 0 {
			path = fmt.Sprintf("%s.%s", currentPath, requestField)
		}
		err := &Error{
			ErrType: "INVALID_FIELD_ERR",
			Path:    path,
			Message: fmt.Sprintf("Invalid field <%s>", path),
		}
		return err
	}

	return nil
}

func (g *Validate) validateObject(fieldType interface{}, requestField string, requestData map[string]interface{}, path string) error {
	fieldObjectType, isFieldTypeObject := fieldType.(map[string]interface{})
	bodyObjectValue, _ := requestData[requestField].(map[string]interface{})
	if isFieldTypeObject {
		result := g.typeCheck(bodyObjectValue, fieldObjectType, path)
		if result != nil {
			return result
		}
	}
	return nil
}
