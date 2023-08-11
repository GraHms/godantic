package godantic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestForbidExtraFields tests the CheckTypeCompatibility function.
func TestForbidExtraFields(t *testing.T) {
	// Create a Validate instance.
	g := &Validate{}

	// Create a request data map.
	requestData := map[string]interface{}{
		"firstName":  "John",
		"lastName":   "Doe",
		"age":        25,
		"extraField": "extraValue",
	}

	// Create a reference data map.
	referenceData := map[string]interface{}{
		"firstName": "string",
		"lastName":  "string",
		"age":       25,
	}

	// with the expected message.
	err := g.CheckTypeCompatibility(requestData, referenceData)
	assert.EqualError(t, err, "Invalid field <extraField>")
}

// TestForbidExtraFields tests the CheckTypeCompatibility function.
func TestForbidExtraFieldsSuccess(t *testing.T) {
	// Create a Validate instance.
	g := &Validate{}

	// Create a request data map.
	requestData := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
		"age":       25,
	}

	// Create a reference data map.
	referenceData := map[string]interface{}{
		"firstName": "string",
		"lastName":  "string",
		"age":       25,
	}

	err := g.CheckTypeCompatibility(requestData, referenceData)
	assert.Equal(t, err, nil)
}

// TestForbidExtraFields tests the CheckTypeCompatibility function.
func TestForbidExtraFieldsInObject(t *testing.T) {
	// Create a Validate instance.
	g := &Validate{}

	// Create a request data map.
	requestData := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
		"age":       25,
		"object": map[string]interface{}{
			"name":       "what",
			"extraField": "extraValue",
		},
	}

	// Create a reference data map.
	referenceData := map[string]interface{}{
		"firstName": "string",
		"lastName":  "string",
		"age":       25,
		"object": map[string]interface{}{
			"name": "what",
		},
	}

	// Call the CheckTypeCompatibility function and check that it returns an error
	// with the expected message.
	err := g.CheckTypeCompatibility(requestData, referenceData)
	assert.EqualError(t, err, "Invalid field <object.extraField>")
}

// TestForbidExtraFields tests the CheckTypeCompatibility function.
func TestForbidExtraFieldsInListObject(t *testing.T) {
	// Create a Validate instance.
	g := &Validate{}

	// Create a request data map.
	requestData := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
		"age":       25,
		"object": []interface{}{
			map[string]interface{}{
				"name":       "what",
				"extraField": "extraValue",
			},
		},
	}

	// Create a reference data map.
	referenceData := map[string]interface{}{
		"firstName": "string",
		"lastName":  "string",
		"age":       25,
		"object": []interface{}{
			map[string]interface{}{
				"name": "what",
			},
		},
	}

	// Call the CheckTypeCompatibility function and check that it returns an error
	// with the expected message.
	err := g.CheckTypeCompatibility(requestData, referenceData)
	assert.EqualError(t, err, "Invalid field <object.extraField>")
}

func TestShouldCheckTypeCompatibility(t *testing.T) {
	// 1. Basic Compatibility
	t.Run("compatible basic types", func(t *testing.T) {
		v := &Validate{}
		reqData := map[string]interface{}{"key": "value"}
		refData := map[string]interface{}{"key": ""}
		err := v.CheckTypeCompatibility(reqData, refData)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// 2. Nested Map Compatibility
	t.Run("compatible nested map", func(t *testing.T) {
		v := &Validate{}
		reqData := map[string]interface{}{"parent": map[string]interface{}{"child": "value"}}
		refData := map[string]interface{}{"parent": map[string]interface{}{"child": ""}}
		err := v.CheckTypeCompatibility(reqData, refData)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// 3. List Compatibility
	t.Run("compatible list types", func(t *testing.T) {
		v := &Validate{}
		reqData := map[string]interface{}{"key": []interface{}{"value1", "value2"}}
		refData := map[string]interface{}{"key": []interface{}{""}}
		err := v.CheckTypeCompatibility(reqData, refData)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	// 4. Invalid Fields
	t.Run("extra field in reqData", func(t *testing.T) {
		v := &Validate{}
		reqData := map[string]interface{}{"key": "value", "extra": "value"}
		refData := map[string]interface{}{"key": ""}
		err := v.CheckTypeCompatibility(reqData, refData)
		if err == nil {
			t.Error("Expected an error for extra field, got none")
		}
	})

	// 5. Default Types (Example for string)
	t.Run("basic string type", func(t *testing.T) {
		v := &Validate{}
		reqData := map[string]interface{}{"key": "value"}
		refData := map[string]interface{}{"key": "default"}
		err := v.CheckTypeCompatibility(reqData, refData)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
