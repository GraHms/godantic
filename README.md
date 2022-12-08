# Godantic
Godantic is a Go package for validating JSON input data in HTTP requests. It provides a simple way to validate that the input data conforms to a given JSON schema, and to ensure that the input data does not contain any extra fields that are not present in the schema.

## Usage
To use Godantic, you first need to create a Godantic instance:

```Go
g := &godantic.NewGodantic(".")
ErrorPathSeparator: , // The separator used to join field names in error messages
}
```
Next, you can use the ForbidExtra method to check if the input data contains any extra fields that are not present in the reference data:

```Go
referenceData := map[string]interface{}{
"name": "John Doe",
"age":  30,
}

inputData := map[string]interface{}{
"name": "John Doe",
"age":  30,
"city": "New York", // This field is not present in the reference data
}

err := g.Fields.ForbidExtra(inputData, referenceData, "")
if err != nil {
// Handle the error
}

```

If the input data contains any extra fields, ForbidExtra will return an error with the field names joined together with the separator specified in the ErrorPathSeparator field of the Godantic instance.

You can also use the ReflectStruct method to validate that the input data conforms to a given JSON schema:

```Go
type User struct {
Name string `json:"name" binding:"required"`
Age  int    `json:"age"`
}

inputData := map[string]interface{}{
"name": "John Doe",
"age":  "30", // This field should be an int, not a string
}

var user User
err := g.Fields.ValidateStruct(inputData, "")
if err != nil {
// Handle the error
}
```
If the input data is not a valid JSON representation of the given struct, ReflectStruct will return an error with the field names and types.

## License
Godantic is released under the MIT license. See LICENSE for more details.