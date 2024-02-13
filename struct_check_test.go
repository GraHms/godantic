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

func TestReflectStructSlice(t *testing.T) {
	g := &Validate{}

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

func TestStuckCheck(t *testing.T) {
	t.Run("should explicitly ignore struct fields", func(t *testing.T) {

		g := &Validate{}

		// Test that the function returns an error if a required field is missing.
		type Address struct {
			Street *string `json:"street"`
			City   *string `json:"city" binding:"ignore"`
		}
		type User struct {
			Name    *string  `json:"name" binding:"required"`
			Status  *string  `json:"status" binding:"ignore"`
			Address *Address `json:"address"`
		}
		data := "some val"
		val1 := User{Name: &data, Address: &Address{City: &data, Street: &data}}
		err1 := g.InspectStruct(val1)
		assert.Error(t, err1)
		assert.Equal(t, "INVALID_FIELD_ERR", err1.(*Error).ErrType)
		assert.Equal(t, "address.city", err1.(*Error).Path)

	})
}
