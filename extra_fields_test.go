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
