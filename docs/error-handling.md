# Error Handling

Comprehensive guide to error handling in Go Validation library.

## Table of Contents

- [Error Codes](#error-codes)
- [Error Structure](#error-structure)
- [Error Methods](#error-methods)
- [Custom Error Messages](#custom-error-messages)
- [Error Formatting](#error-formatting)
- [Advanced Error Handling](#advanced-error-handling)
- [Flatten Errors for Form Validation](#flatten-errors-for-form-validation)

## Error Codes

All error codes are available as constants for type-safe usage:

```go
gozod.ErrCodeRequired          // "required"
gozod.ErrCodeInvalidType       // "invalid_type"
gozod.ErrCodeTooSmall          // "too_small"
gozod.ErrCodeTooBig            // "too_big"
gozod.ErrCodeInvalidString     // "invalid_string"
gozod.ErrCodeInvalidEnumValue  // "invalid_enum_value"
gozod.ErrCodeUnrecognizedKeys  // "unrecognized_keys"
```

## Error Structure

Errors are returned as `*ValidationErrors` which provides:

```go
type ValidationError struct {
    Path    string // Field path (e.g., "user.email", "items[0].name")
    Message string // Human-readable error message
    Code    string // Error code (use constants like ErrCodeTooSmall, ErrCodeInvalidType, etc.)
}

type ValidationErrors struct {
    Errors []ValidationError
}
```

## Error Methods

### FormatErrors

Get formatted error string for logging or display.

```go
func (ve *ValidationErrors) FormatErrors() string
```

**Example:**
```go
errors := schema.Validate(data, nil)
if errors != nil {
    fmt.Println(errors.FormatErrors())
}
```

**Output:**
```
Validation failed with 3 error(s):
  1. [email] Invalid email format at email
  2. [age] Number at age must be greater than or equal to 18, got 15
  3. [name] String at name must be at least 3 character(s) long, got 2
```

### FormatErrorsJSON

Get errors as JSON structure (perfect for API responses).

```go
func (ve *ValidationErrors) FormatErrorsJSON() map[string]any
```

**Returns:**
```json
{
  "errors": [
    {
      "path": "email",
      "message": "Invalid email format at email",
      "code": "invalid_string"
    },
    {
      "path": "age",
      "message": "Number at age must be greater than or equal to 18, got 15",
      "code": "too_small"
    }
  ],
  "count": 2
}
```

**Example:**
```go
errors := schema.Validate(data, nil)
if errors != nil {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(errors.FormatErrorsJSON())
}
```

### GetErrorsByPath

Get all errors for a specific field path.

```go
func (ve *ValidationErrors) GetErrorsByPath(path string) []ValidationError
```

**Example:**
```go
errors := schema.Validate(data, nil)
if errors != nil {
    emailErrors := errors.GetErrorsByPath("email")
    for _, err := range emailErrors {
        fmt.Printf("Email error: %s\n", err.Message)
    }
}
```

### GetErrorsByCode

Get all errors with a specific error code.

```go
func (ve *ValidationErrors) GetErrorsByCode(code string) []ValidationError
```

**Example:**
```go
errors := schema.Validate(data, nil)
if errors != nil {
    tooSmallErrors := errors.GetErrorsByCode(gozod.ErrCodeTooSmall)
    for _, err := range tooSmallErrors {
        fmt.Printf("Too small: %s\n", err.Message)
    }
}
```

### Flatten

Flatten errors into formErrors and fieldErrors structure.

```go
func (ve *ValidationErrors) Flatten() FlattenErrorResult
```

**Returns:**
```go
type FlattenErrorResult struct {
    FormErrors  []string            // Errors without a specific field path
    FieldErrors map[string][]string // Errors grouped by field path
}
```

**Example:**
```go
errors := schema.Validate(data, nil)
if errors != nil {
    flattened := errors.Flatten()
    // Returns: FlattenErrorResult{
    //   FormErrors: []string{},
    //   FieldErrors: map[string][]string{
    //     "name": []string{"String must contain at most 1 character(s)"},
    //     "test": []string{"Array must contain at most 1 element(s)", "String must contain at most 2 character(s)"},
    //   },
    // }
}
```

## Custom Error Messages

### Per-Field Custom Errors

Set custom error messages for specific error codes:

```go
emailSchema := gozod.String().
    Email().
    CustomError(gozod.ErrCodeInvalidString, "Please provide a valid email address").
    CustomError(gozod.ErrCodeRequired, "Email is required")
```

### Dynamic Error Formatters

Use a function to format all errors dynamically:

```go
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
```

### API Request with Custom Errors

```go
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

## Error Formatting

### Example Error Output

```
Validation failed with 3 error(s):
  1. [email] Invalid email format at email
  2. [age] Number at age must be greater than or equal to 18, got 15
  3. [name] String at name must be at least 3 character(s) long, got 2
```

### JSON Error Output

```json
{
  "errors": [
    {
      "path": "email",
      "message": "Invalid email format at email",
      "code": "invalid_string"
    },
    {
      "path": "age",
      "message": "Number at age must be greater than or equal to 18, got 15",
      "code": "too_small"
    },
    {
      "path": "name",
      "message": "String at name must be at least 3 character(s) long, got 2",
      "code": "too_small"
    }
  ],
  "count": 3
}
```

## Advanced Error Handling

```go
errors := schema.Validate(data, nil)

// Check for errors
if errors != nil {
    // Get JSON format for API response
    jsonErrors := errors.FormatErrorsJSON()
    // Returns: {"errors": [...], "count": 3}

    // Get specific error types
    tooSmallErrors := errors.GetErrorsByCode(gozod.ErrCodeTooSmall)
    emailErrors := errors.GetErrorsByPath("email")

    // Format for logging
    fmt.Println(errors.FormatErrors())

    // Flatten errors for form validation
    flattened := errors.Flatten()
    // Returns: FlattenErrorResult{
    //   FormErrors: []string{},
    //   FieldErrors: map[string][]string{
    //     "name": []string{"String must contain at most 1 character(s)"},
    //     "test": []string{"Array must contain at most 1 element(s)", "String must contain at most 2 character(s)"},
    //   },
    // }
}
```

## Flatten Errors for Form Validation

The `Flatten()` method is perfect for form validation where you need to separate form-level errors from field-level errors:

```go
errors := schema.Validate(data, nil)

if errors != nil {
    flattened := errors.Flatten()

    // flattened.FormErrors contains errors without a specific field path
    // flattened.FieldErrors contains errors grouped by field path

    // Example output:
    // {
    //   "formErrors": [],
    //   "fieldErrors": {
    //     "name": ["String must contain at most 1 character(s)"],
    //     "test": [
    //       "Array must contain at most 1 element(s)",
    //       "String must contain at most 2 character(s)"
    //     ]
    //   }
    // }

    // Perfect for JSON API responses
    jsonBytes, _ := json.Marshal(flattened)
    w.Write(jsonBytes)
}
```

## See Also

- [Examples](examples.md) - See error handling in action
- [API Reference](api-reference.md) - Complete API documentation
- [Use Cases](use-cases.md) - Real-world error handling examples
