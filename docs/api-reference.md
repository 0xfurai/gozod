# API Reference

Complete API documentation for Go Validation library.

## Table of Contents

- [Core Functions](#core-functions)
- [String Schema](#string-schema)
- [Number Schema](#number-schema)
- [Object Schema](#object-schema)
- [Array Schema](#array-schema)
- [Boolean Schema](#boolean-schema)

## Core Functions

### Validate

All schemas implement the `Validate` method which validates a value against the schema and returns validation errors.

```go
func (schema Schema) Validate(value any, path []any) *ValidationErrors
```

**Parameters:**
- `value` - The value to validate
- `path` - The path for error reporting (use `nil` for top-level validation)

**Returns:**
- `*ValidationErrors` - Validation errors if validation fails, `nil` if valid

**Examples:**

Validate a string:
```go
emailSchema := gozod.String().Email()
errors := emailSchema.Validate("user@example.com", nil)
if errors != nil {
    fmt.Println(errors.FormatErrors())
}
```

Validate a map:
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
```

Validate a struct directly:
```go
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
```

## String Schema

### String

Create a new string schema. Fields are required by default unless marked as nilable.

```go
func String() *StringSchema
```

### Nilable

Allow null/nil values for this field.

```go
func (s *StringSchema) Nilable() *StringSchema
```

### Min

Set minimum length requirement.

```go
func (s *StringSchema) Min(length int) *StringSchema
```

### Max

Set maximum length requirement.

```go
func (s *StringSchema) Max(length int) *StringSchema
```

### Email

Validate email format.

```go
func (s *StringSchema) Email() *StringSchema
```

### URL

Validate URL format.

```go
func (s *StringSchema) URL() *StringSchema
```

### Regex

Validate against a regular expression pattern.

```go
func (s *StringSchema) Regex(pattern string, message ...string) *StringSchema
```

**Parameters:**
- `pattern` - Regular expression pattern
- `message` - Optional custom error message

### OneOf

Value must be one of the provided options.

```go
func (s *StringSchema) OneOf(options ...string) *StringSchema
```

### NotOneOf

Value must not be one of the provided options.

```go
func (s *StringSchema) NotOneOf(options ...string) *StringSchema
```

### StartsWith

String must start with the given prefix.

```go
func (s *StringSchema) StartsWith(prefix string) *StringSchema
```

### EndsWith

String must end with the given suffix.

```go
func (s *StringSchema) EndsWith(suffix string) *StringSchema
```

### Includes

String must include the given substring.

```go
func (s *StringSchema) Includes(substring string) *StringSchema
```

### CustomError

Set a custom error message for a specific error code.

```go
func (s *StringSchema) CustomError(code, message string) *StringSchema
```

### SetErrorFormatter

Set a function to format all errors dynamically.

```go
func (s *StringSchema) SetErrorFormatter(formatter CustomErrorFunc) *StringSchema
```

### Refine

Add a custom validation function. This allows you to implement any custom validation logic that isn't covered by the built-in validators.

```go
func (s *StringSchema) Refine(validator RefineFunc) *StringSchema
```

**Parameters:**
- `validator` - A function that receives the value and returns `(bool, string)` where:
  - First return value (`bool`): `true` if validation passes, `false` if it fails
  - Second return value (`string`): Custom error message (optional, can be empty string for default message)

**Example:**
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

**Note:** Refine functions are only called after type validation passes. If the value doesn't match the expected type, refine functions won't be executed.

### SuperRefine

Add an advanced custom validation function with fine-grained control over error reporting. Similar to Zod's `superRefine`, this method provides a context object that allows you to add errors with custom paths, codes, and metadata.

```go
func (s *StringSchema) SuperRefine(validator SuperRefineFunc) *StringSchema
```

**Parameters:**
- `validator` - A function that receives the value and a `SuperRefineContext` object

**SuperRefineContext Methods:**
- `AddIssue(path []any, code, message string)` - Add a validation error with custom path and code
- `AddIssueWithMeta(path []any, code, message string, meta map[string]any)` - Add a validation error with metadata

**Example:**
```go
// Basic superRefine
schema := gozod.String().SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    str := value.(string)
    if len(str) < 5 {
        ctx.AddIssue([]any{}, gozod.ErrCodeTooSmall, "String must be at least 5 characters")
    }
})

// Multiple errors
schema := gozod.String().SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    str := value.(string)
    if len(str) < 3 {
        ctx.AddIssue([]any{}, gozod.ErrCodeTooSmall, "Too short")
    }
    if str == "forbidden" {
        ctx.AddIssue([]any{}, gozod.ErrCodeInvalidString, "Forbidden value")
    }
})

// With metadata
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
```

**Note:** SuperRefine functions are only called after type validation passes. The context's `AddIssue` methods allow you to add errors to any path, making it ideal for cross-field validation and complex business logic.

## Number Schema

### Int

Create an integer schema. Only accepts integers, rejects floats.

```go
func Int() *IntSchema
```

### Float

Create a float schema. Only accepts floats, rejects integers.

```go
func Float() *FloatSchema
```

### Min

Set minimum value requirement.

```go
func (s *IntSchema) Min(value float64) *IntSchema
func (s *FloatSchema) Min(value float64) *FloatSchema
```

### Max

Set maximum value requirement.

```go
func (s *IntSchema) Max(value float64) *IntSchema
func (s *FloatSchema) Max(value float64) *FloatSchema
```

### Positive

Value must be positive (> 0).

```go
func (s *IntSchema) Positive() *IntSchema
func (s *FloatSchema) Positive() *FloatSchema
```

### Negative

Value must be negative (< 0).

```go
func (s *IntSchema) Negative() *IntSchema
func (s *FloatSchema) Negative() *FloatSchema
```

### NonNegative

Value must be non-negative (>= 0).

```go
func (s *IntSchema) NonNegative() *IntSchema
func (s *FloatSchema) NonNegative() *FloatSchema
```

### NonPositive

Value must be non-positive (<= 0).

```go
func (s *IntSchema) NonPositive() *IntSchema
func (s *FloatSchema) NonPositive() *FloatSchema
```

### MultipleOf

Value must be a multiple of the given value.

```go
func (s *IntSchema) MultipleOf(value float64) *IntSchema
func (s *FloatSchema) MultipleOf(value float64) *FloatSchema
```

### Nilable

Allow null/nil values for this field.

```go
func (s *IntSchema) Nilable() *IntSchema
func (s *FloatSchema) Nilable() *FloatSchema
```

### CustomError

Set a custom error message for a specific error code.

```go
func (s *IntSchema) CustomError(code, message string) *IntSchema
func (s *FloatSchema) CustomError(code, message string) *FloatSchema
```

### SetErrorFormatter

Set a function to format all errors dynamically.

```go
func (s *IntSchema) SetErrorFormatter(formatter CustomErrorFunc) *IntSchema
func (s *FloatSchema) SetErrorFormatter(formatter CustomErrorFunc) *FloatSchema
```

### Refine

Add a custom validation function for numbers.

```go
func (s *IntSchema) Refine(validator RefineFunc) *IntSchema
func (s *FloatSchema) Refine(validator RefineFunc) *FloatSchema
```

**Example:**
```go
// Validate that number is even
evenSchema := gozod.Int().Refine(func(value any) (bool, string) {
    // Handle different integer types
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

### SuperRefine

Add an advanced custom validation function for numbers with fine-grained error control.

```go
func (s *IntSchema) SuperRefine(validator SuperRefineFunc) *IntSchema
func (s *FloatSchema) SuperRefine(validator SuperRefineFunc) *FloatSchema
```

**Example:**
```go
// Custom error code
schema := gozod.Int().SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    num := value.(int64)
    if num%2 != 0 {
        ctx.AddIssue([]any{}, "odd_number", "Number must be even")
    }
})

// Multiple conditions
schema := gozod.Int().SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    num := value.(int64)
    if num < 0 {
        ctx.AddIssue([]any{}, gozod.ErrCodeTooSmall, "Must be non-negative")
    }
    if num > 1000 {
        ctx.AddIssue([]any{}, gozod.ErrCodeTooBig, "Must not exceed 1000")
    }
})
```

## Object Schema

### Map

Create an object/map schema. Extra keys are allowed by default.

```go
func Map(shape map[string]Schema) *MapSchema
```

**Parameters:**
- `shape` - Map of field names to their schemas

**Example:**
```go
userSchema := gozod.Map(map[string]gozod.Schema{
    "name": gozod.String().Min(3),
    "email": gozod.String().Email(),
})
```

### Strict

Reject unknown keys that are not defined in the schema.

```go
func (s *MapSchema) Strict() *MapSchema
```

### Nilable

Allow null/nil values for this field.

```go
func (s *MapSchema) Nilable() *MapSchema
```

### CustomError

Set a custom error message for a specific error code.

```go
func (s *MapSchema) CustomError(code, message string) *MapSchema
```

### SetErrorFormatter

Set a function to format all errors dynamically.

```go
func (s *MapSchema) SetErrorFormatter(formatter CustomErrorFunc) *MapSchema
```

### Refine

Add a custom validation function for objects/maps. This is particularly useful for cross-field validation.

```go
func (s *MapSchema) Refine(validator RefineFunc) *MapSchema
```

**Example:**
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
```

## Array Schema

### Array

Create an array schema with element validation.

```go
func Array(elementSchema Schema) *ArraySchema
```

**Parameters:**
- `elementSchema` - Schema for validating each array element

**Example:**
```go
tagsSchema := gozod.Array(gozod.String().Min(2))
```

### Min

Set minimum array length requirement.

```go
func (s *ArraySchema) Min(length int) *ArraySchema
```

### Max

Set maximum array length requirement.

```go
func (s *ArraySchema) Max(length int) *ArraySchema
```

### NonEmpty

Array must not be empty.

```go
func (s *ArraySchema) NonEmpty() *ArraySchema
```

### Nilable

Allow null/nil values for this field.

```go
func (s *ArraySchema) Nilable() *ArraySchema
```

### CustomError

Set a custom error message for a specific error code.

```go
func (s *ArraySchema) CustomError(code, message string) *ArraySchema
```

### SetErrorFormatter

Set a function to format all errors dynamically.

```go
func (s *ArraySchema) SetErrorFormatter(formatter CustomErrorFunc) *ArraySchema
```

### Refine

Add a custom validation function for arrays.

```go
func (s *ArraySchema) Refine(validator RefineFunc) *ArraySchema
```

**Example:**
```go
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
```

### SuperRefine

Add an advanced custom validation function for arrays with fine-grained error control. Useful for adding errors to specific array indices.

```go
func (s *ArraySchema) SuperRefine(validator SuperRefineFunc) *ArraySchema
```

**Example:**
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
            // Error appears on the specific array index
            ctx.AddIssue([]any{i}, gozod.ErrCodeCustomValidation, "Duplicate value found")
            return
        }
        seen[num] = true
    }
})
```

## Boolean Schema

### Bool

Create a boolean schema.

```go
func Bool() *BoolSchema
```

### Nilable

Allow null/nil values for this field.

```go
func (s *BoolSchema) Nilable() *BoolSchema
```

### CustomError

Set a custom error message for a specific error code.

```go
func (s *BoolSchema) CustomError(code, message string) *BoolSchema
```

### SetErrorFormatter

Set a function to format all errors dynamically.

```go
func (s *BoolSchema) SetErrorFormatter(formatter CustomErrorFunc) *BoolSchema
```

### Refine

Add a custom validation function for booleans.

```go
func (s *BoolSchema) Refine(validator RefineFunc) *BoolSchema
```

**Example:**
```go
// Custom boolean validation
schema := gozod.Bool().Refine(func(value any) (bool, string) {
    b := value.(bool)
    // Example: require true for certain conditions
    if !b {
        return false, "This field must be true"
    }
    return true, ""
})
```

### SuperRefine

Add an advanced custom validation function for booleans with fine-grained error control.

```go
func (s *BoolSchema) SuperRefine(validator SuperRefineFunc) *BoolSchema
```

**Example:**
```go
// Custom validation with metadata
schema := gozod.Bool().SuperRefine(func(value any, ctx *gozod.SuperRefineContext) {
    b := value.(bool)
    if !b {
        meta := map[string]any{
            "reason": "Terms must be accepted",
        }
        ctx.AddIssueWithMeta([]any{}, gozod.ErrCodeCustomValidation, "This field must be true", meta)
    }
})
```

## See Also

- [Examples](examples.md) - Comprehensive validation examples
- [Error Handling](error-handling.md) - Error handling guide
- [Use Cases](use-cases.md) - Real-world examples
