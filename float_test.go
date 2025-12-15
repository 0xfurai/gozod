package gozod

import (
	"testing"
)

func TestFloat_Required(t *testing.T) {
	schema := Float()

	// Valid float
	err := schema.Validate(42.5, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid float, got: %v", err)
	}

	// Nil value should fail
	err = schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error for nil value")
	}
}

func TestFloat_OnlyAcceptsFloats(t *testing.T) {
	schema := Float()

	// Valid: floats
	validFloats := []any{
		float32(42.5),
		float64(42.5),
	}

	for _, val := range validFloats {
		err := schema.Validate(val, nil)
		if err != nil {
			t.Errorf("Expected valid float %v (%T), got error: %v", val, val, err)
		}
	}

	// Invalid: integers
	invalidInts := []any{
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

	for _, val := range invalidInts {
		err := schema.Validate(val, nil)
		if err == nil {
			t.Errorf("Expected error for integer %v (%T), but it passed", val, val)
			continue
		}
		if err.Errors[0].Code != ErrCodeInvalidType {
			t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
		}
	}
}

func TestFloatSchema_Min(t *testing.T) {
	schema := Float().Min(10.5)

	// Valid: meets minimum
	err := schema.Validate(15.5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	err = schema.Validate(10.5, nil)
	if err != nil {
		t.Errorf("Expected no errors for boundary value, got: %v", err)
	}

	// Invalid: below minimum
	err = schema.Validate(5.5, nil)
	if err == nil {
		t.Error("Expected error for value below minimum")
		return
	}
	if err.Errors[0].Code != ErrCodeTooSmall {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooSmall, err.Errors[0].Code)
	}
}

func TestFloatSchema_Max(t *testing.T) {
	schema := Float().Max(10.5)

	// Valid: meets maximum
	err := schema.Validate(5.5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	err = schema.Validate(10.5, nil)
	if err != nil {
		t.Errorf("Expected no errors for boundary value, got: %v", err)
	}

	// Invalid: above maximum
	err = schema.Validate(15.5, nil)
	if err == nil {
		t.Error("Expected error for value above maximum")
		return
	}
	if err.Errors[0].Code != ErrCodeTooBig {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooBig, err.Errors[0].Code)
	}
}

func TestFloatSchema_Positive(t *testing.T) {
	schema := Float().Positive()

	// Valid: positive number
	err := schema.Validate(5.5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: zero
	err = schema.Validate(0.0, nil)
	if err == nil {
		t.Error("Expected error for zero with Positive()")
	}

	// Invalid: negative
	err = schema.Validate(-5.5, nil)
	if err == nil {
		t.Error("Expected error for negative number with Positive()")
	}
}

func TestFloatSchema_Negative(t *testing.T) {
	schema := Float().Negative()

	// Valid: negative number
	err := schema.Validate(-5.5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: zero
	err = schema.Validate(0.0, nil)
	if err == nil {
		t.Error("Expected error for zero with Negative()")
	}

	// Invalid: positive
	err = schema.Validate(5.5, nil)
	if err == nil {
		t.Error("Expected error for positive number with Negative()")
	}
}

func TestFloatSchema_NonNegative(t *testing.T) {
	schema := Float().NonNegative()

	// Valid: zero
	err := schema.Validate(0.0, nil)
	if err != nil {
		t.Errorf("Expected no errors for zero, got: %v", err)
	}

	// Valid: positive
	err = schema.Validate(5.5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: negative
	err = schema.Validate(-5.5, nil)
	if err == nil {
		t.Error("Expected error for negative number with NonNegative()")
	}
}

func TestFloatSchema_NonPositive(t *testing.T) {
	schema := Float().NonPositive()

	// Valid: zero
	err := schema.Validate(0.0, nil)
	if err != nil {
		t.Errorf("Expected no errors for zero, got: %v", err)
	}

	// Valid: negative
	err = schema.Validate(-5.5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: positive
	err = schema.Validate(5.5, nil)
	if err == nil {
		t.Error("Expected error for positive number with NonPositive()")
	}
}

func TestFloatSchema_MultipleOf(t *testing.T) {
	schema := Float().MultipleOf(2.5)

	// Valid: multiple of 2.5
	err := schema.Validate(5.0, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	err = schema.Validate(7.5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	err = schema.Validate(0.0, nil)
	if err != nil {
		t.Errorf("Expected no errors for zero, got: %v", err)
	}

	// Invalid: not a multiple
	err = schema.Validate(6.2, nil)
	if err == nil {
		t.Error("Expected error for number not multiple of 2.5")
	}

	err = schema.Validate(8.1, nil)
	if err == nil {
		t.Error("Expected error for number not multiple of 2.5")
	}
}

func TestFloatSchema_ChainedValidators(t *testing.T) {
	schema := Float().
		Min(18.5).
		Max(120.5).
		NonNegative()

	// Valid: meets all criteria
	err := schema.Validate(25.5, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: below minimum
	err = schema.Validate(15.5, nil)
	if err == nil {
		t.Error("Expected error for value below minimum")
	}

	// Invalid: above maximum
	err = schema.Validate(150.5, nil)
	if err == nil {
		t.Error("Expected error for value above maximum")
	}
}

func TestFloatSchema_InvalidType(t *testing.T) {
	schema := Float()

	// Invalid types
	invalidValues := []any{
		"123.5",
		true,
		[]float64{1.5, 2.5, 3.5},
		map[string]float64{"key": 1.5},
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

func TestFloatSchema_CustomError(t *testing.T) {
	schema := Float().
		Min(10.5).
		CustomError(ErrCodeTooSmall, "Custom min error")

	err := schema.Validate(5.5, nil)
	if err == nil {
		t.Error("Expected error")
		return
	}
	if err.Errors[0].Message != "Custom min error" {
		t.Errorf("Expected custom error message, got: %s", err.Errors[0].Message)
	}
}

func TestFloatSchema_SetErrorFormatter(t *testing.T) {
	schema := Float().
		Min(10.5).
		SetErrorFormatter(func(path []any, code, defaultMessage string) string {
			return "Formatted: " + defaultMessage
		})

	err := schema.Validate(5.5, nil)
	if err == nil {
		t.Error("Expected error")
		return
	}
	if err.Errors[0].Message[:11] != "Formatted: " {
		t.Errorf("Expected formatted error message, got: %s", err.Errors[0].Message)
	}
}

func TestFloatSchema_Type(t *testing.T) {
	floatSchema := Float()
	if floatSchema.Type() != "float" {
		t.Errorf("Expected type 'float', got '%s'", floatSchema.Type())
	}
}

func TestFloatSchema_ZeroValue(t *testing.T) {
	schema := Float()

	// Zero should be valid by default
	err := schema.Validate(0.0, nil)
	if err != nil {
		t.Errorf("Expected zero to be valid, got: %v", err)
	}

	// With Positive(), zero should fail
	schema2 := Float().Positive()
	err = schema2.Validate(0.0, nil)
	if err == nil {
		t.Error("Expected error for zero with Positive()")
	}
}
