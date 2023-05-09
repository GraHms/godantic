package godantic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReflectStruct(t *testing.T) {
	g := &Godentic{}

	// Test that the function returns an error if a required field is missing.
	type testStruct1 struct {
		Field1 *string `json:"field_1" binding:"required"`
		Field2 *string `json:"field_2"`
	}
	val1 := testStruct1{}
	err1 := g.InspectStruct(val1, "")
	assert.Error(t, err1)
	assert.Equal(t, "requiredFieldErr", err1.(*GodanticError).ErrType)
	assert.Equal(t, "field_1", err1.(*GodanticError).Path)

	// Test that the function returns an error if a string field is empty.
	type testStruct2 struct {
		Field1 *string `json:"field_1"`
		Field2 *string `json:"field_2"`
	}
	someString := "hello"
	someEmptyString := ""
	val2 := testStruct2{Field1: &someString, Field2: &someEmptyString}
	err2 := g.InspectStruct(val2, "")
	assert.Error(t, err2)
	assert.Equal(t, "emptyStrFieldErr", err2.(*GodanticError).ErrType)
	assert.Equal(t, "field_2", err2.(*GodanticError).Path)

}

func TestReflectStructSlice(t *testing.T) {
	g := &Godentic{}

	// Test that the function returns an error if a required field is missing.
	type sliceStruct struct {
		Field1 *string `json:"field_1" binding:"required"`
		Field2 *string `json:"field_2"`
	}
	type testStruct1 struct {
		Field1     *string         `json:"field_1" binding:"required"`
		Field2     *string         `json:"field_2"`
		SliceField *[]*sliceStruct `json:"slice_field"`
	}
	val1 := testStruct1{
		SliceField: &[]*sliceStruct{},
	}
	err1 := g.InspectStruct(val1, "")
	assert.Error(t, err1)
	assert.Equal(t, "requiredFieldErr", err1.(*GodanticError).ErrType)
	assert.Equal(t, "field_1", err1.(*GodanticError).Path)

	// Test that the function returns an error if a string field is empty.
	type testStruct2 struct {
		Field1 *string `json:"field_1"`
		Field2 *string `json:"field_2"`
	}
	someString := "hello"
	someEmptyString := ""
	val2 := testStruct2{Field1: &someString, Field2: &someEmptyString}
	err2 := g.InspectStruct(val2, "")
	assert.Error(t, err2)
	assert.Equal(t, "emptyStrFieldErr", err2.(*GodanticError).ErrType)
	assert.Equal(t, "field_2", err2.(*GodanticError).Path)

}
