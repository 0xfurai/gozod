# Gozod - Zod-like Validation Library for Go

**Gozod** is a powerful, type-safe validation library for Go inspired by [Zod](https://zod.dev), featuring clear error messages and a fluent API.

## Features

- 🎯 **Type-safe schemas** - Define validation rules with a fluent API
- 📝 **Clear error messages** - Human-readable error messages like Zod
- 🔗 **Chainable validators** - Build complex validation rules easily
- 🎨 **Multiple data types** - String, Int, Float, Bool, Object, Array
- ✅ **Nilable support** - Handle nilable fields
- 🗺️ **Map to Struct conversion** - Validate and convert API payloads to structs
- 🔍 **Flexible validation** - Validate maps, structs, or any any
- 🎨 **Custom error messages** - Easy customization of error messages
- 🔧 **Custom validation with Refine** - Add custom validation logic like Zod's refine
- 🎯 **Advanced validation with SuperRefine** - Fine-grained error control with custom paths and codes, similar to Zod's superRefine
- 📊 **Rich error handling** - Group, filter, and format errors for API responses
- 🚀 **Zero dependencies** - Pure Go, no external dependencies

## Installation

```bash
go get github.com/0xfurai/gozod
```

Or use it directly in your project:

```go
import "gozod"
```

## Quick Start

```go
package main

import (
    "encoding/json"
    "fmt"
    "gozod"
)

// Define your struct
type User struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
    IsActive bool   `json:"isActive"`
}

func main() {
    // Type-safe schemas - Define validation rules with a fluent API
    userSchema := gozod.Struct(map[string]gozod.Schema{
        "name": gozod.String().
            Min(3).
            Max(50),
        "email": gozod.String().
            Email(),
        "age": gozod.Int().
            Min(18).
            Max(120),
        "isActive": gozod.Bool().Nilable(),
    })

    // Valid user
    validUser := User{
        Name:     "John Doe",
        Email:    "john@example.com",
        Age:      25,
        IsActive: true,
    }

    errors := userSchema.Validate(validUser, nil)
    if errors != nil {
        fmt.Println("Validation failed:", errors.FormatErrors())
    } else {
        fmt.Println("✅ Valid user!")
    }

    // Invalid user (demonstrates error handling)
    invalidUser := User{
        Name:  "Jo",              // Too short (min 3)
        Email: "invalid-email",   // Invalid email format
        Age:   15,                // Too young (min 18)
    }

    errors = userSchema.Validate(invalidUser, nil)
    if errors != nil {
        // Get formatted errors for display
        fmt.Println("\n❌ Validation errors:")
        fmt.Println(errors.FormatErrors())

        // Or get JSON format for API responses
        jsonErrors, _ := json.MarshalIndent(errors.FormatErrorsJSON(), "", "  ")
        fmt.Println("\nJSON format:", string(jsonErrors))
    }
}
```

## Documentation

- 📚 [Examples](docs/examples.md) - Comprehensive validation examples for all types
- 📖 [API Reference](docs/api-reference.md) - Complete API documentation
- ⚠️ [Error Handling](docs/error-handling.md) - Error handling guide and customization
- 🎯 [Use Cases](docs/use-cases.md) - Real-world use case examples
- 🔄 [Comparison with Zod](docs/comparison.md) - Feature comparison with Zod

## Running Examples

```bash
# Run comprehensive examples
go run examples/main.go

# Or run your own code
go run main.go
```

## License

MIT
