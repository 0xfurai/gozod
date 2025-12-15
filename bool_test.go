package gozod

import (
	"testing"
)

func TestBoolSchema_Required(t *testing.T) {
	schema := Bool()

	// Valid bool
	err := schema.Validate(true, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid bool, got: %v", err)
	}

	err = schema.Validate(false, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid bool, got: %v", err)
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

func TestBoolSchema_Nilable(t *testing.T) {
	schema := Bool().Nilable()

	// Valid bool
	err := schema.Validate(true, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid bool, got: %v", err)
	}

	// Nil value should pass
	err = schema.Validate(nil, nil)
	if err != nil {
		t.Error("Expected no errors for nil value with nilable schema")
	}
}

func TestBoolSchema_InvalidType(t *testing.T) {
	schema := Bool()

	// Invalid types
	invalidValues := []any{
		"true",
		"false",
		1,
		0,
		"yes",
		"no",
		[]bool{true},
		map[string]bool{"key": true},
	}

	for _, val := range invalidValues {
		err := schema.Validate(val, nil)
		if err == nil {
			t.Errorf("Expected error for invalid type %T (%v), but it passed", val, val)
			continue
		}
		if err.Errors[0].Code != ErrCodeInvalidType {
			t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
		}
	}
}

func TestBoolSchema_CustomError(t *testing.T) {
	schema := Bool().
		CustomError(ErrCodeRequired, "Custom required error")

	err := schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error")
		return
	}
	if err.Errors[0].Message != "Custom required error" {
		t.Errorf("Expected custom error message, got: %s", err.Errors[0].Message)
	}
}

func TestBoolSchema_SetErrorFormatter(t *testing.T) {
	schema := Bool().
		SetErrorFormatter(func(path []any, code, defaultMessage string) string {
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

func TestBoolSchema_Type(t *testing.T) {
	schema := Bool()
	if schema.Type() != "bool" {
		t.Errorf("Expected type 'bool', got '%s'", schema.Type())
	}
}

func TestBoolSchema_TrueAndFalse(t *testing.T) {
	schema := Bool()

	// Both true and false should be valid
	err := schema.Validate(true, nil)
	if err != nil {
		t.Errorf("Expected true to be valid, got: %v", err)
	}

	err = schema.Validate(false, nil)
	if err != nil {
		t.Errorf("Expected false to be valid, got: %v", err)
	}
}
