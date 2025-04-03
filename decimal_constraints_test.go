package godantic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecimalConstraints(t *testing.T) {
	g := &Validate{}

	t.Run("should allow valid float with max_digits and decimal_places", func(t *testing.T) {
		type S struct {
			Value *float64 `json:"value" max_digits:"6" decimal_places:"2"`
		}
		val := 1234.56
		assert.NoError(t, g.InspectStruct(S{Value: &val}))
	})

	t.Run("should fail when exceeding max_digits", func(t *testing.T) {
		type S struct {
			Value *float64 `json:"value" max_digits:"5"`
		}
		val := 12345.6 // 6 digits total
		err := g.InspectStruct(S{Value: &val})
		assert.Error(t, err)
		assert.Equal(t, "MAX_DIGITS_ERR", err.(*Error).ErrType)
	})

	t.Run("should fail when exceeding decimal_places", func(t *testing.T) {
		type S struct {
			Value *float64 `json:"value" decimal_places:"2"`
		}
		val := 123.456 // 3 decimal places
		err := g.InspectStruct(S{Value: &val})
		assert.Error(t, err)
		assert.Equal(t, "DECIMAL_PLACES_ERR", err.(*Error).ErrType)
	})

	t.Run("should ignore trailing zeros in decimal places", func(t *testing.T) {
		type S struct {
			Value *float64 `json:"value" decimal_places:"2"`
		}
		val := 12.3000 // 1 decimal digit
		assert.NoError(t, g.InspectStruct(S{Value: &val}))
	})

	t.Run("should ignore leading zero in int part", func(t *testing.T) {
		type S struct {
			Value *float64 `json:"value" max_digits:"3"`
		}
		val := 0.99 // only 2 digits (0 is ignored)
		assert.NoError(t, g.InspectStruct(S{Value: &val}))
	})

	t.Run("should allow exact match with max_digits and decimal_places", func(t *testing.T) {
		type S struct {
			Value *float64 `json:"value" max_digits:"6" decimal_places:"2"`
		}
		val := 9999.99 // 6 digits, 2 decimals
		assert.NoError(t, g.InspectStruct(S{Value: &val}))
	})

	t.Run("should fail with negative number if it exceeds digits", func(t *testing.T) {
		type S struct {
			Value *float64 `json:"value" max_digits:"3"`
		}
		val := -123.45 // 5 digits (ignore '-')
		err := g.InspectStruct(S{Value: &val})
		assert.Error(t, err)
		assert.Equal(t, "MAX_DIGITS_ERR", err.(*Error).ErrType)
	})
}
