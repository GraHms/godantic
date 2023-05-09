package godantic

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestForbidExtraFields tests the ForbidExtraFields function.
func TestForbidExtraFields(t *testing.T) {
	// Create a Godentic instance.
	g := &Godentic{}

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

	// Convert the reference data map to JSON.
	referenceDataJSON, _ := json.Marshal(referenceData)

	// Create a map that represents the expected result of unmarshalling the
	// reference data JSON.
	referenceDataMap := map[string]interface{}{}
	_ = json.Unmarshal(referenceDataJSON, &referenceDataMap)

	// Call the ForbidExtraFields function and check that it returns an error
	// with the expected message.
	err := g.ForbidExtraFields(requestData, referenceDataMap, "")
	assert.EqualError(t, err, "Invalid field <extraField>")
}

// TestForbidExtraFields tests the ForbidExtraFields function.
func TestForbidExtraFieldsSuccess(t *testing.T) {
	// Create a Godentic instance.
	g := &Godentic{}

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

	// Convert the reference data map to JSON.
	referenceDataJSON, _ := json.Marshal(referenceData)

	// Create a map that represents the expected result of unmarshalling the
	// reference data JSON.
	referenceDataMap := map[string]interface{}{}
	_ = json.Unmarshal(referenceDataJSON, &referenceDataMap)

	// Call the ForbidExtraFields function and check that it returns an error
	// with the expected message.
	err := g.ForbidExtraFields(requestData, referenceDataMap, "")
	assert.Equal(t, err, nil)
}

// TestForbidExtraFields tests the ForbidExtraFields function.
func TestForbidExtraFieldsInObject(t *testing.T) {
	// Create a Godentic instance.
	g := &Godentic{}

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

	// Convert the reference data map to JSON.
	referenceDataJSON, _ := json.Marshal(referenceData)

	// Create a map that represents the expected result of unmarshalling the
	// reference data JSON.
	referenceDataMap := map[string]interface{}{}
	_ = json.Unmarshal(referenceDataJSON, &referenceDataMap)

	// Call the ForbidExtraFields function and check that it returns an error
	// with the expected message.
	err := g.ForbidExtraFields(requestData, referenceDataMap, "")
	assert.EqualError(t, err, "Invalid field <object.extraField>")
}

// TestForbidExtraFields tests the ForbidExtraFields function.
func TestForbidExtraFieldsInListObject(t *testing.T) {
	// Create a Godentic instance.
	g := &Godentic{}

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

	// Convert the reference data map to JSON.
	referenceDataJSON, _ := json.Marshal(referenceData)

	// Create a map that represents the expected result of unmarshalling the
	// reference data JSON.
	referenceDataMap := map[string]interface{}{}
	_ = json.Unmarshal(referenceDataJSON, &referenceDataMap)

	// Call the ForbidExtraFields function and check that it returns an error
	// with the expected message.
	err := g.ForbidExtraFields(requestData, referenceDataMap, "")
	assert.EqualError(t, err, "Invalid field <object.extraField>")
}
