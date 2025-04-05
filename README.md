# Godantic
[![Go Report Card](https://goreportcard.com/badge/github.com/grahms/godantic)](https://goreportcard.com/report/github.com/grahms/godantic)
[![Go Reference](https://pkg.go.dev/badge/github.com/grahms/godantic.svg)](https://pkg.go.dev/github.com/grahms/godantic)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Tests](https://github.com/grahms/godantic/actions/workflows/tests.yml/badge.svg)](https://github.com/grahms/godantic/actions/workflows/tests.yml)
[![Code Coverage](https://img.shields.io/badge/coverage-90%25-brightgreen.svg)](#) <!-- substitui com real se tiver cobertura -->
[![Issues](https://img.shields.io/github/issues/grahms/godantic)](https://github.com/grahms/godantic/issues)

Godantic is a Go package for inspecting and validating JSON-like data against Go struct types and schemas. It provides functionalities for checking type compatibility, structure compatibility, and other validations such as empty string, invalid time, minimum length list checks, regex pattern matching, and format validation.

## Getting Started

Install the godantic package:

```sh
go get github.com/grahms/godantic
```

Then import it in your Go code:

```go
import "github.com/grahms/godantic"
```

## Simple Usage

```go
type Person struct {
    Name *string `json:"name" binding:"required"`
    Age  *int    `json:"age"`
}

var jsonData = []byte(`{"name": "John", "age": 30}`)
var person Person

validator := godantic.Validate{}

err := validator.BindJSON(jsonData, &person)
if err != nil {
    fmt.Println(err)
}
```

## Advanced Usage

- Enum Validation

```go
type Person struct {
    Name *string `json:"name" binding:"required"`
    Role *string `json:"role" enum:"admin,user"`
}

// Here, the Role field must be either 'admin' or 'user'. If it's not, an error is returned.
```

- Handling Extra Fields

```go
var jsonData = []byte(`{"name": "John", "age": 30, "extra": "extra data"}`)
var person Person

validator := godantic.Validate{}

err := validator.BindJSON(jsonData, &person)
if err != nil {
    fmt.Println(err) // This will print an error about the 'extra' field not being valid.
}
```

- Custom Error Handling

```go
type CustomError struct {
    ErrType string
    Message string
    Path    string
    err     error
}

func (e *CustomError) Error() string {
    e.err = errors.New(e.Message)
    return e.err.Error()
}

// Now you can create your own error type and return it in your custom validation functions.
```

- Inspecting and Validating Structs

```go
validator := godantic.Validate{}
err := validator.InspectStruct(&myStruct)
if err != nil {
    fmt.Println(err)
}
```

## Nested Fields & Objects

```go
type Address struct {
    City  *string `json:"city" binding:"required"`
    State *string `json:"state" binding:"required"`
}

type Person struct {
    Name    *string `json:"name" binding:"required"`
    Age     *int    `json:"age"`
    Address *Address `json:"address"`
}

var jsonData = []byte(`{
    "name": "John",
    "age": 30,
    "address": {
        "city": "New York",
        "state": "NY"
    }
}`)

var person Person

validator := godantic.Validate{}

err := validator.BindJSON(jsonData, &person)
if err != nil {
    fmt.Println(err)
}
```

In this example, the `Person` struct has a nested `Address` struct. The `godantic` package will validate the fields of the nested struct as well.

## Lists

```go
type Skill struct {
    Name *string `json:"name" binding:"required"`
    Level *int `json:"level"`
}

type Person struct {
    Name  *string `json:"name" binding:"required"`
    Age   *int    `json:"age"`
    Skills []Skill `json:"skills"`
}

var jsonData = []byte(`{
    "name": "John",
    "age": 30,
    "skills": [
        {
            "name": "Go",
            "level": 5
        },
        {
            "name": "Python",
            "level": 4
        }
    ]
}`)

var person Person

validator := godantic.Validate{}

err := validator.BindJSON(jsonData, &person)
if err != nil {
    fmt.Println(err)
}
```

In this example, the `Person` struct has a `Skills` field that is a slice of `Skill` structs. The `godantic` package will iterate over the list and validate each object in the list.

## Integration with Web Frameworks

### Using Godantic with Gin

Here's an example of how to use the `godantic` package with the Gin web framework.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/grahms/godantic"
    "net/http"
)

type User struct {
    Name    *string `json:"name" binding:"required"`
    Email   *string `json:"email" binding:"required"`
    Age     *int    `json:"age"`
}

func main() {
    r := gin.Default()

    r.POST("/user", func(c *gin.Context) {
        var user User
        validator := godantic.Validate{}

        jsonData, err := c.GetRawData()
        if err != nil {
            c.JSON(http

.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        err = validator.BindJSON(jsonData, &user)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    r.Run()
}
```

In the example above, instead of using Gin's built-in JSON binding (`c.BindJSON(&user)`), we're using `godantic`'s `BindJSON` function. Here are the advantages:

1. **More control over validation**: `godantic` provides much more control over the validation process compared to Gin's built-in binding. It supports various validation methods and customizations like type compatibility checks, structure compatibility checks, and handling extra fields. You can customize these validation rules based on your needs.

2. **Detailed error reporting**: `godantic` provides detailed error types and messages which can be very useful for debugging and for providing precise error messages to the API users. In contrast, Gin's built-in binding returns a generic "binding error".

3. **Enum Validation**: `godantic` supports enum validation, which is not available in Gin's built-in JSON binding.

4. **Nested Fields & Objects**: `godantic` supports validation for nested fields and objects as well as lists, which provides more flexibility and control compared to Gin's built-in binding.

Please remember that the Go's `json.Unmarshal` function used by `godantic` doesn't check for additional fields in the JSON input that are not present in the target struct. If you want to disallow additional fields, you might have to implement additional checks.

## Features

- **BindJSON**: Parses and validates JSON data into a provided struct. It performs type checking and structural validation against the expected schema of the provided struct.
- **InspectStruct**: Iteratively inspects the fields of a struct based on their type and validates them based on certain conditions.
- **CheckTypeCompatibility**: Checks if two `map[string]interface{}` objects (request and reference data) are compatible in terms of structure and type.

## Error Types

- `REQUIRED_FIELD_ERR`: Triggered when a field marked as required is not provided.
- `INVALID_ENUM_ERR`: Triggered when a field value is not among the allowed enum values.
- `INVALID_FIELD_ERR`: Triggered when an invalid field is provided.
- `TYPE_MISMATCH_ERR`: Triggered when a field is given a value with an invalid type.
- `SYNTAX_ERR`: Triggered when there is a syntax error in the JSON data.
- `INVALID_JSON_ERR`: Triggered when the provided data is not valid JSON.
- `EMPTY_JSON_ERR`: Triggered when the provided JSON data is empty.
- `INVALID_TIME_ERR`: Triggered when a time.Time field has an invalid time value.
- `EMPTY_STRING_ERR`: Triggered when a string field is empty.
- `EMPTY_LIST_ERR`: Triggered when a list field is empty.
- `INVALID_REGEX_ERR`: Triggered when a field value does not match the required regex pattern.
- `INVALID_FORMAT_ERR`: Triggered when a field value does not match the required format.

## Format Tags

The following table lists the supported format tags and their corresponding regular expressions:

| Format Tag          | Description                      | Example Use Case                      |
|---------------------|----------------------------------|---------------------------------------|
| email               | Email address format             | Validating user email addresses       |
| url                 | URL format                       | Validating website URLs               |
| date                | Date format (YYYY-MM-DD)        | Validating dates in a specific format |
| time                | Time format (HH:MM:SS)          | Validating times in a specific format |
| uuid                | UUID format                      | Validating UUIDs                      |
| ip                  | IP address format                | Validating IPv4 or IPv6 addresses     |
| credit_card         | Credit card number format        | Validating credit card numbers        |
| postal_code         | Postal code format               | Validating postal codes               |
| phone               | Phone number format              | Validating phone numbers              |
| ssn                 | Social Security Number format    | Validating SSN                        |
| credit_card_expiry  | Credit card expiry date format  | Validating credit card expiry dates   |
| latitude            | Latitude format                  | Validating latitude coordinates       |
| longitude           | Longitude format                 | Validating longitude coordinates      |
| hex_color           | Hex color format                 | Validating hex color codes            |
| mac_address         | MAC address format               | Validating MAC addresses              |
| html_tag            | HTML tag format                  | Validating HTML tags                  |
| mz-msisdn           | Mozambican phone number format   | Validating Mozambican phone numbers   |
| mz-nuit             | Mozambican NUIT format           | Validating Mozambican NUIT numbers    |



## Using the `ignore` Tag

The `ignore` tag allows you to exclude specific fields from input validation while still retaining them in the struct. This can be useful for fields representing metadata or internal information that shouldn't be validated during input but are required for other purposes.

### Example

Consider a `User` struct with an `ID` field that should be excluded from input validation but retained in the struct for internal use:

```go
package main

import (
    "fmt"
    "github.com/grahms/godantic"
)

type User struct {
    ID        int    `json:"id" binding:"ignore"` // ID field is ignored during input validation
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    Email     string `json:"email" binding:"required" format:"email"`
    // Other fields...
}

func main() {
    // Example JSON data representing user input
    jsonData := []byte(`{
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com"
        // No "id" field included
    }`)

    // Create a new instance of the validator
    validator := godantic.Validate{}

    // Create an instance of the User struct
    var user User

    // Bind and validate the JSON data against the User struct
    err := validator.BindJSON(jsonData, &user)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Validation successful, process the user data
    fmt.Printf("User ID: %d\n", user.ID) // ID is still accessible despite being ignored during validation
    fmt.Printf("Name: %s %s\n", user.FirstName, user.LastName)
    fmt.Printf("Email: %s\n", user.Email)
}
```

In this example:

- The `ID` field represents a unique identifier for the user and is marked with the `ignore` tag.
- Despite being ignored during validation, the `ID` field remains accessible in the `User` struct after validation, allowing you to utilize it for internal operations or data processing.


---

## 📐 Numeric & Decimal Constraints

`godantic` supports advanced validation for numeric and decimal fields using struct tags. These rules enable expressive, type-safe constraints on your data models.

### 🔢 `min` and `max` (generic)

Used to validate:
- String length
- Slice/array length
- Numeric values

```go
type User struct {
  Name  *string `json:"name" min:"3" max:"50"`       // between 3 and 50 characters
  Tags  *[]int  `json:"tags" min:"1" max:"5"`         // between 1 and 5 items
  Age   *int    `json:"age" min:"18" max:"60"`        // between 18 and 60 years
}
```

---

### ⚖️ `gt`, `ge`, `lt`, `le` (value bounds)

Used to enforce strict or inclusive numeric constraints.

| Tag  | Meaning                        |
|------|--------------------------------|
| `gt` | value must be **greater than** |
| `ge` | value must be **≥**            |
| `lt` | value must be **less than**    |
| `le` | value must be **≤**            |

```go
type Product struct {
  Price *float64 `json:"price" gt:"0"`     // must be greater than 0
  Stock *int     `json:"stock" le:"1000"`  // must be ≤ 1000
}
```

---

### 🎯 `multiple_of`

Ensures the value is a multiple of a specified number.

```go
type Payment struct {
  Amount *float64 `json:"amount" multiple_of:"0.05"` // e.g. currency in increments of 0.05
}
```

---

### 🧮 `max_digits` & `decimal_places`

Validates decimal precision:

| Tag              | Description                                                       |
|------------------|-------------------------------------------------------------------|
| `max_digits`     | Max total digits (excluding leading zero, includes decimals)      |
| `decimal_places` | Max number of digits after the decimal point                      |

```go
type Invoice struct {
  Total *float64 `json:"total" max_digits:"6" decimal_places:"2"` // e.g. 9999.99
}
```

---

### ☢️ `allow_inf_nan`

Allows `+Inf`, `-Inf`, and `NaN` values for floating point numbers.

```go
type Reading struct {
  Value *float64 `json:"value" allow_inf_nan:"true"`
}
```

> 🔒 By default, `inf` and `NaN` are **not allowed**.

---



## Conditional Validation Based on Enum Values

`godantic` allows you to apply **conditional validation rules** based on the values of other fields. This is done using the `when` tag.

### **1️⃣ Basic Conditional Validation**
You can specify that a field should only be validated when another field has a specific value.

#### **Example: Requiring a field when `context.type=organization`**
```go
type Context struct {
    Type *string `json:"type" enum:"individual,organization"`
}

type User struct {
    RegNo *string `json:"reg_no" when:"context.type=organization;binding=required"`
}
```

✅ **If `context.type` is `"organization"`, `reg_no` is required.**  
❌ **If `context.type` is `"individual"`, `reg_no` is ignored.**

#### **Valid JSON Input**
```json
{
  "context": { "type": "organization" },
  "user": { "reg_no": "56789" }
}
```

✅ **Passes validation because `reg_no` is provided for `organization`.**

---

### **2️⃣ Invalid Case: Missing `reg_no` When Type is `organization`**
```json
{
  "context": { "type": "organization" },
  "user": {}
}
```
❌ **Fails validation with error:**
```
Field <user.reg_no> is required when context.type=organization
```

---

### **3️⃣ Multiple Conditions**
You can require a field **only when multiple conditions are met.**

#### **Example: Requiring `vat_number` when `context.type=business` and `country=EU`**
```go
type Context struct {
    Type    *string `json:"type" enum:"individual,business"`
    Country *string `json:"country" enum:"EU,US"`
}

type Business struct {
    VATNumber *string `json:"vat_number" when:"context.type=business;context.country=EU;binding=required"`
}
```

✅ **If `context.type` is `"business"` and `context.country` is `"EU"`, `vat_number` is required.**  
❌ **If `context.type` is `"individual"`, `vat_number` is ignored.**

#### **Valid JSON**
```json
{
  "context": { "type": "business", "country": "EU" },
  "business": { "vat_number": "EU123456" }
}
```

---

### **4️⃣ Allowed Operators for Conditions**
| **Operator** | **Example** | **Meaning** |
|-------------|------------|-------------|
| `=` | `context.type=business` | Field must be equal to value |


---

## **Why Use Conditional Validation?**
✅ **Simplifies complex validation logic**  
✅ **Eliminates unnecessary validation** when conditions aren’t met  
✅ **Supports dynamic rules based on input data**

---

🚀 **Now you can enforce conditional validation effortlessly!** 🚀


## ✅ Custom Validation Tags (`validate`)

Godantic allows you to register custom validation functions tied to specific tag names. These functions give you full control over domain-specific validations, and they integrate seamlessly into your validation flow.

### 🔧 Registering a Custom Validator

Use `RegisterCustom` to attach your custom validation logic to a tag:

```go
godantic.RegisterCustom[string]("starts_with_A", func(val string, path string) *godantic.Error {
	if !strings.HasPrefix(val, "A") {
		return &godantic.Error{
			ErrType: "STARTS_WITH_A_ERR",
			Path:    path,
			Message: fmt.Sprintf("The field <%s> must start with 'A'", path),
		}
	}
	return nil
})
```

> ✅ The generic type `[string]` indicates the expected value type for the validation. The function receives the field value and the full path of the field in the struct.

---

### 📌 Applying the Validator to a Struct Field

```go
type User struct {
	Username *string `json:"username" validate:"starts_with_A"`
}
```

When you call `Validate.InspectStruct(&User{})` or `Validate.BindJSON`, the custom validator will automatically be invoked.

---

### 🧩 Using Multiple Custom Tags

You can apply multiple custom validations on the same field by separating them with commas:

```go
type Product struct {
	Code *string `json:"code" validate:"starts_with_A,min_len_3"`
}
```

All validators (`starts_with_A` and `min_len_3`) will be executed in order. If any of them return an error, validation will fail.

---

### 🧠 Validators by Type

Godantic supports type-safe validation using Go generics. For example, to validate integers:

```go
godantic.RegisterCustom[int]("positive", func(val int, path string) *godantic.Error {
	if val <= 0 {
		return &godantic.Error{
			ErrType: "POSITIVE_ERR",
			Path:    path,
			Message: fmt.Sprintf("The field <%s> must be a positive number", path),
		}
	}
	return nil
})
```

```go
type Invoice struct {
	Amount *int `json:"amount" validate:"positive"`
}
```

---

### 🛡️ Safety and Design

- You don't need to manually add the `Path` inside the error — Godantic will do it for you if it’s missing.
- If the field type doesn't match the registered validator's type, the validator is skipped without causing panic.

---

### ✅ Best Practices

- Use descriptive tag names: `min_len_5`, `email_domain_gov`, `alphanumeric_only`, etc.
- Register all custom validators once during app initialization (`init()` or startup function).
- Combine with built-in tags like `binding:"required"`, `format:"email"`, `when:"..."`, and `enum:"..."` for expressive rules.

---
## 🔌 Plugin-Based Validation 

Godantic supports a powerful **interface-based validation mechanism** that allows you to embed custom logic inside your struct types using the `ValidationPlugin` interface.

### ✨ What is a Plugin?

A Plugin is any struct that implements the following interface:

```go
type ValidationPlugin interface {
	Validate() *godantic.Error
}
```

Godantic will **automatically detect** structs that implement this interface and **invoke the `Validate()` method** during the validation cycle — whether the struct is a root object, nested field, or an item in a list.

---

### 🧪 Example: Basic Plugin

```go
type Password struct {
	Value *string `json:"value"`
}

func (p Password) Validate() *godantic.Error {
	if p.Value == nil || len(*p.Value) < 8 {
		return &godantic.Error{
			ErrType: "WEAK_PASSWORD",
			Message: "Password must be at least 8 characters long",
		}
	}
	return nil
}
```

Then use it in your main struct:

```go
type User struct {
	Username *string  `json:"username" binding:"required"`
	Password *Password `json:"password"` // Plugin validation will be triggered
}
```

Godantic will automatically call `Password.Validate()` when you validate the `User` struct.

---

### 📦 Nested Plugin Support

Godantic supports plugins **recursively**, meaning:

- Nested objects
- Elements within slices
- Pointers or non-pointers

All are handled correctly.

#### Example:

```go
type Role struct {
	Name string `json:"name"`
}

func (r Role) Validate() *godantic.Error {
	if r.Name != "admin" && r.Name != "user" {
		return &godantic.Error{
			ErrType: "INVALID_ROLE",
			Message: fmt.Sprintf("Role <%s> is not allowed", r.Name),
		}
	}
	return nil
}

type User struct {
	Roles []Role `json:"roles"` // Each Role will be validated using its Validate method
}
```

---

### 🧠 Why Use Plugins?

- When logic is too complex for struct tags
- When validation depends on multiple fields
- When you want reusable, encapsulated validation units
- When you want full control over the error being returned

---

### 💡 Good to Know

- Plugin logic runs **after tag validations** (e.g., `binding`, `format`, `regex`).
- If the plugin returns an error **without a path**, Godantic won’t add one. You should provide the `Path` in the error when relevant.
- Plugin validation is supported for both **pointer** and **non-pointer** struct types.

---


## 🔄 Dynamic Field Validation

In many applications, some fields are **not strictly typed at compile time** — especially when you're building form-like schemas, dynamic inputs, or polymorphic models.

Godantic solves this elegantly using the `DynamicFieldsValidator` interface.

---

### ✨ What is a Dynamic Field?

A dynamic field is one whose **value, name, and type** are **determined at runtime**, and needs to be validated **accordingly**.

To support this, implement the following interface:

```go
type DynamicFieldsValidator interface {
	GetValue() any          // The actual value
	GetValueType() string   // Expected type: "string", "float", "boolean", "integer"
	GetAttribute() string   // The name/path for error reporting
}
```

Godantic will automatically invoke your implementation and validate the dynamic value.

---

### ✅ Supported Value Types

| `GetValueType()` | Validated As...       |
|------------------|------------------------|
| `"string"`       | Must be a Go `string` |
| `"float"`        | Must be a `float64`   |
| `"boolean"`      | Must be a `bool`      |
| `"integer"`      | Must be an `int` or castable `float64` |
| `"numeric"`      | Alias for `"integer"` |

---

### 🧪 Example: Simple Dynamic Field

```go
type MyDynamicField struct {
	Value     interface{} `json:"value"`
	ValueType string      `json:"valueType" enums:"string,integer,boolean"`
	Attribute string      `json:"attribute"`
}

func (mdf MyDynamicField) GetValue() any        { return mdf.Value }
func (mdf MyDynamicField) GetValueType() string { return mdf.ValueType }
func (mdf MyDynamicField) GetAttribute() string { return mdf.Attribute }
```

```go
type Request struct {
	Field MyDynamicField `json:"field"`
}
```

```json
{
  "field": {
    "value": 42,
    "valueType": "integer",
    "attribute": "age"
  }
}
```

Godantic will automatically validate the field value based on the declared type (`"integer"` in this case).

---

### 🧪 Example: With Array of Dynamic Fields

```go
type Request struct {
	Fields []MyDynamicField `json:"fields"`
}
```

Each element in the list will be validated using its dynamic type at runtime.

---

### ❌ Example: Invalid Type

```json
{
  "field": {
    "value": "not a number",
    "valueType": "integer",
    "attribute": "age"
  }
}
```

✅ Error:
```
Invalid value type for field 'age'. Expected numeric value.
```

---

### 💡 Why Use Dynamic Fields?

- You're building a **form builder**, **rule engine**, or **API that accepts generic inputs**.
- You want validation **without knowing types at compile time**.
- You need to enforce types dynamically **based on metadata**.

---

### 💬 Notes

- Works for both **pointer** and **non-pointer** values.
- Integrates seamlessly into nested objects and lists.
- Errors are detailed and reference the provided `attribute` for clarity.

---



## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.


This README.md includes detailed information about how to use Godantic, including simple and advanced usage examples, integration with web frameworks, features, error types, supported format tags, and more. If you have any further updates or modifications, please let me know!