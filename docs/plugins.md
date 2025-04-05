# Plugin-Based Validation

## Interface

```go
type ValidationPlugin interface {
    Validate() *godantic.Error
}
```

Godantic will auto-detect and run this method on nested structs or list elements.