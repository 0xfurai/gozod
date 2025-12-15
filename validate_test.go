package gozod

import (
	"reflect"
	"testing"
)

func TestValidate_StringSchema(t *testing.T) {
	schema := String()

	// Valid string
	err := schema.Validate("test", nil)
	if err != nil {
		t.Errorf("Expected no errors for valid string, got: %v", err)
	}

	// Invalid string - nil value
	err = schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error for nil value")
		return
	}
	if err.Errors[0].Code != ErrCodeRequired {
		t.Errorf("Expected error code %s, got %s", ErrCodeRequired, err.Errors[0].Code)
	}

	// Invalid string - wrong type
	err = schema.Validate(123, nil)
	if err == nil {
		t.Error("Expected error for wrong type")
		return
	}
	if err.Errors[0].Code != ErrCodeInvalidType {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
	}

	// Invalid string - validation failure
	schema = String().Min(5)
	err = schema.Validate("abc", nil)
	if err == nil {
		t.Error("Expected error for string too short")
		return
	}
	if err.Errors[0].Code != ErrCodeTooSmall {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooSmall, err.Errors[0].Code)
	}
}

func TestValidate_IntSchema(t *testing.T) {
	schema := Int()

	// Valid integer
	err := schema.Validate(42, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid integer, got: %v", err)
	}

	// Invalid - nil value
	err = schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error for nil value")
	}

	// Invalid - wrong type
	err = schema.Validate("not a number", nil)
	if err == nil {
		t.Error("Expected error for wrong type")
		return
	}
	if err.Errors[0].Code != ErrCodeInvalidType {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
	}

	// Invalid - validation failure
	schema = Int().Min(10)
	err = schema.Validate(5, nil)
	if err == nil {
		t.Error("Expected error for number too small")
		return
	}
	if err.Errors[0].Code != ErrCodeTooSmall {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooSmall, err.Errors[0].Code)
	}
}

func TestValidate_FloatSchema(t *testing.T) {
	schema := Float()

	// Valid float
	err := schema.Validate(3.14, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid float, got: %v", err)
	}

	// Invalid - wrong type (integer when float expected)
	err = schema.Validate(42, nil)
	if err == nil {
		t.Error("Expected error for integer when float expected")
		return
	}
	if err.Errors[0].Code != ErrCodeInvalidType {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
	}
}

func TestValidate_BoolSchema(t *testing.T) {
	schema := Bool()

	// Valid boolean
	err := schema.Validate(true, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid boolean, got: %v", err)
	}

	err = schema.Validate(false, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid boolean, got: %v", err)
	}

	// Invalid - nil value
	err = schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error for nil value")
	}

	// Invalid - wrong type
	err = schema.Validate("not a bool", nil)
	if err == nil {
		t.Error("Expected error for wrong type")
		return
	}
	if err.Errors[0].Code != ErrCodeInvalidType {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
	}
}

func TestValidate_ArraySchema(t *testing.T) {
	schema := Array(String())

	// Valid array
	err := schema.Validate([]string{"a", "b", "c"}, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid array, got: %v", err)
	}

	// Invalid - nil value
	err = schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error for nil value")
	}

	// Invalid - wrong type
	err = schema.Validate("not an array", nil)
	if err == nil {
		t.Error("Expected error for wrong type")
		return
	}
	if err.Errors[0].Code != ErrCodeInvalidType {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
	}

	// Invalid - validation failure
	schema = Array(String()).Min(3)
	err = schema.Validate([]string{"a", "b"}, nil)
	if err == nil {
		t.Error("Expected error for array too short")
		return
	}
	if err.Errors[0].Code != ErrCodeTooSmall {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooSmall, err.Errors[0].Code)
	}
}

func TestValidate_MapSchema(t *testing.T) {
	schema := Map(map[string]Schema{
		"name": String(),
		"age":  Int(),
	})

	// Valid map
	err := schema.Validate(map[string]any{
		"name": "John",
		"age":  30,
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid map, got: %v", err)
	}

	// Invalid - nil value
	err = schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error for nil value")
	}

	// Invalid - wrong type
	err = schema.Validate("not a map", nil)
	if err == nil {
		t.Error("Expected error for wrong type")
		return
	}
	if err.Errors[0].Code != ErrCodeInvalidType {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
	}

	// Invalid - validation failure (missing required field)
	err = schema.Validate(map[string]any{
		"name": "John",
	}, nil)
	if err == nil {
		t.Error("Expected error for missing required field")
		return
	}
	if err.Errors[0].Code != ErrCodeRequired {
		t.Errorf("Expected error code %s, got %s", ErrCodeRequired, err.Errors[0].Code)
	}
}

func TestValidate_NilableSchema(t *testing.T) {
	schema := String().Nilable()

	// Valid string
	err := schema.Validate("test", nil)
	if err != nil {
		t.Errorf("Expected no errors for valid string, got: %v", err)
	}

	// Nilable should allow nil even when required is true
	err = schema.Validate(nil, nil)
	if err != nil {
		t.Error("Expected no errors for nil value with nilable schema (nilable allows nil even when required)")
	}
}

func TestValidate_ChainedValidators(t *testing.T) {
	schema := String().Min(3).Max(20).Email()

	// Valid email
	err := schema.Validate("test@example.com", nil)
	if err != nil {
		t.Errorf("Expected no errors for valid email, got: %v", err)
	}

	// Invalid - too short
	err = schema.Validate("ab", nil)
	if err == nil {
		t.Error("Expected error for string too short")
	}

	// Invalid - not an email
	err = schema.Validate("notanemail", nil)
	if err == nil {
		t.Error("Expected error for invalid email")
	}
}

func TestValidate_ComplexNestedSchema(t *testing.T) {
	schema := Map(map[string]Schema{
		"users": Array(Map(map[string]Schema{
			"name":  String(),
			"email": String().Email(),
			"age":   Int().Min(18),
		})),
		"count": Int(),
	})

	// Valid nested structure
	err := schema.Validate(map[string]any{
		"users": []map[string]any{
			{
				"name":  "John",
				"email": "john@example.com",
				"age":   25,
			},
			{
				"name":  "Jane",
				"email": "jane@example.com",
				"age":   30,
			},
		},
		"count": 2,
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid nested structure, got: %v", err)
	}

	// Invalid - missing required field
	err = schema.Validate(map[string]any{
		"users": []map[string]any{
			{
				"name": "John",
				"age":  25,
			},
		},
		"count": 1,
	}, nil)
	if err == nil {
		t.Error("Expected error for missing required field")
	}

	// Invalid - invalid email
	err = schema.Validate(map[string]any{
		"users": []map[string]any{
			{
				"name":  "John",
				"email": "notanemail",
				"age":   25,
			},
		},
		"count": 1,
	}, nil)
	if err == nil {
		t.Error("Expected error for invalid email")
	}

	// Invalid - age too small
	err = schema.Validate(map[string]any{
		"users": []map[string]any{
			{
				"name":  "John",
				"email": "john@example.com",
				"age":   15,
			},
		},
		"count": 1,
	}, nil)
	if err == nil {
		t.Error("Expected error for age too small")
	}
}

func TestValidate_ErrorPath(t *testing.T) {
	schema := Map(map[string]Schema{
		"user": Map(map[string]Schema{
			"profile": Map(map[string]Schema{
				"name": String(),
			}),
		}),
	})

	// Missing nested field
	err := schema.Validate(map[string]any{
		"user": map[string]any{
			"profile": map[string]any{},
		},
	}, nil)
	if err == nil {
		t.Error("Expected error for missing nested field")
		return
	}

	// Verify error path
	if len(err.Errors) > 0 {
		expectedPath := []any{"user", "profile", "name"}
		if !PathEqual(err.Errors[0].Path, expectedPath) {
			t.Errorf("Expected path %v, got %v", expectedPath, err.Errors[0].Path)
		}
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	// Empty string with required - empty string is a valid string value
	// It will only fail if there's a Min constraint
	schema := String().Min(1)
	err := schema.Validate("", nil)
	if err == nil {
		t.Error("Expected error for empty string with min length constraint")
	}

	// Empty array with required - empty array is a valid array value
	// It will only fail if there's a Min or NonEmpty constraint
	arraySchema := Array(String()).NonEmpty()
	err = arraySchema.Validate([]string{}, nil)
	if err == nil {
		t.Error("Expected error for empty array with NonEmpty constraint")
	}

	// Empty map with required - empty map is valid, but missing required fields will fail
	mapSchema := Map(map[string]Schema{
		"field": String(),
	})
	err = mapSchema.Validate(map[string]any{}, nil)
	if err == nil {
		t.Error("Expected error for empty map missing required field")
	}
}

func TestValidate_ZeroValue(t *testing.T) {
	// Zero integer with positive requirement
	var schema Schema
	schema = Int().Positive()
	err := schema.Validate(0, nil)
	if err == nil {
		t.Error("Expected error for zero value with positive requirement")
	}

	// Zero float with positive requirement
	schema = Float().Positive()
	err = schema.Validate(0.0, nil)
	if err == nil {
		t.Error("Expected error for zero float with positive requirement")
	}
}

func TestValidate_ReturnsValidationErrors(t *testing.T) {
	schema := String()
	err := schema.Validate(nil, nil)

	// Verify it returns ValidationErrors type (non-nil when there are errors)
	if err == nil {
		t.Error("Expected ValidationErrors, got nil")
		return
	}

	// Verify error structure
	if len(err.Errors) == 0 {
		t.Error("Expected at least one error")
		return
	}

	error := err.Errors[0]
	if error.Code == "" {
		t.Error("Expected error code to be set")
	}
	if error.Message == "" {
		t.Error("Expected error message to be set")
	}
}

func TestValidate_Refine(t *testing.T) {
	// Test refine with Int schema
	intSchema := Int().Refine(func(value any) (bool, string) {
		// Handle different integer types
		var num int64
		switch v := value.(type) {
		case int:
			num = int64(v)
		case int64:
			num = v
		case int32:
			num = int64(v)
		case int16:
			num = int64(v)
		case int8:
			num = int64(v)
		default:
			return false, "Invalid integer type"
		}
		return num%2 == 0, "Number must be even"
	})

	err := intSchema.Validate(4, nil)
	if err != nil {
		t.Errorf("Expected no errors for even number, got: %v", err)
	}

	err = intSchema.Validate(5, nil)
	if err == nil {
		t.Error("Expected error for odd number")
		return
	}
	if err.Errors[0].Code != ErrCodeCustomValidation {
		t.Errorf("Expected error code %s, got %s", ErrCodeCustomValidation, err.Errors[0].Code)
	}

	// Test refine with Map schema
	mapSchema := Map(map[string]Schema{
		"password": String().Min(8),
		"confirm":  String().Min(8),
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

	// Valid: passwords match
	err = mapSchema.Validate(map[string]any{
		"password": "password123",
		"confirm":  "password123",
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors for matching passwords, got: %v", err)
	}

	// Invalid: passwords don't match
	err = mapSchema.Validate(map[string]any{
		"password": "password123",
		"confirm":  "different123",
	}, nil)
	if err == nil {
		t.Error("Expected error for non-matching passwords")
		return
	}
	// Should have the refine error
	found := false
	for _, e := range err.Errors {
		if e.Code == ErrCodeCustomValidation && e.Message == "Passwords do not match" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find refine error for non-matching passwords")
	}

	// Test refine with Array schema
	arraySchema := Array(Int()).Refine(func(value any) (bool, string) {
		// Handle different slice types using reflect
		val := reflect.ValueOf(value)
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			return false, "Value is not a slice or array"
		}

		if val.Len() == 0 {
			return true, ""
		}

		// Check if all elements are unique
		seen := make(map[int64]bool)
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i).Interface()
			// Handle different integer types from array elements
			var num int64
			switch n := elem.(type) {
			case int:
				num = int64(n)
			case int64:
				num = n
			case int32:
				num = int64(n)
			case int16:
				num = int64(n)
			case int8:
				num = int64(n)
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

	// Valid: unique values
	err = arraySchema.Validate([]int{1, 2, 3}, nil)
	if err != nil {
		t.Errorf("Expected no errors for unique array, got: %v", err)
	}

	// Invalid: duplicate values
	err = arraySchema.Validate([]int{1, 2, 2}, nil)
	if err == nil {
		t.Error("Expected error for array with duplicates")
		return
	}
	found = false
	for _, e := range err.Errors {
		if e.Code == ErrCodeCustomValidation && e.Message == "Array must contain unique values" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find refine error for duplicate values")
	}
}
