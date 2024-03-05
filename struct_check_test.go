package godantic

import (
	"github.com/stretchr/testify/assert"
	"regexp"
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
func TestAlphanumericValidation(t *testing.T) {
	g := &Validate{}

	// Define a struct with a field requiring validation for alphanumeric characters
	type testStruct struct {
		MyField *string `json:"my_field" binding:"required" regex:"^[A-Za-z0-9]([A-Za-z0-9_-]*[A-Za-z0-9])?$"`
	}

	// Test case: MyField contains only alphanumeric characters
	validValue := "abc123"
	val1 := testStruct{MyField: &validValue}
	err1 := g.InspectStruct(val1)
	assert.NoError(t, err1)

	// Test case: MyField contains non-alphanumeric characters
	invalidValue := "abc$123"
	val2 := testStruct{MyField: &invalidValue}
	err2 := g.InspectStruct(val2)
	assert.Error(t, err2)
	assert.Equal(t, "INVALID_PATTERN_ERR", err2.(*Error).ErrType)
	assert.Equal(t, "my_field", err2.(*Error).Path)
}

func TestRegexValidation(t *testing.T) {
	g := &Validate{}

	// Define a struct for email validation
	type EmailStruct struct {
		Email *string `json:"email" binding:"required" regex:"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
	}

	t.Run("Valid email format", func(t *testing.T) {
		validEmail := "test@example.com"
		val1 := EmailStruct{Email: &validEmail}
		err1 := g.InspectStruct(val1)
		assert.NoError(t, err1)
	})

	t.Run("Invalid email format", func(t *testing.T) {
		invalidEmail := "invalid-email"
		val2 := EmailStruct{Email: &invalidEmail}
		err2 := g.InspectStruct(val2)
		assert.Error(t, err2)
	})

	// Define a struct for phone number validation
	type PhoneStruct struct {
		Phone *string `json:"phone" binding:"required" regex:"^\\+[1-9]\\d{1,14}$"`
	}

	t.Run("Valid phone number format", func(t *testing.T) {
		validPhone := "+1234567890"
		val3 := PhoneStruct{Phone: &validPhone}
		err3 := g.InspectStruct(val3)
		assert.NoError(t, err3)
	})

	t.Run("Invalid phone number format", func(t *testing.T) {
		invalidPhone := "1234567890"
		val4 := PhoneStruct{Phone: &invalidPhone}
		err4 := g.InspectStruct(val4)
		assert.Error(t, err4)
	})

	// Define a struct for string length validation
	type LengthStruct struct {
		Length *string `json:"length" binding:"required" regex:"^.{8,12}$"`
	}

	t.Run("Valid length", func(t *testing.T) {
		validLength := "abcdefgh"
		val5 := LengthStruct{Length: &validLength}
		err5 := g.InspectStruct(val5)
		assert.NoError(t, err5)
	})

	t.Run("Invalid length", func(t *testing.T) {
		invalidLength := "abc"
		val6 := LengthStruct{Length: &invalidLength}
		err6 := g.InspectStruct(val6)
		assert.Error(t, err6)
	})

	// Define a struct for numerical range validation
	type RangeStruct struct {
		Range *string `json:"range" binding:"required" regex:"^([1-9]|[1-5][0-9]|6[0-4])$"`
	}

	t.Run("Valid range", func(t *testing.T) {
		validRange := "50"
		val7 := RangeStruct{Range: &validRange}
		err7 := g.InspectStruct(val7)
		assert.NoError(t, err7)
	})

	t.Run("Invalid range", func(t *testing.T) {
		invalidRange := "70"
		val8 := RangeStruct{Range: &invalidRange}
		err8 := g.InspectStruct(val8)
		assert.Error(t, err8)
	})

	// Define a struct for URL validation
	type URLStruct struct {
		URL *string `json:"url" binding:"required" regex:"^(http|https)://[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}(?:/[^\\s]*)?$"`
	}

	t.Run("Valid URL format", func(t *testing.T) {
		validURL := "https://example.com"
		val9 := URLStruct{URL: &validURL}
		err9 := g.InspectStruct(val9)
		assert.NoError(t, err9)
	})

	t.Run("Invalid URL format", func(t *testing.T) {
		invalidURL := "invalid-url"
		val10 := URLStruct{URL: &invalidURL}
		err10 := g.InspectStruct(val10)
		assert.Error(t, err10)
	})

	// Define a struct for date validation
	type DateStruct struct {
		Date *string `json:"date" binding:"required" regex:"^\\d{4}-\\d{2}-\\d{2}$"`
	}

	t.Run("Valid date format", func(t *testing.T) {
		validDate := "2024-02-29"
		val11 := DateStruct{Date: &validDate}
		err11 := g.InspectStruct(val11)
		assert.NoError(t, err11)
	})

	t.Run("Invalid date format", func(t *testing.T) {
		invalidDate := "invalid-date"
		val12 := DateStruct{Date: &invalidDate}
		err12 := g.InspectStruct(val12)
		assert.Error(t, err12)
	})

	// Define a struct for custom pattern validation
	type CustomStruct struct {
		Custom *string `json:"custom" binding:"required" regex:"^\\d{3}-\\d{2}-\\d{4}$"`
	}

	t.Run("Valid custom format", func(t *testing.T) {
		validCustom := "123-45-6789"
		val13 := CustomStruct{Custom: &validCustom}
		err13 := g.InspectStruct(val13)
		assert.NoError(t, err13)
	})

	t.Run("Invalid custom format", func(t *testing.T) {
		invalidCustom := "invalid-custom"
		val14 := CustomStruct{Custom: &invalidCustom}
		err14 := g.InspectStruct(val14)
		assert.Error(t, err14)
	})

}

func TestStructComposition(t *testing.T) {
	g := &Validate{}

	// Define a struct for email validation
	type Address struct {
		Street *string `json:"street" binding:"required"`
		City   *string `json:"city"`
	}

	type User struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
		Address
	}

	t.Run("Valid email format", func(t *testing.T) {
		validEmail := "test@example.com"
		val1 := User{Email: &validEmail, Address: Address{Street: &validEmail}}
		err1 := g.InspectStruct(val1)
		assert.NoError(t, err1)
	})
}

func TestFormatValidation(t *testing.T) {
	g := &Validate{}

	// Define a struct for email validation
	type EmailStruct struct {
		Email *string `json:"email" binding:"required" format:"email"`
	}

	t.Run("Valid email format", func(t *testing.T) {
		validEmail := "test@example.com"
		val1 := EmailStruct{Email: &validEmail}
		err1 := g.InspectStruct(val1)
		assert.NoError(t, err1)
	})

	t.Run("Invalid email format", func(t *testing.T) {
		invalidEmail := "invalid-email"
		val2 := EmailStruct{Email: &invalidEmail}
		err2 := g.InspectStruct(val2)
		assert.Error(t, err2)
	})

	// Define a struct for URL validation
	type URLStruct struct {
		URL *string `json:"url" binding:"required" format:"url"`
	}

	t.Run("Valid URL format", func(t *testing.T) {
		validURL := "https://example.com"
		val3 := URLStruct{URL: &validURL}
		err3 := g.InspectStruct(val3)
		assert.NoError(t, err3)
	})

	t.Run("Invalid URL format", func(t *testing.T) {
		invalidURL := "invalid-url"
		val4 := URLStruct{URL: &invalidURL}
		err4 := g.InspectStruct(val4)
		assert.Error(t, err4)
	})
}

func TestGetFormatRegex(t *testing.T) {
	testCases := []struct {
		formatTag string
		input     string
		expected  bool
	}{
		{"email", "test@example.com", true},
		{"email", "invalid-email", false},
		{"url", "http://example.com", true},
		{"url", "invalid-url", false},
		{"date", "2024-02-29", true},
		{"date", "invalid-date", false},
		{"time", "12:34:56", true},
		{"time", "25:00:00", false},
		{"uuid", "123e4567-e89b-12d3-a456-426614174000", true},
		{"uuid", "invalid-uuid", false},
		{"ip", "192.168.1.1", true},
		{"ip", "invalid-ip", false},
		{"credit_card", "1234-5678-9012-3456", true},
		{"credit_card", "invalid-credit-card", false},
		{"postal_code", "12345", true},
		{"postal_code", "invalid-postal-code", false},
		{"phone", "+1234567890", true},
		{"phone", "1234567890", false},
		{"ssn", "123-45-6789", true},
		{"ssn", "invalid-ssn", false},
		{"credit_card_expiry", "02/2024", true},
		{"credit_card_expiry", "invalid-expiry-date", false},
		{"latitude", "45.678", true},
		{"latitude", "invalid-latitude", false},
		{"longitude", "-123.456", true},
		{"longitude", "invalid-longitude", false},
		{"hex_color", "#FFFFFF", true},
		{"hex_color", "invalid-hex-color", false},
		{"mac_address", "00:0a:95:9d:68:16", true},
		{"mac_address", "invalid-mac-address", false},
		{"mz-msisdn", "258123456789", true},
		{"mz-msisdn", "123456789", false},
		{"mz-nuit", "123456789", true},
		{"mz-nuit", "invalid-nuit", false},
	}

	for _, tc := range testCases {
		t.Run(tc.formatTag, func(t *testing.T) {
			regexPattern := getFormatRegex(tc.formatTag)
			match := regexPattern != "" && regexp.MustCompile(regexPattern).MatchString(tc.input)
			if match != tc.expected {
				t.Errorf("For format tag %s and input %s, expected match: %v, got: %v", tc.formatTag, tc.input, tc.expected, match)
			}
		})
	}
}
