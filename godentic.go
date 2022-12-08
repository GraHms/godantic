package godantic

import "github.com/grahms/godantic/fields"

// Godantic is a struct that can be used to call the ForbidExtraFields function.
type Godantic struct {
	ErrorPathSeparator string
	Fields             *fields.Fields
}

func NewGodantic(ErrorPathSeparator string) *Godantic {
	f := fields.NewFields(ErrorPathSeparator)
	return &Godantic{
		ErrorPathSeparator: ErrorPathSeparator,
		Fields:             f,
	}
}
