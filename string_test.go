package gozod

import (
	"testing"
)

func TestStringSchema_Required(t *testing.T) {
	schema := String()

	// Valid string
	err := schema.Validate("test", nil)
	if err != nil {
		t.Errorf("Expected no errors for valid string, got: %v", err)
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

func TestStringSchema_Nilable(t *testing.T) {
	schema := String().Nilable()

	// Valid string
	err := schema.Validate("test", nil)
	if err != nil {
		t.Errorf("Expected no errors for valid string, got: %v", err)
	}

	// Nil value should pass (nilable allows nil)
	err = schema.Validate(nil, nil)
	if err != nil {
		t.Error("Expected no errors for nil value with nilable schema")
	}
}

func TestStringSchema_Min(t *testing.T) {
	schema := String().Min(5)

	// Valid: string meets minimum
	err := schema.Validate("hello", nil)
	if err != nil {
		t.Errorf("Expected no errors for string meeting minimum, got: %v", err)
	}

	// Invalid: string too short
	err = schema.Validate("hi", nil)
	if err == nil {
		t.Error("Expected error for string below minimum")
		return
	}
	if err.Errors[0].Code != ErrCodeTooSmall {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooSmall, err.Errors[0].Code)
	}
}

func TestStringSchema_Max(t *testing.T) {
	schema := String().Max(5)

	// Valid: string meets maximum
	err := schema.Validate("hello", nil)
	if err != nil {
		t.Errorf("Expected no errors for string meeting maximum, got: %v", err)
	}

	// Invalid: string too long
	err = schema.Validate("hello world", nil)
	if err == nil {
		t.Error("Expected error for string above maximum")
		return
	}
	if err.Errors[0].Code != ErrCodeTooBig {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooBig, err.Errors[0].Code)
	}
}

func TestStringSchema_MinMax(t *testing.T) {
	schema := String().Min(3).Max(10)

	// Valid: within range
	err := schema.Validate("hello", nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Too short
	err = schema.Validate("hi", nil)
	if err == nil {
		t.Error("Expected error for string too short")
	}

	// Too long
	err = schema.Validate("hello world this is too long", nil)
	if err == nil {
		t.Error("Expected error for string too long")
	}
}

func TestStringSchema_Email(t *testing.T) {
	schema := String().Email()

	// Valid emails
	validEmails := []string{
		"user@example.com",
		"test.email@domain.co.uk",
		"user+tag@example.com",
		"user_name@example-domain.com",
	}

	for _, email := range validEmails {
		err := schema.Validate(email, nil)
		if err != nil {
			t.Errorf("Expected valid email %s, got error: %v", email, err)
		}
	}

	// Invalid emails
	invalidEmails := []string{
		"notanemail",
		"@example.com",
		"user@",
		"user@.com",
		"user @example.com",
		"user@example",
	}

	for _, email := range invalidEmails {
		err := schema.Validate(email, nil)
		if err == nil {
			t.Errorf("Expected invalid email %s to fail, but it passed", email)
			continue
		}
		if err.Errors[0].Code != ErrCodeInvalidString {
			t.Errorf("Expected error code %s, got %s", ErrCodeInvalidString, err.Errors[0].Code)
		}
	}
}

func TestStringSchema_URL(t *testing.T) {
	schema := String().URL()

	// Valid URLs
	validURLs := []string{
		"http://example.com",
		"https://example.com",
		"http://example.com/path",
		"https://example.com/path?query=value",
	}

	for _, url := range validURLs {
		err := schema.Validate(url, nil)
		if err != nil {
			t.Errorf("Expected valid URL %s, got error: %v", url, err)
		}
	}

	// Invalid URLs
	invalidURLs := []string{
		"not a url",
		"example.com",
		"ftp://example.com",
		"http://",
		"http:// ",
	}

	for _, url := range invalidURLs {
		err := schema.Validate(url, nil)
		if err == nil {
			t.Errorf("Expected invalid URL %s to fail, but it passed", url)
		}
	}
}

func TestStringSchema_Regex(t *testing.T) {
	schema := String().Regex(`^[A-Z][a-z]+$`)

	// Valid: matches pattern
	err := schema.Validate("Hello", nil)
	if err != nil {
		t.Errorf("Expected no errors for matching pattern, got: %v", err)
	}

	// Invalid: doesn't match pattern
	err = schema.Validate("hello", nil)
	if err == nil {
		t.Error("Expected error for non-matching pattern")
	}

	// Test with custom message
	schema2 := String().Regex(`^\d+$`, "Must be digits only")
	err = schema2.Validate("abc", nil)
	if err == nil {
		t.Error("Expected error for non-matching pattern")
		return
	}
	if err.Errors[0].Message != "Must be digits only" {
		t.Errorf("Expected custom message, got: %s", err.Errors[0].Message)
	}
}

func TestStringSchema_OneOf(t *testing.T) {
	schema := String().OneOf("red", "green", "blue")

	// Valid: one of the options
	err := schema.Validate("red", nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	err = schema.Validate("green", nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: not one of the options
	err = schema.Validate("yellow", nil)
	if err == nil {
		t.Error("Expected error for value not in oneOf")
		return
	}
	if err.Errors[0].Code != ErrCodeInvalidEnumValue {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidEnumValue, err.Errors[0].Code)
	}
}

func TestStringSchema_NotOneOf(t *testing.T) {
	schema := String().NotOneOf("admin", "root", "system")

	// Valid: not in the forbidden list
	err := schema.Validate("user", nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: in the forbidden list
	err = schema.Validate("admin", nil)
	if err == nil {
		t.Error("Expected error for value in notOneOf")
	}
}

func TestStringSchema_StartsWith(t *testing.T) {
	schema := String().StartsWith("https://")

	// Valid: starts with prefix
	err := schema.Validate("https://example.com", nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: doesn't start with prefix
	err = schema.Validate("http://example.com", nil)
	if err == nil {
		t.Error("Expected error for string not starting with prefix")
	}
}

func TestStringSchema_EndsWith(t *testing.T) {
	schema := String().EndsWith(".com")

	// Valid: ends with suffix
	err := schema.Validate("example.com", nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: doesn't end with suffix
	err = schema.Validate("example.org", nil)
	if err == nil {
		t.Error("Expected error for string not ending with suffix")
	}
}

func TestStringSchema_Includes(t *testing.T) {
	schema := String().Includes("test")

	// Valid: includes substring
	err := schema.Validate("this is a test string", nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: doesn't include substring
	err = schema.Validate("this is a string", nil)
	if err == nil {
		t.Error("Expected error for string not including substring")
	}
}

func TestStringSchema_InvalidType(t *testing.T) {
	schema := String()

	// Invalid types
	invalidValues := []any{
		123,
		45.67,
		true,
		[]string{"test"},
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

func TestStringSchema_ChainedValidators(t *testing.T) {
	schema := String().
		Min(5).
		Max(20).
		Email().
		StartsWith("user@")

	// Valid: meets all criteria
	err := schema.Validate("user@example.com", nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: too short
	err = schema.Validate("u@e.c", nil)
	if err == nil {
		t.Error("Expected error for string too short")
	}

	// Invalid: doesn't start with prefix
	err = schema.Validate("admin@example.com", nil)
	if err == nil {
		t.Error("Expected error for string not starting with prefix")
	}
}

func TestStringSchema_CustomError(t *testing.T) {
	schema := String().
		Min(5).
		CustomError(ErrCodeTooSmall, "Custom min length error")

	err := schema.Validate("hi", nil)
	if err == nil {
		t.Error("Expected error")
		return
	}
	if err.Errors[0].Message != "Custom min length error" {
		t.Errorf("Expected custom error message, got: %s", err.Errors[0].Message)
	}
}

func TestStringSchema_SetErrorFormatter(t *testing.T) {
	schema := String().
		Min(5).
		SetErrorFormatter(func(path []any, code, defaultMessage string) string {
			return "Formatted: " + defaultMessage
		})

	err := schema.Validate("hi", nil)
	if err == nil {
		t.Error("Expected error")
		return
	}
	if err.Errors[0].Message[:11] != "Formatted: " {
		t.Errorf("Expected formatted error message, got: %s", err.Errors[0].Message)
	}
}

func TestStringSchema_Type(t *testing.T) {
	schema := String()
	if schema.Type() != "string" {
		t.Errorf("Expected type 'string', got '%s'", schema.Type())
	}
}

func TestStringSchema_EmptyString(t *testing.T) {
	schema := String()

	// Empty string should be valid by default
	err := schema.Validate("", nil)
	if err != nil {
		t.Errorf("Expected empty string to be valid, got: %v", err)
	}

	// With Min(1), empty string should fail
	schema2 := String().Min(1)
	err = schema2.Validate("", nil)
	if err == nil {
		t.Error("Expected error for empty string with Min(1)")
	}
}

func TestStringSchema_Refine(t *testing.T) {
	// Test refine with custom validation
	schema := String().Refine(func(value any) (bool, string) {
		str := value.(string)
		if len(str) > 0 && str[0] == 'A' {
			return true, ""
		}
		return false, "String must start with 'A'"
	})

	// Valid: string starts with 'A'
	err := schema.Validate("Apple", nil)
	if err != nil {
		t.Errorf("Expected no errors for string starting with 'A', got: %v", err)
	}

	// Invalid: string doesn't start with 'A'
	err = schema.Validate("Banana", nil)
	if err == nil {
		t.Error("Expected error for string not starting with 'A'")
		return
	}
	if err.Errors[0].Code != ErrCodeCustomValidation {
		t.Errorf("Expected error code %s, got %s", ErrCodeCustomValidation, err.Errors[0].Code)
	}
	if err.Errors[0].Message != "String must start with 'A'" {
		t.Errorf("Expected error message 'String must start with 'A'', got '%s'", err.Errors[0].Message)
	}

	// Test refine with default error message
	schema2 := String().Refine(func(value any) (bool, string) {
		str := value.(string)
		return str != "forbidden", ""
	})

	err = schema2.Validate("forbidden", nil)
	if err == nil {
		t.Error("Expected error for forbidden string")
		return
	}
	if err.Errors[0].Code != ErrCodeCustomValidation {
		t.Errorf("Expected error code %s, got %s", ErrCodeCustomValidation, err.Errors[0].Code)
	}

	// Test multiple refinements
	schema3 := String().Refine(func(value any) (bool, string) {
		str := value.(string)
		return len(str) >= 3, "String must be at least 3 characters"
	}).Refine(func(value any) (bool, string) {
		str := value.(string)
		return str != "bad", "String cannot be 'bad'"
	})

	// Valid
	err = schema3.Validate("good", nil)
	if err != nil {
		t.Errorf("Expected no errors for valid string, got: %v", err)
	}

	// Invalid: too short
	err = schema3.Validate("ab", nil)
	if err == nil {
		t.Error("Expected error for string too short")
		return
	}
	if len(err.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(err.Errors))
	}

	// Invalid: forbidden value
	err = schema3.Validate("bad", nil)
	if err == nil {
		t.Error("Expected error for forbidden string")
		return
	}
	if len(err.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(err.Errors))
	}

	// Test refine doesn't run if type check fails
	schema4 := String().Refine(func(value any) (bool, string) {
		return false, "This should not run"
	})

	err = schema4.Validate(123, nil)
	if err == nil {
		t.Error("Expected error for wrong type")
		return
	}
	if err.Errors[0].Code != ErrCodeInvalidType {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
	}
	// Should not have the refine error
	for _, e := range err.Errors {
		if e.Message == "This should not run" {
			t.Error("Refine should not run when type check fails")
		}
	}
}
