# Examples

This guide provides comprehensive examples for all validation types supported by Go Validation.

## Table of Contents

- [String Validation](#string-validation)
- [Number Validation](#number-validation)
- [Object Validation](#object-validation)
- [Array Validation](#array-validation)
- [Nilable Fields](#nilable-fields)
- [Nested Objects](#nested-objects)
- [API Request Validation](#api-request-validation)
- [Validate Map Without Conversion](#validate-map-without-conversion)
- [Validate Struct Directly](#validate-struct-directly)
- [Custom Error Messages](#custom-error-messages)
- [Custom Validation with Refine](#custom-validation-with-refine)
- [Advanced Validation with SuperRefine](#advanced-validation-with-superrefine)

## String Validation

### Basic String Validation

```go
// Basic string validation
nameSchema := gozod.String().
    Min(3).
    Max(50)
```

### Email Validation

```go
// Email validation
emailSchema := gozod.String().
    Email().
    Min(5)
```

### URL Validation

```go
// URL validation
urlSchema := gozod.String().
    URL()
```

### Regex Validation

```go
// Regex validation
passwordSchema := gozod.String().
    Min(8).
    Regex(`[A-Z]`, "Must contain uppercase letter").
    Regex(`[0-9]`, "Must contain number")
```

### Enum Validation

```go
// Enum validation
statusSchema := gozod.String().
    OneOf("pending", "approved", "rejected")
```

### String Operations

```go
// String operations
domainSchema := gozod.String().
    StartsWith("https://").
    EndsWith(".com")
```

## Number Validation

### Integer Validation

```go
// Integer validation (only accepts integers)
ageSchema := gozod.Int().
    Min(18).
    Max(120)
```

### Float Validation

```go
// Float validation (only accepts floats, rejects integers)
priceSchema := gozod.Float().
    Positive().
    Max(1000.0)
```

### Advanced Number Validators

```go
// Even number validation
evenSchema := gozod.Int().
    MultipleOf(2)

// Positive number validation
positiveSchema := gozod.Int().
    Positive()

// Non-negative number validation
nonNegativeSchema := gozod.Int().
    NonNegative()
```

## Object Validation

```go
userSchema := gozod.Map(map[string]gozod.Schema{
    "name": gozod.String().
        Min(3).
        Max(50),
    "email": gozod.String().
        Email(),
    "age": gozod.Int().
        Min(18).
        Max(120),
    "isActive": gozod.Bool(),
})

// Validate a map
user := map[string]any{
    "name":     "John Doe",
    "email":    "john@example.com",
    "age":      25,
    "isActive": true,
}

errors := userSchema.Validate(user, nil)
if errors != nil {
    fmt.Println(errors.FormatErrors())
}
```

## Array Validation

```go
tagsSchema := gozod.Array(
    gozod.String().Min(2).Max(20),
).Min(1).Max(10)

tags := []any{"golang", "validation", "zod"}
errors := tagsSchema.Validate(tags, nil)
if errors != nil {
    fmt.Println(errors.FormatErrors())
}
```

## Nilable Fields

```go
profileSchema := gozod.Map(map[string]gozod.Schema{
    "name": gozod.String().Min(2),
    "bio":  gozod.String().Max(500).Nilable(),
    "age":  gozod.Int().Min(0).Nilable(),
})
```

## Nested Objects

```go
postSchema := gozod.Map(map[string]gozod.Schema{
    "title": gozod.String().Min(5),
    "author": gozod.Map(map[string]gozod.Schema{
        "name":  gozod.String().Min(2),
        "email": gozod.String().Email(),
    }),
    "tags": gozod.Array(
        gozod.String().Min(2),
    ).Nilable(),
})
```

## API Request Validation

Perfect for validating HTTP request payloads:

```go
// Define schema
userSchema := gozod.Map(map[string]gozod.Schema{
    "name":     gozod.String().Min(3).Max(50),
    "email":    gozod.String().Email(),
    "age":      gozod.Int().Min(18).Max(120),
    "isActive": gozod.Bool().Nilable(),
})

// Parse JSON payload (from HTTP request)
var payload map[string]any
json.Unmarshal(requestBody, &payload)

// Validate payload
errors := userSchema.Validate(payload, nil)
if errors != nil {
    // Return validation errors to client
    return errors.FormatErrorsJSON()
}

// Use validated payload
fmt.Printf("Payload: %+v\n", payload)
```

## Validate Map

When you need to validate a map:

```go
loginSchema := gozod.Map(map[string]gozod.Schema{
    "email":    gozod.String().Email(),
    "password": gozod.String().Min(8),
})

loginData := map[string]any{
    "email":    "user@example.com",
    "password": "secret123",
}

errors := loginSchema.Validate(loginData, nil)
if errors != nil {
    return errors.FormatErrorsJSON()
}
```

## Validate Struct Directly

Validate existing structs:

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}

user := User{
    Name:  "John Doe",
    Email: "john@example.com",
    Age:   25,
}

userSchema := gozod.Struct(gozod.Shape{
    "name":  gozod.String().Min(2),
    "email": gozod.String().Email(),
    "age":   gozod.Int().Min(18),
})

errors := userSchema.Validate(user, nil)
if errors != nil {
    return errors.FormatErrorsJSON()
}
```

## Custom Validation with Refine

The `Refine` method allows you to add custom validation logic that isn't covered by built-in validators. This is similar to Zod's `refine` method.

### String Refine

```go
// Custom validation: string must start with 'A'
schema := gozod.String().Refine(func(value any) (bool, string) {
    str := value.(string)
    if len(str) > 0 && str[0] == 'A' {
        return true, ""
    }
    return false, "String must start with 'A'"
})

// Multiple refinements can be chained
schema := gozod.String().
    Min(5).
    Refine(func(value any) (bool, string) {
        str := value.(string)
        return str != "forbidden", "This value is not allowed"
    }).
    Refine(func(value any) (bool, string) {
        str := value.(string)
        return len(str)%2 == 0, "String length must be even"
    })
```

### Number Refine

```go
// Validate that number is even
evenSchema := gozod.Int().Refine(func(value any) (bool, string) {
    var num int64
    switch v := value.(type) {
    case int:
        num = int64(v)
    case int64:
        num = v
    default:
        return false, "Invalid integer type"
    }
    return num%2 == 0, "Number must be even"
})

// Validate that number is a perfect square
perfectSquareSchema := gozod.Int().Refine(func(value any) (bool, string) {
    var num int64
    switch v := value.(type) {
    case int:
        num = int64(v)
    case int64:
        num = v
    default:
        return false, "Invalid integer type"
    }
    sqrt := int64(float64(num))
    return sqrt*sqrt == num, "Number must be a perfect square"
})
```

### Object/Map Refine (Cross-field Validation)

```go
// Password confirmation validation
passwordSchema := gozod.Map(map[string]gozod.Schema{
    "password": gozod.String().Min(8),
    "confirm":  gozod.String().Min(8),
}).Refine(func(value any) (bool, string) {
    m := value.(map[string]any)
    password, ok1 := m["password"].(string)
    confirm, ok2 := m["confirm"].(string)
    if !ok1 || !ok2 {
        return true, "" // Let type validation handle this
    }
    if password != confirm {
        return false, "Passwords do not match"
    }
    return true, ""
})

// Date range validation
dateRangeSchema := gozod.Map(map[string]gozod.Schema{
    "startDate": gozod.String(),
    "endDate":   gozod.String(),
}).Refine(func(value any) (bool, string) {
    m := value.(map[string]any)
    startDate, _ := m["startDate"].(string)
    endDate, _ := m["endDate"].(string)
    // Parse and compare dates (simplified example)
    if startDate > endDate {
        return false, "Start date must be before end date"
    }
    return true, ""
})

// Usage
errors := passwordSchema.Validate(map[string]any{
    "password": "password123",
    "confirm":  "different123",
}, nil)
if errors != nil {
    fmt.Println(errors.FormatErrors())
}
```

### Array Refine

```go
import (
    "reflect"
    "gozod"
)

// Validate that array contains unique values
uniqueArraySchema := gozod.Array(gozod.Int()).Refine(func(value any) (bool, string) {
    // Use reflect to handle different slice types
    val := reflect.ValueOf(value)
    if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
        return false, "Value is not a slice or array"
    }

    seen := make(map[int64]bool)
    for i := 0; i < val.Len(); i++ {
        elem := val.Index(i).Interface()
        var num int64
        switch n := elem.(type) {
        case int:
            num = int64(n)
        case int64:
            num = n
        default:
            return false, "Invalid integer type in array"
        }
        if seen[num] {
            return false, "Array must contain unique values"
        }
        seen[num] = true
    }
    return true, ""
})

// Usage
errors := uniqueArraySchema.Validate([]int{1, 2, 2}, nil)
if errors != nil {
    fmt.Println(errors.FormatErrors())
}
```

### Important Notes

- **Refine functions are only called after type validation passes.** If the value doesn't match the expected type, refine functions won't be executed.
- **Multiple refinements can be chained** - each refine function is executed in order.
- **Custom error messages** - You can provide a custom error message, or return an empty string to use the default message.
- **Type handling** - When working with numbers or arrays, you may need to handle different Go types (int, int64, etc.) in your refine function.

## Advanced Validation with SuperRefine

The `SuperRefine` method provides advanced validation capabilities similar to Zod's `superRefine`. Unlike `Refine`, which only allows you to return a boolean and message, `SuperRefine` gives you a context object that lets you add errors with custom paths, codes, and metadata. This is especially useful for:

- **Cross-field validation** - Add errors to specific fields
- **Multiple errors** - Add multiple errors from a single validation
- **Custom error paths** - Control exactly where errors appear
- **Error metadata** - Attach additional data to errors

### Basic SuperRefine

```go
// Simple superRefine with custom error
schema := gozod.String().SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    str := value.(string)
    if len(str) < 5 {
        ctx.AddIssue([]any{}, gozod.ErrCodeTooSmall, "String must be at least 5 characters")
    }
    if str == "forbidden" {
        ctx.AddIssue([]any{}, gozod.ErrCodeInvalidString, "This value is forbidden")
    }
})
```

### Cross-Field Validation with Custom Paths

The main advantage of `SuperRefine` is the ability to add errors to specific field paths, making it perfect for cross-field validation:

```go
// Password confirmation with error on specific field
passwordSchema := gozod.Map(map[string]gozod.Schema{
    "password": gozod.String().Min(8),
    "confirm":  gozod.String().Min(8),
}).SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    m := value.(map[string]any)
    password, ok1 := m["password"].(string)
    confirm, ok2 := m["confirm"].(string)
    if !ok1 || !ok2 {
        return // Let type validation handle this
    }
    if password != confirm {
        // Add error specifically to the "confirm" field
        ctx.AddIssue([]any{"confirm"}, gozod.ErrCodeCustomValidation, "Passwords do not match")
    }
})

// Usage
errors := passwordSchema.Validate(map[string]any{
    "password": "password123",
    "confirm":  "different123",
}, nil)
// Error will be on the "confirm" field path
```

### Multiple Errors from Single Validation

`SuperRefine` allows you to add multiple errors in a single validation pass:

```go
// Date range validation with multiple error conditions
dateRangeSchema := gozod.Map(map[string]gozod.Schema{
    "start": gozod.Int(),
    "end":   gozod.Int(),
}).SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    m := value.(map[string]any)
    start, ok1 := m["start"].(int64)
    end, ok2 := m["end"].(int64)
    if !ok1 || !ok2 {
        return
    }

    // Add error to "start" field if start >= end
    if start >= end {
        ctx.AddIssue([]any{"start"}, gozod.ErrCodeTooBig, "Start must be less than end")
    }

    // Add error to "end" field if range is too large
    if end-start > 100 {
        ctx.AddIssue([]any{"end"}, gozod.ErrCodeTooBig, "Range must not exceed 100")
    }
})
```

### Array Element Validation

You can add errors to specific array indices:

```go
import (
    "reflect"
    "gozod"
)

// Validate unique values with error on duplicate element
uniqueArraySchema := gozod.Array(gozod.Int()).SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    val := reflect.ValueOf(value)
    if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
        return
    }

    seen := make(map[int64]bool)
    for i := 0; i < val.Len(); i++ {
        elem := val.Index(i).Interface()
        var num int64
        switch n := elem.(type) {
        case int:
            num = int64(n)
        case int64:
            num = n
        default:
            continue
        }
        if seen[num] {
            // Add error to the specific array index where duplicate was found
            ctx.AddIssue([]any{i}, gozod.ErrCodeCustomValidation, "Duplicate value found")
            return
        }
        seen[num] = true
    }
})

// Usage
errors := uniqueArraySchema.Validate([]int{1, 2, 2}, nil)
// Error will be on path [2] (the duplicate element)
```

### Custom Error Codes

You can use custom error codes for better error categorization:

```go
// Custom error code for business logic validation
schema := gozod.Int().SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    num := value.(int64)
    if num%2 != 0 {
        ctx.AddIssue([]any{}, "odd_number", "Number must be even")
    }
})
```

### Error Metadata

Attach additional metadata to errors for richer error information:

```go
schema := gozod.String().SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    str := value.(string)
    if len(str) < 3 {
        meta := map[string]any{
            "minLength":    3,
            "actualLength": len(str),
        }
        ctx.AddIssueWithMeta([]any{}, gozod.ErrCodeTooSmall, "String too short", meta)
    }
})

// When validation fails, the error will include metadata
errors := schema.Validate("ab", nil)
// Error.Meta will contain {"minLength": 3, "actualLength": 2}
```

### Nested Path Errors

Add errors to deeply nested paths:

```go
schema := gozod.Map(map[string]gozod.Schema{
    "user": gozod.Map(map[string]gozod.Schema{
        "profile": gozod.Map(map[string]gozod.Schema{
            "age": gozod.Int(),
        }),
    }),
}).SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    m := value.(map[string]any)
    if user, ok := m["user"].(map[string]any); ok {
        if profile, ok := user["profile"].(map[string]any); ok {
            if age, ok := profile["age"].(int64); ok {
                if age < 18 {
                    // Add error to nested path: user.profile.age
                    ctx.AddIssue([]any{"user", "profile", "age"}, gozod.ErrCodeTooSmall, "Age must be at least 18")
                }
            }
        }
    }
})
```

### SuperRefine vs Refine

| Feature | Refine | SuperRefine |
|---------|--------|-------------|
| Error path control | ❌ Always on current field | ✅ Can specify any path |
| Multiple errors | ❌ One error per validation | ✅ Can add multiple errors |
| Error codes | ❌ Always `ErrCodeCustomValidation` | ✅ Can use any error code |
| Error metadata | ❌ Not supported | ✅ Supported via `AddIssueWithMeta` |
| Use case | Simple validation with custom message | Complex validation with fine-grained control |

**When to use Refine:**
- Simple custom validation
- Single error message needed
- Error should appear on the current field

**When to use SuperRefine:**
- Cross-field validation
- Need to add errors to specific fields
- Multiple errors from one validation
- Need custom error codes or metadata
- Complex business logic validation

## Custom Error Messages

Customize error messages for better user experience:

```go
// Per-field custom errors
emailSchema := gozod.String().
    Email().
    CustomError(gozod.ErrCodeInvalidString, "Please provide a valid email address").
    CustomError(gozod.ErrCodeRequired, "Email is required")

// Using error formatter for dynamic messages
ageSchema := gozod.Int().
    Min(18).
    Max(120).
    SetErrorFormatter(func(path, code, defaultMsg string) string {
        switch code {
        case gozod.ErrCodeTooSmall:
            return "You must be at least 18 years old"
        case gozod.ErrCodeTooBig:
            return "Age cannot exceed 120 years"
        default:
            return defaultMsg
        }
    })

// API request with custom errors
createUserSchema := gozod.Map(map[string]gozod.Schema{
    "name": gozod.String().
        Min(3).
        CustomError(gozod.ErrCodeTooSmall, "Name must be at least 3 characters").
        CustomError(gozod.ErrCodeRequired, "Name is required"),
    "email": gozod.String().
        Email().
        CustomError(gozod.ErrCodeInvalidString, "Invalid email format"),
})
```

## See Also

- [API Reference](api-reference.md) - Complete API documentation
- [Error Handling](error-handling.md) - Learn more about error handling
- [Use Cases](use-cases.md) - Real-world examples
