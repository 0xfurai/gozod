package gozod

import (
	"testing"
)

func TestInt_Required(t *testing.T) {
	schema := Int()

	// Valid int
	err := schema.Validate(42, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid int, got: %v", err)
	}

	// Nil value should fail
	err = schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error for nil value")
	}
}

func TestInt_OnlyAcceptsIntegers(t *testing.T) {
	schema := Int()

	// Valid: various integer types
	validInts := []any{
		int(42),
		int8(42),
		int16(42),
		int32(42),
		int64(42),
		uint(42),
		uint8(42),
		uint16(42),
		uint32(42),
		uint64(42),
	}

	for _, val := range validInts {
		err := schema.Validate(val, nil)
		if err != nil {
			t.Errorf("Expected valid int %v (%T), got error: %v", val, val, err)
		}
	}

	// Invalid: floats
	invalidFloats := []any{
		float32(42.5),
		float64(42.5),
	}

	for _, val := range invalidFloats {
		err := schema.Validate(val, nil)
		if err == nil {
			t.Errorf("Expected error for float %v (%T), but it passed", val, val)
			continue
		}
		if err.Errors[0].Code != ErrCodeInvalidType {
			t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
		}
	}
}

func TestIntSchema_Min(t *testing.T) {
	schema := Int().Min(10)

	// Valid: meets minimum
	err := schema.Validate(15, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	err = schema.Validate(10, nil)
	if err != nil {
		t.Errorf("Expected no errors for boundary value, got: %v", err)
	}

	// Invalid: below minimum
	err = schema.Validate(5, nil)
	if err == nil {
		t.Error("Expected error for value below minimum")
		return
	}
	if err.Errors[0].Code != ErrCodeTooSmall {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooSmall, err.Errors[0].Code)
	}
}

func TestIntSchema_Max(t *testing.T) {
	schema := Int().Max(10)

	// Valid: meets maximum
	err := schema.Validate(5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	err = schema.Validate(10, nil)
	if err != nil {
		t.Errorf("Expected no errors for boundary value, got: %v", err)
	}

	// Invalid: above maximum
	err = schema.Validate(15, nil)
	if err == nil {
		t.Error("Expected error for value above maximum")
		return
	}
	if err.Errors[0].Code != ErrCodeTooBig {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooBig, err.Errors[0].Code)
	}
}

func TestIntSchema_Positive(t *testing.T) {
	schema := Int().Positive()

	// Valid: positive number
	err := schema.Validate(5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: zero
	err = schema.Validate(0, nil)
	if err == nil {
		t.Error("Expected error for zero with Positive()")
	}

	// Invalid: negative
	err = schema.Validate(-5, nil)
	if err == nil {
		t.Error("Expected error for negative number with Positive()")
	}
}

func TestIntSchema_Negative(t *testing.T) {
	schema := Int().Negative()

	// Valid: negative number
	err := schema.Validate(-5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: zero
	err = schema.Validate(0, nil)
	if err == nil {
		t.Error("Expected error for zero with Negative()")
	}

	// Invalid: positive
	err = schema.Validate(5, nil)
	if err == nil {
		t.Error("Expected error for positive number with Negative()")
	}
}

func TestIntSchema_NonNegative(t *testing.T) {
	schema := Int().NonNegative()

	// Valid: zero
	err := schema.Validate(0, nil)
	if err != nil {
		t.Errorf("Expected no errors for zero, got: %v", err)
	}

	// Valid: positive
	err = schema.Validate(5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: negative
	err = schema.Validate(-5, nil)
	if err == nil {
		t.Error("Expected error for negative number with NonNegative()")
	}
}

func TestIntSchema_NonPositive(t *testing.T) {
	schema := Int().NonPositive()

	// Valid: zero
	err := schema.Validate(0, nil)
	if err != nil {
		t.Errorf("Expected no errors for zero, got: %v", err)
	}

	// Valid: negative
	err = schema.Validate(-5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: positive
	err = schema.Validate(5, nil)
	if err == nil {
		t.Error("Expected error for positive number with NonPositive()")
	}
}

func TestIntSchema_MultipleOf(t *testing.T) {
	schema := Int().MultipleOf(5)

	// Valid: multiple of 5
	err := schema.Validate(10, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	err = schema.Validate(15, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	err = schema.Validate(0, nil)
	if err != nil {
		t.Errorf("Expected no errors for zero, got: %v", err)
	}

	// Invalid: not a multiple
	err = schema.Validate(7, nil)
	if err == nil {
		t.Error("Expected error for number not multiple of 5")
	}

	err = schema.Validate(13, nil)
	if err == nil {
		t.Error("Expected error for number not multiple of 5")
	}
}

func TestIntSchema_ChainedValidators(t *testing.T) {
	schema := Int().
		Min(18).
		Max(120).
		NonNegative()

	// Valid: meets all criteria
	err := schema.Validate(25, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: below minimum
	err = schema.Validate(15, nil)
	if err == nil {
		t.Error("Expected error for value below minimum")
	}

	// Invalid: above maximum
	err = schema.Validate(150, nil)
	if err == nil {
		t.Error("Expected error for value above maximum")
	}
}

func TestIntSchema_InvalidType(t *testing.T) {
	schema := Int()

	// Invalid types
	invalidValues := []any{
		"123",
		true,
		[]int{1, 2, 3},
		map[string]int{"key": 1},
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

func TestIntSchema_CustomError(t *testing.T) {
	schema := Int().
		Min(10).
		CustomError(ErrCodeTooSmall, "Custom min error")

	err := schema.Validate(5, nil)
	if err == nil {
		t.Error("Expected error")
		return
	}
	if err.Errors[0].Message != "Custom min error" {
		t.Errorf("Expected custom error message, got: %s", err.Errors[0].Message)
	}
}

func TestIntSchema_SetErrorFormatter(t *testing.T) {
	schema := Int().
		Min(10).
		SetErrorFormatter(func(path []any, code, defaultMessage string) string {
			return "Formatted: " + defaultMessage
		})

	err := schema.Validate(5, nil)
	if err == nil {
		t.Error("Expected error")
		return
	}
	if err.Errors[0].Message[:11] != "Formatted: " {
		t.Errorf("Expected formatted error message, got: %s", err.Errors[0].Message)
	}
}

func TestIntSchema_Type(t *testing.T) {
	intSchema := Int()
	if intSchema.Type() != "int" {
		t.Errorf("Expected type 'int', got '%s'", intSchema.Type())
	}
}

func TestIntSchema_ZeroValue(t *testing.T) {
	schema := Int()

	// Zero should be valid by default
	err := schema.Validate(0, nil)
	if err != nil {
		t.Errorf("Expected zero to be valid, got: %v", err)
	}

	// With Positive(), zero should fail
	schema2 := Int().Positive()
	err = schema2.Validate(0, nil)
	if err == nil {
		t.Error("Expected error for zero with Positive()")
	}
}
