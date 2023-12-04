package godantic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReflectStruct(t *testing.T) {
	g := &Validate{}

	// Test that the function returns an error if a required field is missing.
	type testStruct1 struct {
		Field1 *string `json:"field_1" binding:"required"`
		Field2 *string `json:"field_2"`
	}
	val1 := testStruct1{}
	err1 := g.InspectStruct(val1)
	assert.Error(t, err1)
	assert.Equal(t, "REQUIRED_FIELD_ERR", err1.(*Error).ErrType)
	assert.Equal(t, "field_1", err1.(*Error).Path)

	// Test that the function returns an error if a string field is empty.
	type testStruct2 struct {
		Field1 *string `json:"field_1"`
		Field2 *string `json:"field_2"`
	}
	someString := "hello"
	someEmptyString := ""
	val2 := testStruct2{Field1: &someString, Field2: &someEmptyString}
	err2 := g.InspectStruct(val2)
	assert.Error(t, err2)
	assert.Equal(t, "EMPTY_STRING_ERR", err2.(*Error).ErrType)
	assert.Equal(t, "field_2", err2.(*Error).Path)

}

func TestShouldValidateDefaultValue(t *testing.T) {
	g := &Validate{}

	// Test that the function returns an error if a required field is missing.
	type Product struct {
		Type  string  `json:"type" binding:"required"`
		Name  string  `json:"name"`
		Price *string `json:"price"`
	}
	p := Product{Type: "cheap", Name: "orange"}
	err := g.InspectStruct(p)
	assert.Nil(t, err)
	assert.Equal(t, "orange", p.Name)

}
func TestReflectStructSlice(t *testing.T) {
	g := &Validate{}

	// Test that the function returns an error if a required field is missing.
	type sliceStruct struct {
		Field1 *string `json:"field_1" binding:"required"`
		Field2 *string `json:"field_2"`
	}
	type testStruct1 struct {
		Field1     string          `json:"field_1" binding:"required"`
		Field2     *string         `json:"field_2"`
		SliceField *[]*sliceStruct `json:"slice_field"`
	}
	val1 := testStruct1{
		SliceField: &[]*sliceStruct{},
	}
	err1 := g.InspectStruct(val1)
	assert.Error(t, err1)
	assert.Equal(t, "REQUIRED_FIELD_ERR", err1.(*Error).ErrType)
	assert.Equal(t, "field_1", err1.(*Error).Path)

	// Test that the function returns an error if a string field is empty.
	type testStruct2 struct {
		Field1 *string `json:"field_1"`
		Field2 *string `json:"field_2"`
	}
	someString := "hello"
	someEmptyString := ""
	val2 := testStruct2{Field1: &someString, Field2: &someEmptyString}
	err2 := g.InspectStruct(val2)
	assert.Error(t, err2)
	assert.Equal(t, "EMPTY_STRING_ERR", err2.(*Error).ErrType)
	assert.Equal(t, "field_2", err2.(*Error).Path)

}

func TestReflectStructMap(t *testing.T) {
	g := &Validate{}

	// Test that the function returns an error if a required field is missing.
	type mapStruct struct {
		Field1 *string `json:"field_1" binding:"required"`
		Field2 *string `json:"field_2"`
	}
	type testStruct1 struct {
		Field1   *string                `json:"field_1" binding:"required"`
		Field2   *string                `json:"field_2"`
		MapField *map[string]*mapStruct `json:"map_field"`
	}
	val1 := testStruct1{
		MapField: &map[string]*mapStruct{},
	}
	err1 := g.InspectStruct(val1)
	assert.Error(t, err1)
	assert.Equal(t, "REQUIRED_FIELD_ERR", err1.(*Error).ErrType)
	assert.Equal(t, "field_1", err1.(*Error).Path)

	// Test that the function returns an error if a string field is empty.
	type testStruct2 struct {
		Field1 *string `json:"field_1"`
		Field2 *string `json:"field_2"`
	}
	someString := "hello"
	someEmptyString := ""
	val2 := testStruct2{Field1: &someString, Field2: &someEmptyString}
	err2 := g.InspectStruct(val2)
	assert.Error(t, err2)
	assert.Equal(t, "EMPTY_STRING_ERR", err2.(*Error).ErrType)
	assert.Equal(t, "field_2", err2.(*Error).Path)

}

func TestShouldValidateList(t *testing.T) {
	g := &Validate{}

	// Test that the function returns an error if a required field is missing.
	type testStruct1 struct {
		Field1 *string `json:"field_1" binding:"required"`
		Field2 *string `json:"field_2"`
	}
	val1 := []*testStruct1{
		{},
	}
	err1 := g.InspectStruct(val1)
	assert.Error(t, err1)
	assert.Equal(t, "REQUIRED_FIELD_ERR", err1.(*Error).ErrType)

}

func TestShouldValidateListMinLen(t *testing.T) {
	g := &Validate{}

	// Test that the function returns an error if a required field is missing.
	type testStruct1 struct {
		Field1 *string `json:"field_1" binding:"required"`
		Field2 *string `json:"field_2"`
	}
	var val1 []*testStruct1
	err1 := g.InspectStruct(val1)
	assert.Error(t, err1)
	assert.Equal(t, "EMPTY_LIST_ERR", err1.(*Error).ErrType)

}
