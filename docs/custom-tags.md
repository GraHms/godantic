# Custom Tag Functions

Use the `validate` tag to attach custom logic to a field:

```go
type User struct {
  Name string `validate:"starts_with_A,min_len_3"`
}
```

Register:

```go
godantic.RegisterCustom[string]("starts_with_A", func(val string, path string) *godantic.Error {
    if !strings.HasPrefix(val, "A") {
        return &godantic.Error{...}
    }
    return nil
})
```