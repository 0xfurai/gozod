package gozod

import (
	"testing"
)

func TestStructSchema_Required(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	schema := Struct(Shape{
		"name": String(),
	})

	// Valid struct
	user := User{Name: "John"}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid struct, got: %v", err)
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

func TestStructSchema_Nilable(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	schema := Struct(Shape{
		"name": String(),
	}).Nilable()

	// Valid struct
	user := User{Name: "John"}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors for valid struct, got: %v", err)
	}

	// Nil value should pass
	err = schema.Validate(nil, nil)
	if err != nil {
		t.Error("Expected no errors for nil value with nilable schema")
	}
}

func TestStructSchema_FieldValidation(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	schema := Struct(Shape{
		"name": String().Min(3),
		"age":  Int().Min(0),
	})

	// Valid: all fields pass
	user := User{Name: "John", Age: 30}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Invalid: name too short
	user = User{Name: "Jo", Age: 30}
	err = schema.Validate(user, nil)
	if err == nil {
		t.Error("Expected error for invalid name")
	}

	// Invalid: age negative
	user = User{Name: "John", Age: -5}
	err = schema.Validate(user, nil)
	if err == nil {
		t.Error("Expected error for invalid age")
	}
}

func TestStructSchema_MissingFields(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email,omitempty"`
	}

	schema := Struct(Shape{
		"name":  String(),
		"email": String().Nilable(),
	})

	// Valid: required field present, nilable field missing (zero value)
	user := User{Name: "John"}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Valid: both fields present
	user = User{Name: "John", Email: "john@example.com"}
	err = schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}

func TestStructSchema_RequiredFieldMissing(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	schema := Struct(Shape{
		"name":  String(),
		"email": String().Min(1), // Required and must be non-empty
	})

	// Empty string for required field with Min(1) constraint
	user := User{Name: "John"} // Email is empty string (zero value)
	err := schema.Validate(user, nil)
	if err == nil {
		t.Error("Expected error for empty required field")
	}
}

func TestStructSchema_JSONTags(t *testing.T) {
	type User struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Age       int    `json:"age"`
	}

	schema := Struct(Shape{
		"first_name": String().Min(2),
		"last_name":  String().Min(2),
		"age":        Int().Min(0),
	})

	// Valid struct with JSON tags
	user := User{
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
	}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}

func TestStructSchema_NoJSONTags(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	schema := Struct(Shape{
		"name": String().Min(2),
		"age":  Int().Min(0),
	})

	// Valid struct without JSON tags (uses camelCase field names)
	user := User{Name: "John", Age: 30}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}

func TestStructSchema_StrictMode(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Extra string `json:"extra"` // Not in schema
	}

	schema := Struct(Shape{
		"name":  String(),
		"email": String(),
	}).Strict()

	// Struct with extra field should fail in strict mode
	user := User{Name: "John", Email: "john@example.com", Extra: "extra"}
	err := schema.Validate(user, nil)
	if err == nil {
		t.Error("Expected error for extra field in strict mode")
		return
	}

	// Check that error is about unrecognized keys
	found := false
	for _, validationErr := range err.Errors {
		if validationErr.Code == ErrCodeUnrecognizedKeys {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected ErrCodeUnrecognizedKeys error")
	}
}

func TestStructSchema_StrictMode_NoExtraFields(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	schema := Struct(Shape{
		"name":  String(),
		"email": String(),
	}).Strict()

	// Struct with no extra fields should pass
	user := User{Name: "John", Email: "john@example.com"}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}

func TestStructSchema_PointerStruct(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	schema := Struct(Shape{
		"name": String(),
	})

	// Valid pointer to struct
	user := &User{Name: "John"}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Nil pointer should fail (not nilable)
	err = schema.Validate((*User)(nil), nil)
	if err == nil {
		t.Error("Expected error for nil pointer")
	}
}

func TestStructSchema_PointerStruct_Nilable(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	schema := Struct(Shape{
		"name": String(),
	}).Nilable()

	// Nil pointer should pass when nilable
	err := schema.Validate((*User)(nil), nil)
	if err != nil {
		t.Errorf("Expected no errors for nil pointer with nilable schema, got: %v", err)
	}
}

func TestStructSchema_PointerFields(t *testing.T) {
	type User struct {
		Name  string  `json:"name"`
		Email *string `json:"email,omitempty"`
	}

	email := "john@example.com"
	schema := Struct(Shape{
		"name":  String(),
		"email": String().Nilable(),
	})

	// Valid struct with pointer field set
	user := User{Name: "John", Email: &email}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Valid struct with nil pointer field
	user = User{Name: "John", Email: nil}
	err = schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors for nil pointer field, got: %v", err)
	}
}

func TestStructSchema_Omitempty(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email,omitempty"`
	}

	schema := Struct(Shape{
		"name":  String(),
		"email": String().Nilable(),
	})

	// Valid: email is empty string but has omitempty, so treated as nil
	user := User{Name: "John", Email: ""}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Valid: email provided
	user = User{Name: "John", Email: "john@example.com"}
	err = schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}

func TestStructSchema_NestedStructs(t *testing.T) {
	type User struct {
		Name    string `json:"name"`
		Address struct {
			Street string `json:"street"`
			City   string `json:"city"`
		} `json:"address"`
	}

	// Note: For nested structs, you'd typically use a nested Struct schema
	// But for this test, we'll just validate the top-level fields
	schema := Struct(Shape{
		"name": String().Min(2),
	})

	user := User{
		Name: "John",
		Address: struct {
			Street string `json:"street"`
			City   string `json:"city"`
		}{
			Street: "123 Main St",
			City:   "New York",
		},
	}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}

func TestStructSchema_InvalidType(t *testing.T) {
	schema := Struct(Shape{
		"name": String(),
	})

	// Invalid: not a struct
	err := schema.Validate("not a struct", nil)
	if err == nil {
		t.Error("Expected error for invalid type")
		return
	}

	if err.Errors[0].Code != ErrCodeInvalidType {
		t.Errorf("Expected error code %s, got %s", ErrCodeInvalidType, err.Errors[0].Code)
	}
}

func TestStructSchema_CustomError(t *testing.T) {
	schema := Struct(Shape{
		"name": String(),
	}).CustomError(ErrCodeRequired, "Name is required")

	err := schema.Validate(nil, nil)
	if err == nil {
		t.Error("Expected error")
		return
	}

	if err.Errors[0].Message != "Name is required" {
		t.Errorf("Expected custom error message, got: %s", err.Errors[0].Message)
	}
}

func TestStructSchema_ComplexExample(t *testing.T) {
	type User struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Email    string  `json:"email"`
		Age      *int    `json:"age,omitempty"`
		IsActive bool    `json:"isActive"`
		Bio      *string `json:"bio,omitempty"`
	}

	age := 30
	bio := "Software developer"
	schema := Struct(Shape{
		"id":       Int().Min(1),
		"name":     String().Min(2).Max(100),
		"email":    String().Email(),
		"age":      Int().Min(0).Max(150).Nilable(),
		"isActive": Bool(),
		"bio":      String().Max(500).Nilable(),
	})

	// Valid complex struct
	user := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      &age,
		IsActive: true,
		Bio:      &bio,
	}
	err := schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}

	// Valid with nilable fields as nil
	user = User{
		ID:       1,
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		IsActive: false,
	}
	err = schema.Validate(user, nil)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
}
