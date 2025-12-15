package gozod

import (
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Path:    []any{"user", "email"},
		Message: "Invalid email",
		Code:    ErrCodeInvalidString,
	}

	if err.Error() != "Invalid email" {
		t.Errorf("Expected error message 'Invalid email', got '%s'", err.Error())
	}
}

func TestValidationErrors_Error(t *testing.T) {
	// Empty errors
	errors := &ValidationErrors{}
	if errors.Error() != "" {
		t.Errorf("Expected empty string for no errors, got '%s'", errors.Error())
	}

	// Single error
	errors = &ValidationErrors{
		Errors: []ValidationError{
			{Message: "Error 1"},
		},
	}
	if errors.Error() != "Error 1" {
		t.Errorf("Expected 'Error 1', got '%s'", errors.Error())
	}

	// Multiple errors
	errors = &ValidationErrors{
		Errors: []ValidationError{
			{Message: "Error 1"},
			{Message: "Error 2"},
			{Message: "Error 3"},
		},
	}
	expected := "Error 1; Error 2; Error 3"
	if errors.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, errors.Error())
	}
}

func TestValidationErrors_Add(t *testing.T) {
	errors := &ValidationErrors{}

	errors.Add([]any{"name"}, ErrCodeRequired, "Name is required")

	if len(errors.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors.Errors))
	}

	if errors.Errors[0].Code != ErrCodeRequired {
		t.Errorf("Expected error code %s, got %s", ErrCodeRequired, errors.Errors[0].Code)
	}

	if errors.Errors[0].Message != "Name is required" {
		t.Errorf("Expected message 'Name is required', got '%s'", errors.Errors[0].Message)
	}
}

func TestValidationErrors_AddWithMeta(t *testing.T) {
	errors := &ValidationErrors{}

	meta := map[string]any{
		"min": 5,
		"max": 10,
	}
	errors.AddWithMeta([]any{"age"}, ErrCodeTooSmall, "Age too small", meta)

	if len(errors.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors.Errors))
	}

	if errors.Errors[0].Meta == nil {
		t.Error("Expected meta to be set")
	}

	if errors.Errors[0].Meta["min"] != 5 {
		t.Errorf("Expected meta min to be 5, got %v", errors.Errors[0].Meta["min"])
	}
}

func TestValidationErrors_FormatErrors(t *testing.T) {
	errors := &ValidationErrors{}

	// Empty errors
	if errors.FormatErrors() != "" {
		t.Errorf("Expected empty string, got '%s'", errors.FormatErrors())
	}

	// With errors
	errors.Add([]any{"name"}, ErrCodeRequired, "Name is required")
	errors.Add([]any{"email"}, ErrCodeInvalidString, "Invalid email")

	formatted := errors.FormatErrors()
	if formatted == "" {
		t.Error("Expected formatted errors")
	}

	// Check that it contains error messages
	if len(formatted) < 20 {
		t.Error("Expected formatted string to contain error details")
	}
}

func TestValidationErrors_FormatErrorsJSON(t *testing.T) {
	errors := &ValidationErrors{}

	// Empty errors
	result := errors.FormatErrorsJSON()
	if result != nil {
		t.Errorf("Expected nil for empty errors, got %v", result)
	}

	// With errors
	errors.Add([]any{"name"}, ErrCodeRequired, "Name is required")
	errors.Add([]any{"email"}, ErrCodeInvalidString, "Invalid email")

	result = errors.FormatErrorsJSON()
	if result == nil {
		t.Error("Expected non-nil result")
	}

	if count, ok := result["count"].(int); !ok || count != 2 {
		t.Errorf("Expected count 2, got %v", result["count"])
	}

	errorsList, ok := result["errors"].([]map[string]any)
	if !ok {
		t.Error("Expected errors to be a slice of maps")
	}

	if len(errorsList) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(errorsList))
	}
}

func TestValidationErrors_GetErrorsByPath(t *testing.T) {
	errors := &ValidationErrors{}

	errors.Add([]any{"name"}, ErrCodeRequired, "Name required")
	errors.Add([]any{"email"}, ErrCodeInvalidString, "Invalid email")
	errors.Add([]any{"name"}, ErrCodeTooSmall, "Name too short")

	// Get errors for "name" path
	nameErrors := errors.GetErrorsByPath([]any{"name"})
	if len(nameErrors) != 2 {
		t.Errorf("Expected 2 errors for name path, got %d", len(nameErrors))
	}

	// Get errors for "email" path
	emailErrors := errors.GetErrorsByPath([]any{"email"})
	if len(emailErrors) != 1 {
		t.Errorf("Expected 1 error for email path, got %d", len(emailErrors))
	}

	// Get errors for non-existent path
	otherErrors := errors.GetErrorsByPath([]any{"other"})
	if len(otherErrors) != 0 {
		t.Errorf("Expected 0 errors for other path, got %d", len(otherErrors))
	}
}

func TestValidationErrors_GetErrorsByCode(t *testing.T) {
	errors := &ValidationErrors{}

	errors.Add([]any{"name"}, ErrCodeRequired, "Name required")
	errors.Add([]any{"email"}, ErrCodeInvalidString, "Invalid email")
	errors.Add([]any{"age"}, ErrCodeRequired, "Age required")

	// Get required errors
	requiredErrors := errors.GetErrorsByCode(ErrCodeRequired)
	if len(requiredErrors) != 2 {
		t.Errorf("Expected 2 required errors, got %d", len(requiredErrors))
	}

	// Get invalid string errors
	invalidStringErrors := errors.GetErrorsByCode(ErrCodeInvalidString)
	if len(invalidStringErrors) != 1 {
		t.Errorf("Expected 1 invalid string error, got %d", len(invalidStringErrors))
	}

	// Get non-existent code
	otherErrors := errors.GetErrorsByCode("nonexistent")
	if len(otherErrors) != 0 {
		t.Errorf("Expected 0 errors for nonexistent code, got %d", len(otherErrors))
	}
}

func TestValidationErrors_Flatten(t *testing.T) {
	errors := &ValidationErrors{}

	errors.Add([]any{}, ErrCodeRequired, "Form error")
	errors.Add([]any{"name"}, ErrCodeRequired, "Name required")
	errors.Add([]any{"email"}, ErrCodeInvalidString, "Invalid email")
	errors.Add([]any{"tags", 0}, ErrCodeRequired, "Tag required")
	errors.Add([]any{"tags", 1}, ErrCodeRequired, "Tag required")

	flattened := errors.Flatten()

	// Check form errors
	if len(flattened.FormErrors) != 1 {
		t.Errorf("Expected 1 form error, got %d", len(flattened.FormErrors))
	}

	// Check field errors
	if len(flattened.FieldErrors) != 3 {
		t.Errorf("Expected 3 field error groups, got %d", len(flattened.FieldErrors))
	}

	// Check that array errors are grouped
	if len(flattened.FieldErrors["tags"]) != 2 {
		t.Errorf("Expected 2 errors for tags, got %d", len(flattened.FieldErrors["tags"]))
	}
}

func TestPathToString(t *testing.T) {
	tests := []struct {
		path     []any
		expected string
	}{
		{[]any{}, ""},
		{[]any{"user"}, "user"},
		{[]any{"user", "email"}, "user.email"},
		{[]any{"test", 0}, "test[0]"},
		{[]any{"test", 1}, "test[1]"},
		{[]any{"user", "tags", 0}, "user.tags[0]"},
		{[]any{"user", "tags", 0, "name"}, "user.tags[0].name"},
		{[]any{0}, "[0]"},
		{[]any{"user", int64(5)}, "user[5]"},
	}

	for _, test := range tests {
		result := PathToString(test.path)
		if result != test.expected {
			t.Errorf("PathToString(%v) = '%s', expected '%s'", test.path, result, test.expected)
		}
	}
}

func TestPathEqual(t *testing.T) {
	tests := []struct {
		path1    []any
		path2    []any
		expected bool
	}{
		{[]any{}, []any{}, true},
		{[]any{"user"}, []any{"user"}, true},
		{[]any{"user", "email"}, []any{"user", "email"}, true},
		{[]any{"user"}, []any{"admin"}, false},
		{[]any{"user", "email"}, []any{"user"}, false},
		{[]any{"test", 0}, []any{"test", 0}, true},
		{[]any{"test", 0}, []any{"test", 1}, false},
	}

	for _, test := range tests {
		result := PathEqual(test.path1, test.path2)
		if result != test.expected {
			t.Errorf("PathEqual(%v, %v) = %v, expected %v", test.path1, test.path2, result, test.expected)
		}
	}
}

func TestPathAppend(t *testing.T) {
	path := []any{"user"}

	// Append string
	newPath := PathAppend(path, "email")
	if len(newPath) != 2 {
		t.Errorf("Expected path length 2, got %d", len(newPath))
	}
	if newPath[0] != "user" || newPath[1] != "email" {
		t.Errorf("Expected path ['user', 'email'], got %v", newPath)
	}

	// Original path should not be modified
	if len(path) != 1 {
		t.Error("Original path should not be modified")
	}

	// Append int
	newPath2 := PathAppend(newPath, 0)
	if len(newPath2) != 3 {
		t.Errorf("Expected path length 3, got %d", len(newPath2))
	}
	if newPath2[2] != 0 {
		t.Errorf("Expected last element to be 0, got %v", newPath2[2])
	}
}

func TestValidationErrors_PathCopy(t *testing.T) {
	errors := &ValidationErrors{}

	originalPath := []any{"user", "email"}
	errors.Add(originalPath, ErrCodeRequired, "Required")

	// Modify original path
	originalPath[0] = "modified"

	// Error path should not be affected
	if errors.Errors[0].Path[0] == "modified" {
		t.Error("Error path should not be affected by modifications to original")
	}
}

func TestValidationErrors_ComplexPaths(t *testing.T) {
	errors := &ValidationErrors{}

	// Add errors with various path structures
	errors.Add([]any{"users", 0, "name"}, ErrCodeRequired, "Name required")
	errors.Add([]any{"users", 0, "email"}, ErrCodeInvalidString, "Invalid email")
	errors.Add([]any{"users", 1, "name"}, ErrCodeRequired, "Name required")

	// Test flattening
	flattened := errors.Flatten()
	if len(flattened.FieldErrors) != 1 {
		t.Errorf("Expected 1 field group (users), got %d", len(flattened.FieldErrors))
	}

	// Array errors should be grouped under base field
	if len(flattened.FieldErrors["users"]) != 3 {
		t.Errorf("Expected 3 errors for users, got %d", len(flattened.FieldErrors["users"]))
	}
}
