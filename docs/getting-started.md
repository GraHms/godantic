# Getting Started

## Installation

```sh
go get github.com/grahms/godantic
```

## Basic Usage

```go
import "github.com/grahms/godantic"

var v godantic.Validate
err := v.BindJSON(jsonData, &myStruct)
```

This will bind and validate the JSON into your struct.