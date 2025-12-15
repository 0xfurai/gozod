package gozod

import (
	"testing"
)

func TestMapSchema_Required(t *testing.T) {
	schema := Map(map[string]Schema{
		"name": String(),
	})

	// Valid map
	err := schema.Validate(map[string]any{"name": "John"}, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid map, got: %v", err)
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

func TestMapSchema_Nilable(t *testing.T) {
	schema := Map(map[string]Schema{
		"name": String(),
	}).Nilable()

	// Valid map
	err := schema.Validate(map[string]any{"name": "John"}, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid map, got: %v", err)
	}

	// Nil value should pass
	err = schema.Validate(nil, nil)
	if err != nil {
		t.Error("Expected no errors for nil value with nilable schema")
	}
}

func TestMapSchema_FieldValidation(t *testing.T) {
	schema := Map(map[string]Schema{
		"name": String().Min(3),
		"age":  Int().Min(0),
	})

	// Valid: all fields pass
	err := schema.Validate(map[string]any{
		"name": "John",
		"age":  30,
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: name too short
	err = schema.Validate(map[string]any{
		"name": "Jo",
		"age":  30,
	}, nil)
	if err == nil {
		t.Error("Expected error for invalid name")
	}

	// Invalid: age negative
	err = schema.Validate(map[string]any{
		"name": "John",
		"age":  -5,
	}, nil)
	if err == nil {
		t.Error("Expected error for invalid age")
	}
}

func TestMapSchema_MissingFields(t *testing.T) {
	schema := Map(map[string]Schema{
		"name":  String(),
		"email": String().Nilable(),
	})

	// Valid: required field present, nilable field missing (treated as nil, which is allowed)
	err := schema.Validate(map[string]any{
		"name": "John",
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Valid: both fields present
	err = schema.Validate(map[string]any{
		"name":  "John",
		"email": "john@example.com",
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: required field missing
	err = schema.Validate(map[string]any{
		"email": "john@example.com",
	}, nil)
	if err == nil {
		t.Error("Expected error for missing required field")
	}
}

func TestMapSchema_NestedMaps(t *testing.T) {
	schema := Map(map[string]Schema{
		"user": Map(map[string]Schema{
			"name": String(),
			"age":  Int(),
		}),
	})

	// Valid: nested map
	err := schema.Validate(map[string]any{
		"user": map[string]any{
			"name": "John",
			"age":  30,
		},
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: nested field missing
	err = schema.Validate(map[string]any{
		"user": map[string]any{
			"name": "John",
		},
	}, nil)
	if err == nil {
		t.Error("Expected error for missing nested field")
	}
}

func TestMapSchema_Strict(t *testing.T) {
	schema := Map(map[string]Schema{
		"name": String(),
	}).Strict()

	// Valid: only known keys
	err := schema.Validate(map[string]any{
		"name": "John",
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: unknown key
	err = schema.Validate(map[string]any{
		"name":  "John",
		"extra": "value",
	}, nil)
	if err == nil {
		t.Error("Expected error for unknown key in strict mode")
		return
	}
	if err.Errors[0].Code != ErrCodeUnrecognizedKeys {
		t.Errorf("Expected error code %s, got %s", ErrCodeUnrecognizedKeys, err.Errors[0].Code)
	}
}

func TestMapSchema_AllowExtra(t *testing.T) {
	// Test default behavior: extra keys are allowed by default
	schema := Map(map[string]Schema{
		"name": String(),
	})

	// Valid: extra keys allowed by default
	err := schema.Validate(map[string]any{
		"name":  "John",
		"extra": "value",
		"more":  "data",
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}

func TestMapSchema_RejectsStructs(t *testing.T) {
	schema := Map(map[string]Schema{
		"name": String(),
		"age":  Int(),
	})

	// Test that structs are rejected
	type Person struct {
		Name string
		Age  int
	}

	person := Person{Name: "John", Age: 30}
	err := schema.Validate(person, nil)
	if err == nil {
		t.Error("Expected error for struct input, MapSchema should only accept maps")
		return
	}
	if err.Errors[0].Code != ErrCodeInvalidType {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
	}
}

func TestMapSchema_InvalidType(t *testing.T) {
	schema := Map(map[string]Schema{
		"name": String(),
	})

	// Invalid types (including structs)
	type Person struct {
		Name string
	}

	invalidValues := []any{
		"not a map",
		123,
		true,
		[]string{"test"},
		Person{Name: "John"}, // Structs should be rejected
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

func TestMapSchema_ComplexSchema(t *testing.T) {
	schema := Map(map[string]Schema{
		"name":  String().Min(3),
		"age":   Int().Min(18).Max(120),
		"email": String().Email().Nilable(),
		"tags":  Array(String()).Nilable(),
		"address": Map(map[string]Schema{
			"street": String(),
			"city":   String(),
		}).Nilable(),
	})

	// Valid: complex object
	err := schema.Validate(map[string]any{
		"name":  "John Doe",
		"age":   30,
		"email": "john@example.com",
		"tags":  []string{"developer", "golang"},
		"address": map[string]any{
			"street": "123 Main St",
			"city":   "New York",
		},
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: multiple errors
	err = schema.Validate(map[string]any{
		"name":  "Jo",            // too short
		"age":   150,             // too old
		"email": "invalid-email", // invalid email
	}, nil)
	if err == nil {
		t.Error("Expected multiple errors")
		return
	}
	if len(err.Errors) < 3 {
		t.Errorf("Expected at least 3 errors, got %d", len(err.Errors))
	}
}

func TestMapSchema_CustomError(t *testing.T) {
	schema := Map(map[string]Schema{
		"name": String(),
	}).CustomError(ErrCodeRequired, "Custom required error")

	err := schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error")
		return
	}
	if err.Errors[0].Message != "Custom required error" {
		t.Errorf("Expected custom error message, got: %s", err.Errors[0].Message)
	}
}

func TestMapSchema_SetErrorFormatter(t *testing.T) {
	schema := Map(map[string]Schema{
		"name": String(),
	}).SetErrorFormatter(func(path []any, code, defaultMessage string) string {
		return "Formatted: " + defaultMessage
	})

	err := schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error")
		return
	}
	if err.Errors[0].Message[:11] != "Formatted: " {
		t.Errorf("Expected formatted error message, got: %s", err.Errors[0].Message)
	}
}

func TestMapSchema_Type(t *testing.T) {
	schema := Map(map[string]Schema{
		"name": String(),
	})
	if schema.Type() != "object" {
		t.Errorf("Expected type 'object', got '%s'", schema.Type())
	}
}

func TestMapSchema_EmptyMap(t *testing.T) {
	schema := Map(map[string]Schema{
		"name": String().Nilable(),
	})

	// Empty map should be valid if all fields are nilable (missing fields treated as nil)
	err := schema.Validate(map[string]any{}, nil)
	if err != nil {
		t.Errorf("Expected empty map to be valid, got: %v", err)
	}

	// Empty map with required field should fail
	schema2 := Map(map[string]Schema{
		"name": String(),
	})
	err = schema2.Validate(map[string]any{}, nil)
	if err == nil {
		t.Error("Expected error for empty map with required field")
	}
}

func TestMapSchema_ErrorPaths(t *testing.T) {
	schema := Map(map[string]Schema{
		"user": Map(map[string]Schema{
			"name": String().Min(3),
		}),
	})

	err := schema.Validate(map[string]any{
		"user": map[string]any{
			"name": "Jo",
		},
	}, nil)

	if err == nil {
		t.Error("Expected error")
		return
	}

	// Check that error path includes nested field
	foundNestedError := false
	for _, validationErr := range err.Errors {
		if len(validationErr.Path) >= 2 {
			if validationErr.Path[0] == "user" && validationErr.Path[1] == "name" {
				foundNestedError = true
				break
			}
		}
	}
	if !foundNestedError {
		t.Error("Expected error path to include nested field path")
	}
}
