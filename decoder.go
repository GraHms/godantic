package godantic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

func decodeError(err error) error {
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		return &Error{
			ErrType: "TYPE_MISMATCH_ERR",
			Path:    e.Field,
			Message: fmt.Sprintf("The field <%s> was given an invalid type, the expected type is `%s`", e.Field, e.Type.String()),
		}

	case *json.SyntaxError:
		return &Error{
			ErrType: "SYNTAX_ERR",
			Path:    e.Error(),
			Message: e.Error(),
		}
	default:
		return nil
	}
}

func decodeJSON(jsonData []byte, obj interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	err := decoder.Decode(obj)
	if err != nil {
		return decodeError(err)
	}
	return nil
}

func (g *Validate) BindJSON(jsonData []byte, obj interface{}) error {
	err := decodeJSON(jsonData, obj)
	if err != nil {
		return err
	}
	var requestDataMap map[string]interface{}
	err = json.Unmarshal(jsonData, &requestDataMap)
	if err != nil {
		return &Error{
			ErrType: "INVALID_JSON_ERR",
			Path:    "",
			Message: "The given data is not a valid JSON",
		}
	}
	if len(requestDataMap) == 0 {
		return &Error{
			ErrType: "EMPTY_JSON_ERR",
			Path:    "",
			Message: "The given json data is empty",
		}
	}
	var referenceDataMap map[string]interface{}
	refDataBytes, _ := json.Marshal(obj)
	_ = json.Unmarshal(refDataBytes, &referenceDataMap)

	err = decodeJSON(jsonData, obj)
	if err != nil {
		return err
	}
	err = g.InspectStruct(obj)
	if err != nil {
		return err
	}

	err = g.CheckTypeCompatibility(requestDataMap, referenceDataMap)
	if err != nil {
		return err
	}

	return nil
}

func (e *Error) Error() string {
	e.err = errors.New(e.Message)
	return e.err.Error()
}

type Error struct {
	ErrType string
	Message string
	Path    string
	err     error
}
