package godantic

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateWithCustomTag(t *testing.T) {
	type DummyStruct struct {
		Slug  string `json:"slug" validate:"slug"`
		Views int    `json:"views" validate:"positive"`
	}

	RegisterCustom[string]("slug", func(val string, path string) *Error {
		if val != "valid-slug" {
			return &Error{
				ErrType: "INVALID_SLUG",
				Message: "Slug format is invalid",
				// Path is omitted on purpose to test fallback
			}
		}
		return nil
	})

	RegisterCustom[int]("positive", func(val int, path string) *Error {
		if val < 0 {
			return &Error{
				ErrType: "NEGATIVE_VALUE_ERR",
				Message: "Must be non-negative",
				Path:    path, // Path is included to test override prevention
			}
		}
		return nil
	})

	v := &Validate{}

	t.Run("valid slug and views", func(t *testing.T) {
		typ := reflect.TypeOf(DummyStruct{})
		slugField := typ.Field(0)
		viewsField := typ.Field(1)

		err := v.validateWithCustomTag("valid-slug", slugField, "slug")
		assert.Nil(t, err)

		err = v.validateWithCustomTag(100, viewsField, "views")
		assert.Nil(t, err)
	})

	t.Run("invalid slug without path in error", func(t *testing.T) {
		typ := reflect.TypeOf(DummyStruct{})
		field := typ.Field(0)

		err := v.validateWithCustomTag("BAD SLUG", field, "slug")
		assert.NotNil(t, err)
		assert.Equal(t, "slug", err.Path)
		assert.Equal(t, "INVALID_SLUG", err.ErrType)
	})

	t.Run("negative views with path in error", func(t *testing.T) {
		typ := reflect.TypeOf(DummyStruct{})
		field := typ.Field(1)

		err := v.validateWithCustomTag(-5, field, "views")
		assert.NotNil(t, err)
		assert.Equal(t, "views", err.Path)
		assert.Equal(t, "NEGATIVE_VALUE_ERR", err.ErrType)
	})
}
