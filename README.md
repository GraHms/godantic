# Godantic

Godantic is a Go package for inspecting and validating JSON-like data against Go struct types and schemas. It provides functionalities for checking type compatibility, structure compatibility, and other validations such as empty string, invalid time, and minimum length list checks.

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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
