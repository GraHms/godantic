package godantic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type MyDynamicField struct {
	Value     interface{} `json:"value"`
	ValueType string      `json:"valueType" enums:"numeric,string,float,boolean"`
	Attribute string      `json:"attribute"`
}

func (mdf MyDynamicField) GetValue() interface{} {
	return mdf.Value
}

func (mdf MyDynamicField) GetValueType() string {
	return mdf.ValueType
}

func (mdf MyDynamicField) GetAttribute() string {
	return mdf.Attribute
}

func TestDynamicField(t *testing.T) {
	g := &Validate{}

	t.Run("should validate dynamic numeric type", func(t *testing.T) {
		err := g.InspectStruct(MyDynamicField{
			Value:     18,
			ValueType: "numeric",
			Attribute: "age",
		})
		assert.Nil(t, err)
	})

	t.Run("should validate dynamic string type", func(t *testing.T) {
		err := g.InspectStruct(MyDynamicField{
			Value:     "John Doe",
			ValueType: "string",
			Attribute: "name",
		})
		assert.Nil(t, err)
	})

	t.Run("should validate dynamic float type", func(t *testing.T) {
		err := g.InspectStruct(MyDynamicField{
			Value:     3.14,
			ValueType: "float",
			Attribute: "pi",
		})
		assert.Nil(t, err)
	})

	t.Run("should validate dynamic boolean type", func(t *testing.T) {
		err := g.InspectStruct(MyDynamicField{
			Value:     true,
			ValueType: "boolean",
			Attribute: "is_active",
		})
		assert.Nil(t, err)
	})

	t.Run("should handle invalid value type", func(t *testing.T) {
		err := g.InspectStruct(MyDynamicField{
			Value:     "invalid",
			ValueType: "invalid_type",
			Attribute: "invalid_attribute",
		})
		e := err.(*Error)
		assert.NotNil(t, err)
		assert.Equal(t, "INVALID_VALUE_TYPE_ERR", e.ErrType)

	})

	t.Run("should handle invalid numeric value", func(t *testing.T) {
		err := g.InspectStruct(MyDynamicField{
			Value:     "invalid_numeric",
			ValueType: "numeric",
			Attribute: "age",
		})
		assert.NotNil(t, err)
		e := err.(*Error)
		assert.Equal(t, "INVALID_VALUE_TYPE_ERR", e.ErrType)
		assert.Contains(t, e.Message, "Expected numeric value")
	})

	t.Run("should handle invalid string value", func(t *testing.T) {
		err := g.InspectStruct(MyDynamicField{
			Value:     123, // invalid type for string
			ValueType: "string",
			Attribute: "name",
		})
		assert.NotNil(t, err)
		e := err.(*Error)
		assert.Equal(t, "INVALID_VALUE_TYPE_ERR", e.ErrType)
		assert.Contains(t, e.Message, "Expected string value")
	})

	t.Run("should handle invalid float value", func(t *testing.T) {
		err := g.InspectStruct(MyDynamicField{
			Value:     "invalid_float", // invalid type for float
			ValueType: "float",
			Attribute: "pi",
		})
		assert.NotNil(t, err)
		e := err.(*Error)
		assert.Equal(t, "INVALID_VALUE_TYPE_ERR", e.ErrType)
		assert.Contains(t, e.Message, "Expected float value")
	})

	t.Run("should handle invalid boolean value", func(t *testing.T) {
		err := g.InspectStruct(MyDynamicField{
			Value:     "invalid_boolean", // invalid type for boolean
			ValueType: "boolean",
			Attribute: "is_active",
		})
		assert.NotNil(t, err)
		e := err.(*Error)
		assert.Equal(t, "INVALID_VALUE_TYPE_ERR", e.ErrType)
		assert.Contains(t, e.Message, "Expected boolean value")
	})
}

type Address struct {
	Street string `json:"street" binding:"required" min:"3"`
	City   string `json:"city" binding:"required"`
	Zip    string `json:"zip" format:"postal_code"`
}

type Metadata struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type ComplexUser struct {
	Name      string         `json:"name" binding:"required" min:"2"`
	Email     string         `json:"email" format:"email"`
	Age       int            `json:"age" min:"18" max:"99"`
	Role      string         `json:"role" enum:"admin,user,guest"`
	Active    bool           `json:"active"`
	Reason    *string        `json:"reason" when:"active=false;binding=required"`
	Tags      []string       `json:"tags" min:"1"`
	Addresses []Address      `json:"addresses" binding:"required"`
	Metadata  []Metadata     `json:"metadata"`
	Dynamic   MyDynamicField `json:"dynamic"`
}

type NestedDynamic struct {
	Dynamic MyDynamicField `json:"dynamic"`
}

func TestComplexValidation(t *testing.T) {
	g := &Validate{}

	t.Run("should not pass with non boolean attribute", func(t *testing.T) {

		validUser := NestedDynamic{
			Dynamic: MyDynamicField{
				Value:     "sup",
				ValueType: "boolean",
				Attribute: "is_confirmed",
			},
		}
		err := g.InspectStruct(validUser)
		assert.Error(t, err)
	})

	t.Run("should fail if dynamic field is wrong type", func(t *testing.T) {
		invalidUser := ComplexUser{
			Name:  "Jo",
			Email: "john@invalid", // invalid email
			Age:   17,             // too young
			Role:  "superadmin",   // invalid enum
			Tags:  []string{},
			Addresses: []Address{
				{Street: "St", City: "", Zip: "ABC123"}, // too short and missing city
			},
			Dynamic: MyDynamicField{
				Value:     "not-a-bool",
				ValueType: "boolean",
				Attribute: "is_confirmed",
			},
		}
		err := g.InspectStruct(invalidUser)
		assert.NotNil(t, err)
	})
}
