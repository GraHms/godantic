package godantic

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestMinMaxValidation(t *testing.T) {
	g := &Validate{}

	t.Run("String length constraints", func(t *testing.T) {
		type TestStruct struct {
			Name *string `json:"name" min:"3" max:"5"`
		}

		short := "Hi"
		long := "HelloWorld"
		valid := "John"

		assert.Error(t, g.InspectStruct(TestStruct{Name: &short}))
		assert.Error(t, g.InspectStruct(TestStruct{Name: &long}))
		assert.NoError(t, g.InspectStruct(TestStruct{Name: &valid}))
	})

	t.Run("Slice length constraints", func(t *testing.T) {
		type TestStruct struct {
			Items *[]string `json:"items" min:"2" max:"4"`
		}

		short := &[]string{"a"}
		long := &[]string{"a", "b", "c", "d", "e"}
		valid := &[]string{"a", "b"}

		assert.Error(t, g.InspectStruct(TestStruct{Items: short}))
		assert.Error(t, g.InspectStruct(TestStruct{Items: long}))
		assert.NoError(t, g.InspectStruct(TestStruct{Items: valid}))
	})

	t.Run("Integer value constraints", func(t *testing.T) {
		type TestStruct struct {
			Age *int `json:"age" min:"18" max:"30"`
		}

		young := 17
		old := 35
		valid := 25

		assert.Error(t, g.InspectStruct(TestStruct{Age: &young}))
		assert.Error(t, g.InspectStruct(TestStruct{Age: &old}))
		assert.NoError(t, g.InspectStruct(TestStruct{Age: &valid}))
	})

	t.Run("Float value constraints", func(t *testing.T) {
		type TestStruct struct {
			Score *float64 `json:"score" min:"1.5" max:"3.5"`
		}

		low := 1.4
		high := 3.6
		valid := 2.7

		assert.Error(t, g.InspectStruct(TestStruct{Score: &low}))
		assert.Error(t, g.InspectStruct(TestStruct{Score: &high}))
		assert.NoError(t, g.InspectStruct(TestStruct{Score: &valid}))
	})
}

func TestNumericConstraints(t *testing.T) {
	g := &Validate{}

	t.Run("should validate greater than (gt)", func(t *testing.T) {
		type S struct {
			Age *int `json:"age" gt:"18"`
		}
		age := 20
		tooYoung := 18
		assert.NoError(t, g.InspectStruct(S{Age: &age}))
		assert.Error(t, g.InspectStruct(S{Age: &tooYoung}))
	})

	t.Run("should validate greater or equal (ge)", func(t *testing.T) {
		type S struct {
			Age *int `json:"age" ge:"18"`
		}
		age := 18
		younger := 17
		assert.NoError(t, g.InspectStruct(S{Age: &age}))
		assert.Error(t, g.InspectStruct(S{Age: &younger}))
	})

	t.Run("should validate less than (lt)", func(t *testing.T) {
		type S struct {
			Age *int `json:"age" lt:"65"`
		}
		age := 60
		tooOld := 70
		assert.NoError(t, g.InspectStruct(S{Age: &age}))
		assert.Error(t, g.InspectStruct(S{Age: &tooOld}))
	})

	t.Run("should validate less or equal (le)", func(t *testing.T) {
		type S struct {
			Age *int `json:"age" le:"30"`
		}
		age := 30
		tooOld := 31
		assert.NoError(t, g.InspectStruct(S{Age: &age}))
		assert.Error(t, g.InspectStruct(S{Age: &tooOld}))
	})

	t.Run("should validate multiple_of", func(t *testing.T) {
		type S struct {
			Num *int `json:"num" multiple_of:"5"`
		}
		valid := 25
		invalid := 23
		assert.NoError(t, g.InspectStruct(S{Num: &valid}))
		assert.Error(t, g.InspectStruct(S{Num: &invalid}))
	})

	t.Run("should validate float multiple_of", func(t *testing.T) {
		type S struct {
			Num *float64 `json:"num" multiple_of:"0.5"`
		}
		valid := 2.5
		invalid := 2.3
		assert.NoError(t, g.InspectStruct(S{Num: &valid}))
		assert.Error(t, g.InspectStruct(S{Num: &invalid}))
	})

	t.Run("should reject NaN and Inf by default", func(t *testing.T) {
		type S struct {
			Value *float64 `json:"value"`
		}
		nan := math.NaN()
		posInf := math.Inf(1)
		negInf := math.Inf(-1)
		assert.Error(t, g.InspectStruct(S{Value: &nan}))
		assert.Error(t, g.InspectStruct(S{Value: &posInf}))
		assert.Error(t, g.InspectStruct(S{Value: &negInf}))
	})

	t.Run("should allow NaN and Inf when allow_inf_nan is true", func(t *testing.T) {
		type S struct {
			Value *float64 `json:"value" allow_inf_nan:"true"`
		}
		nan := math.NaN()
		posInf := math.Inf(1)
		negInf := math.Inf(-1)
		assert.NoError(t, g.InspectStruct(S{Value: &nan}))
		assert.NoError(t, g.InspectStruct(S{Value: &posInf}))
		assert.NoError(t, g.InspectStruct(S{Value: &negInf}))
	})
}
