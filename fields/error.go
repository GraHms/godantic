package fields

import "errors"

type GodanticError struct {
	ErrType string
	Message string
	Path    string
	err     error
}

func (e *GodanticError) Error() string {
	e.err = errors.New(e.Message)
	return e.err.Error()
}
