package gozod

import (
	"fmt"
	"regexp"
	"strings"
)

// Error codes for validation errors
const (
	// ErrCodeRequired indicates a required field is missing
	ErrCodeRequired = "required"

	// ErrCodeInvalidType indicates the value type doesn't match the expected type
	ErrCodeInvalidType = "invalid_type"

	// ErrCodeTooSmall indicates a value is too small (below minimum)
	ErrCodeTooSmall = "too_small"

	// ErrCodeTooBig indicates a value is too big (above maximum)
	ErrCodeTooBig = "too_big"

	// ErrCodeInvalidString indicates a string validation failed (email, URL, regex, etc.)
	ErrCodeInvalidString = "invalid_string"

	// ErrCodeInvalidEnumValue indicates a value is not in the allowed enum values
	ErrCodeInvalidEnumValue = "invalid_enum_value"

	// ErrCodeUnrecognizedKeys indicates unknown keys in strict object validation
	ErrCodeUnrecognizedKeys = "unrecognized_keys"

	// ErrCodeCustomValidation indicates a custom refine validation failed
	ErrCodeCustomValidation = "custom_validation"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Path    []any          // Field path as array of path parts (e.g., ["user", "email"] or ["test", 1])
	Message string         // Human-readable error message
	Code    string         // Error code (use constants like ErrCodeTooSmall, ErrCodeInvalidType, etc.)
	Meta    map[string]any // Additional metadata for the error
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return e.Message
}

// ValidationErrors is a collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError
}

// Error implements the error interface
func (e *ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Message
	}

	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// Add adds a new validation error
func (e *ValidationErrors) Add(path []any, code, message string) {
	e.AddWithMeta(path, code, message, nil)
}

// AddWithMeta adds a new validation error with metadata
func (e *ValidationErrors) AddWithMeta(path []any, code, message string, meta map[string]any) {
	// Make a copy of the path to avoid mutations
	pathCopy := make([]any, len(path))
	copy(pathCopy, path)
	e.Errors = append(e.Errors, ValidationError{
		Path:    pathCopy,
		Code:    code,
		Message: message,
		Meta:    meta,
	})
}

// FormatErrors returns a formatted string of all errors
func (e *ValidationErrors) FormatErrors() string {
	if len(e.Errors) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Validation failed with %d error(s):\n", len(e.Errors)))
	for i, err := range e.Errors {
		pathStr := PathToString(err.Path)
		builder.WriteString(fmt.Sprintf("  %d. [%s] %s\n", i+1, pathStr, err.Message))
	}
	return builder.String()
}

// FormatErrorsJSON returns errors in a structured format (for API responses)
func (e *ValidationErrors) FormatErrorsJSON() map[string]any {
	if len(e.Errors) == 0 {
		return nil
	}

	errors := make([]map[string]any, len(e.Errors))
	for i, err := range e.Errors {
		errorMap := map[string]any{
			"path":    err.Path,
			"code":    err.Code,
			"message": err.Message,
		}
		if err.Meta != nil && len(err.Meta) > 0 {
			errorMap["meta"] = err.Meta
		}
		errors[i] = errorMap
	}

	return map[string]any{
		"errors": errors,
		"count":  len(e.Errors),
	}
}

// GetErrorsByPath returns all errors for a specific path
func (e *ValidationErrors) GetErrorsByPath(path []any) []ValidationError {
	var result []ValidationError
	for _, err := range e.Errors {
		if PathEqual(err.Path, path) {
			result = append(result, err)
		}
	}
	return result
}

// GetErrorsByCode returns all errors with a specific code
func (e *ValidationErrors) GetErrorsByCode(code string) []ValidationError {
	var result []ValidationError
	for _, err := range e.Errors {
		if err.Code == code {
			result = append(result, err)
		}
	}
	return result
}

// FlattenErrorResult represents the flattened error structure
type FlattenErrorResult struct {
	FormErrors  []string            `json:"formErrors"`
	FieldErrors map[string][]string `json:"fieldErrors"`
}

// Flatten returns errors in a flattened format with formErrors and fieldErrors
// formErrors contains errors without a specific field path (empty path)
// fieldErrors contains errors grouped by field path
// Array element errors (e.g., test[0], test[1]) are grouped under the base field name (e.g., test)
func (e *ValidationErrors) Flatten() FlattenErrorResult {
	result := FlattenErrorResult{
		FormErrors:  []string{},
		FieldErrors: make(map[string][]string),
	}

	// Regex to match array indices like [0], [1], etc.
	arrayIndexRegex := regexp.MustCompile(`\[(\d+)\]`)

	for _, err := range e.Errors {
		pathStr := PathToString(err.Path)
		if pathStr == "" {
			// Empty path means form-level error
			result.FormErrors = append(result.FormErrors, err.Message)
		} else {
			// Check if path contains array indices
			if arrayIndexRegex.MatchString(pathStr) {
				// Extract base field name (everything before the first [)
				baseField := arrayIndexRegex.Split(pathStr, 2)[0]
				// Group all array element errors under the base field name
				result.FieldErrors[baseField] = append(result.FieldErrors[baseField], err.Message)
			} else {
				// Regular field-level error
				result.FieldErrors[pathStr] = append(result.FieldErrors[pathStr], err.Message)
			}
		}
	}

	return result
}

// PathToString converts a path array to a string representation
// e.g., ["user", "email"] -> "user.email", ["test", 1] -> "test[1]"
func PathToString(path []any) string {
	if len(path) == 0 {
		return ""
	}
	var parts []string
	for _, part := range path {
		switch v := part.(type) {
		case string:
			parts = append(parts, v)
		case int:
			parts = append(parts, fmt.Sprintf("[%d]", v))
		case int64:
			parts = append(parts, fmt.Sprintf("[%d]", v))
		default:
			parts = append(parts, fmt.Sprintf("[%v]", v))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], "[") {
			result += parts[i]
		} else {
			result += "." + parts[i]
		}
	}
	return result
}

// PathEqual checks if two path arrays are equal
func PathEqual(path1, path2 []any) bool {
	if len(path1) != len(path2) {
		return false
	}
	for i := range path1 {
		if path1[i] != path2[i] {
			return false
		}
	}
	return true
}

// PathAppend appends a new part to a path array and returns a new array
func PathAppend(path []any, part any) []any {
	newPath := make([]any, len(path), len(path)+1)
	copy(newPath, path)
	return append(newPath, part)
}
