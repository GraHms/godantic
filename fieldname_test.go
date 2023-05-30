package godantic

import (
	"reflect"
	"testing"
)

func TestFieldName(t *testing.T) {
	// Define a test struct field with a "json" tag
	field := reflect.StructField{
		Name: "myField",
		Tag:  reflect.StructTag(`json:"myJsonField"`),
	}

	// Test with no tree string
	expectedName := "myJsonField"
	actualName := fieldName(field, "")
	if actualName != expectedName {
		t.Errorf("Expected %s but got %s", expectedName, actualName)
	}

	// Test with a tree string
	expectedName = "myTree.myJsonField"
	actualName = fieldName(field, "myTree")
	if actualName != expectedName {
		t.Errorf("Expected %s but got %s", expectedName, actualName)
	}

	// Test with no json tag
	field.Tag = reflect.StructTag(``)
	expectedName = ""
	actualName = fieldName(field, "")
	if actualName != expectedName {
		t.Errorf("Expected %s but got %s", expectedName, actualName)
	}
}
