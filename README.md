# godantic

godantic is a Go package that provides functionality for decoding JSON data and validating it against a given object structure. It aims to simplify the process of decoding and validating JSON input in Go applications.

## Installation

To use godantic in your Go project, you need to have Go installed and set up on your machine. Then, you can use the go get command to install the package:

```shell
go get github.com/grahms/godantic
```
# Usage

```go
package main

import (
	"fmt"
	"github.com/grahms/godantic"
)

type User struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Age   int    `json:"age"`
}

func main() {
	jsonData := []byte(`{"name": "John Doe", "email": "john@example.com", "age": 30}`)

	var user User
	validator := godantic.Validate{}
	err := validator.BindJSON(jsonData, &user)
	if err != nil {
		fmt.Println("Validation Error:", err)
		return
	}

	fmt.Println("Validated user:", user)
}
```
In the example above, the User struct represents the expected structure of the JSON data. The Validate struct from the godantic package is used to bind and validate the JSON data against the User object. If any validation errors occur, an error is returned. Otherwise, the user object is populated with the validated data.

Make sure to import the "github.com/grahms/godantic" package in your code before using it.
## Advanced Usage

The `godantic` package provides advanced features to handle complex JSON structures, nested objects, and object lists. It ensures type compatibility, required fields, and unknown field validations. This section explains how to utilize these advanced features effectively.

### Error Messages and Handling

When validation errors occur, the package returns an `Error` object that provides detailed information about the error. The `Error` object contains the following fields:

- `ErrType`: The type of the error, which can be used for programmatic handling.
- `Message`: A descriptive error message that provides additional information about the error.
- `Path`: The error path that indicates the location of the error within the JSON structure.

To handle validation errors, you can use a type assertion to check if the error is of type `*godantic.Error`. If it is, you can access the specific fields of the `Error` object for further processing or error reporting.

Here's an example of handling validation errors:

```go
var user User
validator := godantic.Validate{}
err := validator.BindJSON(jsonData, &user)
if err != nil {
	switch e := err.(type) {
	case *godantic.Error:
		fmt.Printf("Validation Error: %s\n", e.Error())
		fmt.Printf("Error Type: %s\n", e.ErrType)
		fmt.Printf("Path: %s\n", e.Path)
	default:
		fmt.Println("Unexpected error:", err)
	}
	return
}
```

In the example above, the error is checked if it is of type `*godantic.Error`. If it is, the error message, error type, and error path are printed to the console. For unexpected errors, a generic error message is printed.

### Nested Paths

When a validation error occurs within nested objects or object lists, the error path reflects the nested structure of the JSON. Each level of nesting is separated by a dot (`.`) in the error path.

For example, if an error occurs in the `number` field of the second phone object, the error path will be `phones[1].number`. This indicates that the error is located within the `number` field of the second phone object in the `phones` array.

Similarly, if an extra unknown field `foo` is provided within the second phone object, the error path will be `phones[1].foo`. This helps identify the exact location of the error within the JSON structure.

By examining the error paths, you can easily pinpoint the specific fields or objects that have validation issues and take appropriate actions based on the error path.

### Customizing Validation Rules

The `godantic` package allows you to define custom validation rules using struct tags. You can utilize struct tags to specify required fields, define custom validation rules, or provide additional metadata for fields.

Here's an example of how to define custom validation rules using struct tags:

```go
type User struct {
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Age       int       `json:"age" binding:"gte=18"`
	CreatedAt time.Time `json:"created_at"`
	Address   Address   `json:"address" binding:"required"`
	Phones    []Phone   `json:"phones" binding:"required,dive"`
}
```

In the example above, the `binding` struct tag is used to define validation rules for each field:

- The `required` rule indicates that the field must be present and non-empty.
- The `email` rule specifies that the `Email` field must be a valid email address.
- The `gte=18` rule specifies that the `Age` field must be greater than or equal to 18.
- The `d

ive` rule is used for the `Phones` field to apply further validation to each element in the array.

By defining custom validation rules using struct tags, you can enforce specific constraints on your JSON structure and ensure data integrity.

Feel free to experiment with different struct tags and validation rules to suit your specific requirements and data validation needs.

---

By utilizing the advanced features of the `godantic` package, you can confidently validate complex JSON structures, handle validation errors effectively, and customize validation rules to ensure the integrity of your data.
## Documentation
For detailed information on how to use godantic, please refer to the GoDoc documentation.

## Contributing
Contributions to godantic are welcome! If you find any issues or have suggestions for improvements, please open an issue on the GitHub repository. Pull requests are also appreciated.

Before contributing, please read the contribution guidelines for this project.

## License
This project is licensed under the MIT License.
