# Dynamic Field Validation

## Interface

```go
type DynamicFieldsValidator interface {
    GetValue() any
    GetValueType() string
    GetAttribute() string
}
```

Godantic will dynamically validate based on the value type (`string`, `integer`, `float`, etc).