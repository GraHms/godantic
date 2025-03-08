package godantic

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func toPtr[T any](val T) *T {
	return &val
}

type Context struct {
	Type *string `json:"type" binding:"required" enum:"individual,organization"`
}
type U struct {
	ID    *string `json:"id" when:"context.type=individual;binding=required"`
	RegNo *string `json:"reg_no" when:"context.type=organization;binding=required"`
	Email *string `json:"email" when:"context.type=individual;binding=regex;regex=^\\S+@\\S+\\.\\S+$"`
}

type Request struct {
	Context Context `json:"context"`
	User    U       `json:"user"`
}

func TestShouldValidateFirstLevelCondition(t *testing.T) {
	jsonStr := `{
		"context": { "type": "individual" },
		"user": { "id": "12345", "email": "test@example.com" }
	}`
	var req Request
	validate := Validate{}
	err := validate.BindJSON([]byte(jsonStr), &req)
	assert.NoError(t, err, "Validation should pass when first-level condition is met")
}

func TestShouldFailWhenMissingRequiredField(t *testing.T) {
	jsonStr := `{
		"context": { "type": "individual" },
		"user": { "email": "test@example.com" }
	}`
	var req Request
	validate := Validate{}
	err := validate.BindJSON([]byte(jsonStr), &req)
	assert.Error(t, err, "Validation should fail when required field is missing")
	assert.Contains(t, err.Error(), "field 'user.id' is required when 'context.type' is 'individual'")
}

//func TestShouldFailWhenRegexDoesNotMatch(t *testing.T) {
//	jsonStr := `{
//		"context": { "type": "individual" },
//		"user": { "id": "12345", "email": "invalid-email" }
//	}`
//	var req Request
//	validate := Validate{}
//	err := validate.BindJSON([]byte(jsonStr), &req)
//	assert.Error(t, err, "Validation should fail when regex does not match")
//	assert.Contains(t, err.Error(), "field 'email' does not match pattern")
//}

func TestShouldValidateOrganizationWithRequiredField(t *testing.T) {
	jsonStr := `{
		"context": { "type": "organization" },
		"user": { "reg_no": "56789" }
	}`
	var req Request
	validate := Validate{}
	err := validate.BindJSON([]byte(jsonStr), &req)
	assert.NoError(t, err, "Validation should pass when organization has required field")
}

func TestShouldFailWhenMissingRegNoForOrganization(t *testing.T) {
	jsonStr := `{
		"context": { "type": "organization" },
		"user": {}
	}`
	var req Request
	validate := Validate{}
	err := validate.BindJSON([]byte(jsonStr), &req)
	assert.Error(t, err, "Validation should fail when 'reg_no' is missing for organization")
	assert.Contains(t, err.Error(), "field 'user.reg_no' is required when 'context.type' is 'organization'")
}

func TestShouldExtractConditionalFieldsWhenPresent(t *testing.T) {
	// Given a request with a context type "individual"
	testData := Request{
		Context: Context{Type: toPtr("individual")},
		User:    U{ID: toPtr("12345"), Email: toPtr("test@example.com")},
	}

	// Convert struct to reflect.Value
	rootVal := reflect.ValueOf(testData)

	// When we extract conditionally required fields
	result := extractEnumValues(rootVal, "")

	// Then it should correctly store the type condition
	assert.Equal(t, "individual", result["context.type"])
}

func TestShouldNotExtractEmptyFields(t *testing.T) {
	// Given a request where context type is empty
	testData := Request{
		Context: Context{Type: nil}, // No type set
		User:    U{},
	}

	// Convert struct to reflect.Value
	rootVal := reflect.ValueOf(testData)

	// When we extract conditional fields
	result := extractEnumValues(rootVal, "")

	// Then it should not store empty values
	_, exists := result["context.type"]
	assert.False(t, exists, "Should not extract empty fields")
}

// ✅ Test: Should extract fields from nested structures
func TestShouldExtractFieldsFromNestedStructs(t *testing.T) {
	// Given a request with a nested struct
	testData := Request{
		Context: Context{Type: toPtr("organization")},
		User:    U{RegNo: toPtr("56789")},
	}

	// Convert struct to reflect.Value
	rootVal := reflect.ValueOf(testData)

	// When we extract conditional fields
	result := extractEnumValues(rootVal, "")

	// Then it should correctly handle nested JSON paths
	assert.Equal(t, "organization", result["context.type"])
}

// ✅ Test: Should handle pointers correctly
func TestShouldHandlePointersCorrectly(t *testing.T) {
	// Given a request with a pointer field
	testData := &Request{
		Context: Context{Type: toPtr("individual")},
		User:    U{ID: toPtr("12345")},
	}

	// Convert struct to reflect.Value
	rootVal := reflect.ValueOf(testData)

	// When we extract conditional fields
	result := extractEnumValues(rootVal, "")

	// Then it should handle pointers and still extract the correct values
	assert.Equal(t, "individual", result["context.type"])
}
