package godantic

import (
	"errors"
)

type godentic struct {
}

func New() *godentic {
	return &godentic{}
}

func (g godentic) ForbidExtraFields(bodyData map[string]interface{}, structure map[string]interface{}, parent string) error {
	for bodyField := range bodyData {
		fieldType, isValidField := structure[bodyField]
		if !isValidField {
			return errors.New(joinKeys(parent, bodyField))
		}
		// verify if field is an object
		fieldObjectType, isFieldTypeObject := fieldType.(map[string]interface{})
		bodyObjectValue, _ := bodyData[bodyField].(map[string]interface{})
		if isFieldTypeObject {
			result := g.ForbidExtraFields(bodyObjectValue, fieldObjectType, bodyField)
			if result != nil {
				return result
			}
		}

		fieldListType, isFieldList := fieldType.([]interface{})
		bodyListValue, isBodyListOk := bodyData[bodyField].([]interface{})
		if isFieldList && isBodyListOk {
			listObjectType := fieldListType[0].(map[string]interface{})
			for _, obj := range bodyListValue {
				listData := obj.(map[string]interface{})
				result := g.ForbidExtraFields(listData, listObjectType, bodyField)
				if result != nil {
					return result
				}
			}
		}
		// verify subObjects

	}
	return nil
}
