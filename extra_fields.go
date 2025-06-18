// Package godantic is written by Ismael GraHms
package godantic

import (
	"fmt"
)

type Validate struct {
	IgnoreRequired     bool
	IgnoreMinLen       bool
	AllowUnknownFields bool
}

func (g *Validate) CheckTypeCompatibility(reqData, refData map[string]any) error {

	return g.typeCheck(reqData, refData, "")

}

func (g *Validate) typeCheck(reqData, refData map[string]any, currentPath string) error {
	for reqField := range reqData {

		if err := g.validateExtra(refData, reqField, currentPath); err != nil {
			return err
		}

		fType := refData[reqField]

		if err := g.validateField(fType, reqData[reqField], reqField); err != nil {
			return err
		}
	}
	return nil
}

func (g *Validate) validateExtra(refData map[string]any, reqField, currentPath string) error {
	if g.AllowUnknownFields || len(refData) == 0 {
		return nil
	}
	_, isValidField := refData[reqField]
	if !isValidField {
		path := reqField
		if len(currentPath) > 0 {
			path = fmt.Sprintf("%s.%s", currentPath, reqField)
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

func (g *Validate) validateObject(fieldType any, reqData map[string]any, reqField, path string) error {
	fieldObjectType, isFieldTypeObject := fieldType.(map[string]any)
	bodyObjectValue, _ := reqData[reqField].(map[string]any)
	if isFieldTypeObject {
		result := g.typeCheck(bodyObjectValue, fieldObjectType, path)
		if result != nil {
			return result
		}
	}
	return nil
}

func (g *Validate) validateField(refType, reqValue any, path string) error {
	switch refTypeAsserted := refType.(type) {
	case Object:
		if _, ok := reqValue.(map[string]any); !ok {
			return fmt.Errorf("expected object for path %s", path)
		}
		// Skip nested validation for Object types
		return nil
	case map[string]any:
		reqMap, ok := reqValue.(map[string]any)
		if !ok {
			return fmt.Errorf("expected map for path %s", path)
		}
		return g.typeCheck(reqMap, refTypeAsserted, path)
	case []any:
		return g.validateList(refTypeAsserted, reqValue, path)
	default:
		//we will Handle other types or validations here, if necessary, in the future
	}
	return nil
}

func (g *Validate) validateList(refListAny, reqValue any, path string) error {
	refList, ok := refListAny.([]any)
	if !ok {
		return fmt.Errorf("expected reference type to be a list for path %s", path)
	}

	reqList, ok := reqValue.([]any)
	if !ok {
		return fmt.Errorf("expected request value to be a list for path %s", path)
	}
	if len(refList) == 0 {
		return nil
	}

	refItem := refList[0]
	for _, item := range reqList {
		if err := g.validateField(refItem, item, path); err != nil {
			return err
		}
	}
	return nil
}

func (g *Validate) constructPath(parent, field string) string {
	if parent == "" {
		return field
	}
	return fmt.Sprintf("%s.%s", parent, field)
}
