package godantic

import (
	"encoding/json"
	"testing"
)

func TestBindJSON(t *testing.T) {
	// Test case 1: Valid JSON data and object
	jsonData := []byte(`{"name": "John", "age": 30}`)
	obj := &struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	v := &Validate{}
	err := v.BindJSON(jsonData, obj)

	if err != nil {
		t.Errorf("Test case 1 failed. Unexpected error: %v", err)
	}

	// Verify the object has been populated correctly
	expectedObj := &struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}
	if !jsonEqual(obj, expectedObj) {
		t.Errorf("Test case 1 failed. Object mismatch. Expected: %+v, Got: %+v", expectedObj, obj)
	}

	// Test case 2: Invalid JSON data
	jsonData = []byte(`{"name": "John", "age": }`) // Invalid JSON, missing value for "age"
	err = v.BindJSON(jsonData, obj)

	if err == nil {
		t.Errorf("Test case 2 failed. Expected an error, but got nil")
	} else if _, ok := err.(*Error); !ok {
		t.Errorf("Test case 2 failed. Expected an *Error, but got %T", err)
	}

	// Test case 3: Empty JSON data
	jsonData = []byte(`{}`)
	err = v.BindJSON(jsonData, obj)

	if err == nil {
		t.Errorf("Test case 3 failed. Expected an error, but got nil")
	} else if _, ok := err.(*Error); !ok {
		t.Errorf("Test case 3 failed. Expected an *Error, but got %T", err)
	}
}

func jsonEqual(a, b interface{}) bool {
	jsonA, _ := json.Marshal(a)
	jsonB, _ := json.Marshal(b)
	return string(jsonA) == string(jsonB)
}
