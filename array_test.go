package gozod

import (
	"testing"
)

func TestArraySchema_Required(t *testing.T) {
	schema := Array(String())

	// Valid array
	err := schema.Validate([]string{"hello", "world"}, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid array, got: %v", err)
	}

	// Nil value should fail
	err = schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error for nil value")
		return
	}

	if err.Errors[0].Code != ErrCodeRequired {
		t.Errorf("Expected error code %s, got %s", ErrCodeRequired, err.Errors[0].Code)
	}
}

func TestArraySchema_Nilable(t *testing.T) {
	schema := Array(String()).Nilable()

	// Valid array
	err := schema.Validate([]string{"hello"}, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid array, got: %v", err)
	}

	// Nil value should pass
	err = schema.Validate(nil, nil)
	if err != nil {
		t.Error("Expected no errors for nil value with nilable schema")
	}
}

func TestArraySchema_Min(t *testing.T) {
	schema := Array(String()).Min(3)

	// Valid: meets minimum
	err := schema.Validate([]string{"a", "b", "c"}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: below minimum
	err = schema.Validate([]string{"a", "b"}, nil)
	if err == nil {
		t.Error("Expected error for array below minimum")
		return
	}
	if err.Errors[0].Code != ErrCodeTooSmall {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooSmall, err.Errors[0].Code)
	}
}

func TestArraySchema_Max(t *testing.T) {
	schema := Array(String()).Max(3)

	// Valid: meets maximum
	err := schema.Validate([]string{"a", "b", "c"}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: above maximum
	err = schema.Validate([]string{"a", "b", "c", "d"}, nil)
	if err == nil {
		t.Error("Expected error for array above maximum")
		return
	}
	if err.Errors[0].Code != ErrCodeTooBig {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooBig, err.Errors[0].Code)
	}
}

func TestArraySchema_NonEmpty(t *testing.T) {
	schema := Array(String()).NonEmpty()

	// Valid: non-empty array
	err := schema.Validate([]string{"hello"}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: empty array
	err = schema.Validate([]string{}, nil)
	if err == nil {
		t.Error("Expected error for empty array")
		return
	}
	if err.Errors[0].Code != ErrCodeTooSmall {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooSmall, err.Errors[0].Code)
	}
}

func TestArraySchema_ElementValidation(t *testing.T) {
	schema := Array(String().Min(3))

	// Valid: all elements pass
	err := schema.Validate([]string{"hello", "world"}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: some elements fail
	err = schema.Validate([]string{"hello", "hi"}, nil)
	if err == nil {
		t.Error("Expected error for invalid element")
		return
	}

	// Check that error path includes index
	foundIndexError := false
	for _, validationErr := range err.Errors {
		if len(validationErr.Path) > 0 {
			if idx, ok := validationErr.Path[len(validationErr.Path)-1].(int); ok && idx == 1 {
				foundIndexError = true
				break
			}
		}
	}
	if !foundIndexError {
		t.Error("Expected error path to include array index")
	}
}

func TestArraySchema_NestedArrays(t *testing.T) {
	schema := Array(Array(Int()))

	// Valid: nested arrays
	err := schema.Validate([][]int{{1, 2}, {3, 4}}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: wrong element type in nested array
	err = schema.Validate([]any{[]int{1, 2}, []any{"invalid"}}, nil)
	if err == nil {
		t.Error("Expected error for invalid nested element")
	}
}

func TestArraySchema_MixedTypes(t *testing.T) {
	// Array of floats (accepts both ints and floats when converted)
	schema := Array(Float())

	// Valid: floats
	err := schema.Validate([]any{1.0, 2.5, 3.0, 4.7}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: wrong type
	err = schema.Validate([]any{1.0, "invalid", 3.0}, nil)
	if err == nil {
		t.Error("Expected error for invalid element type")
	}
}

func TestArraySchema_InvalidType(t *testing.T) {
	schema := Array(String())

	// Invalid types
	invalidValues := []any{
		"not an array",
		123,
		true,
		map[string]string{"key": "value"},
	}

	for _, val := range invalidValues {
		err := schema.Validate(val, nil)
		if err == nil {
			t.Errorf("Expected error for invalid type %T, but it passed", val)
			continue
		}
		if err.Errors[0].Code != ErrCodeInvalidType {
			t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
		}
	}
}

func TestArraySchema_EmptyArray(t *testing.T) {
	schema := Array(String())

	// Empty array should be valid by default
	err := schema.Validate([]string{}, nil)
	if err != nil {
		t.Errorf("Expected empty array to be valid, got: %v", err)
	}

	// With NonEmpty(), empty array should fail
	schema2 := Array(String()).NonEmpty()
	err = schema2.Validate([]string{}, nil)
	if err == nil {
		t.Error("Expected error for empty array with NonEmpty()")
	}
}

func TestArraySchema_ChainedValidators(t *testing.T) {
	schema := Array(String().Min(3)).
		Min(2).
		Max(5).
		NonEmpty()

	// Valid: meets all criteria
	err := schema.Validate([]string{"hello", "world"}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: too short
	err = schema.Validate([]string{"hi"}, nil)
	if err == nil {
		t.Error("Expected error for array below minimum")
	}

	// Invalid: too long
	err = schema.Validate([]string{"hello", "world", "test", "foo", "bar", "baz"}, nil)
	if err == nil {
		t.Error("Expected error for array above maximum")
	}
}

func TestArraySchema_CustomError(t *testing.T) {
	schema := Array(String()).
		Min(2).
		CustomError(ErrCodeTooSmall, "Custom min length error")

	err := schema.Validate([]string{"one"}, nil)
	if err == nil {
		t.Error("Expected error")
		return
	}
	if err.Errors[0].Message != "Custom min length error" {
		t.Errorf("Expected custom error message, got: %s", err.Errors[0].Message)
	}
}

func TestArraySchema_SetErrorFormatter(t *testing.T) {
	schema := Array(String()).
		Min(2).
		SetErrorFormatter(func(path []any, code, defaultMessage string) string {
			return "Formatted: " + defaultMessage
		})

	err := schema.Validate([]string{"one"}, nil)
	if err == nil {
		t.Error("Expected error")
		return
	}
	if err.Errors[0].Message[:11] != "Formatted: " {
		t.Errorf("Expected formatted error message, got: %s", err.Errors[0].Message)
	}
}

func TestArraySchema_Type(t *testing.T) {
	schema := Array(String())
	if schema.Type() != "array" {
		t.Errorf("Expected type 'array', got '%s'", schema.Type())
	}
}

func TestArraySchema_ArrayTypes(t *testing.T) {
	schema := Array(Int())

	// Test with different slice types
	err := schema.Validate([]int{1, 2, 3}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	err = schema.Validate([]int8{1, 2, 3}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	err = schema.Validate([]int64{1, 2, 3}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}

func TestArraySchema_ComplexElementSchema(t *testing.T) {
	// Array of objects
	schema := Array(Map(map[string]Schema{
		"name": String().Min(1),
		"age":  Int().Min(0),
	}))

	// Valid
	err := schema.Validate([]map[string]any{
		{"name": "Alice", "age": 30},
		{"name": "Bob", "age": 25},
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: missing required field
	err = schema.Validate([]map[string]any{
		{"name": "Alice"},
	}, nil)
	if err == nil {
		t.Error("Expected error for missing required field")
	}
}
