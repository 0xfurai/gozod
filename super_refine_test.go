package gozod

import (
	"reflect"
	"testing"
)

func TestSuperRefine_StringSchema(t *testing.T) {
	// Test superRefine with custom error path
	schema := String().SuperRefine(func(value any, ctx *SuperRefineContext) {
		str := value.(string)
		if len(str) < 5 {
			ctx.AddIssue([]any{}, ErrCodeTooSmall, "String must be at least 5 characters")
		}
		if str == "forbidden" {
			ctx.AddIssue([]any{}, ErrCodeInvalidString, "This value is forbidden")
		}
	})

	// Valid string
	err := schema.Validate("hello", nil)
	if err != nil {
		t.Errorf("Expected no errors for valid string, got: %v", err)
	}

	// Invalid - too short
	err = schema.Validate("hi", nil)
	if err == nil {
		t.Error("Expected error for string too short")
		return
	}
	if len(err.Errors) == 0 {
		t.Error("Expected at least one error")
		return
	}
	if err.Errors[0].Code != ErrCodeTooSmall {
		t.Errorf("Expected error code %s, got %s", ErrCodeTooSmall, err.Errors[0].Code)
	}

	// Invalid - forbidden value
	err = schema.Validate("forbidden", nil)
	if err == nil {
		t.Error("Expected error for forbidden value")
		return
	}
	found := false
	for _, e := range err.Errors {
		if e.Code == ErrCodeInvalidString && e.Message == "This value is forbidden" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find error for forbidden value")
	}
}

func TestSuperRefine_MapSchema(t *testing.T) {
	// Test superRefine with cross-field validation (like password confirmation)
	schema := Map(map[string]Schema{
		"password": String().Min(8),
		"confirm":  String().Min(8),
	}).SuperRefine(func(value any, ctx *SuperRefineContext) {
		m := value.(map[string]any)
		password, ok1 := m["password"].(string)
		confirm, ok2 := m["confirm"].(string)
		if !ok1 || !ok2 {
			return // Let type validation handle this
		}
		if password != confirm {
			// Add error to the "confirm" field path
			ctx.AddIssue([]any{"confirm"}, ErrCodeCustomValidation, "Passwords do not match")
		}
	})

	// Valid: passwords match
	err := schema.Validate(map[string]any{
		"password": "password123",
		"confirm":  "password123",
	}, nil)
	if err != nil {
		t.Errorf("Expected no errors for matching passwords, got: %v", err)
	}

	// Invalid: passwords don't match
	err = schema.Validate(map[string]any{
		"password": "password123",
		"confirm":  "different123",
	}, nil)
	if err == nil {
		t.Error("Expected error for non-matching passwords")
		return
	}
	// Should have the superRefine error on the "confirm" field
	found := false
	for _, e := range err.Errors {
		if len(e.Path) > 0 && e.Path[len(e.Path)-1] == "confirm" {
			if e.Code == ErrCodeCustomValidation && e.Message == "Passwords do not match" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected to find superRefine error on confirm field for non-matching passwords")
	}
}

func TestSuperRefine_MapSchema_MultipleErrors(t *testing.T) {
	// Test superRefine adding multiple errors
	schema := Map(map[string]Schema{
		"start": Int(),
		"end":   Int(),
	}).SuperRefine(func(value any, ctx *SuperRefineContext) {
		m := value.(map[string]any)
		start, ok1 := m["start"].(int64)
		end, ok2 := m["end"].(int64)
		if !ok1 || !ok2 {
			return
		}
		if start >= end {
			ctx.AddIssue([]any{"start"}, ErrCodeTooBig, "Start must be less than end")
		}
		if end-start > 100 {
			ctx.AddIssue([]any{"end"}, ErrCodeTooBig, "Range must not exceed 100")
		}
	})

	// Invalid: start >= end
	err := schema.Validate(map[string]any{
		"start": int64(50),
		"end":   int64(30),
	}, nil)
	if err == nil {
		t.Error("Expected error for start >= end")
		return
	}
	found := false
	for _, e := range err.Errors {
		if len(e.Path) > 0 && e.Path[len(e.Path)-1] == "start" {
			if e.Message == "Start must be less than end" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected to find error on start field")
	}

	// Invalid: range too large
	err = schema.Validate(map[string]any{
		"start": int64(0),
		"end":   int64(150),
	}, nil)
	if err == nil {
		t.Error("Expected error for range too large")
		return
	}
	found = false
	for _, e := range err.Errors {
		if len(e.Path) > 0 && e.Path[len(e.Path)-1] == "end" {
			if e.Message == "Range must not exceed 100" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected to find error on end field for range too large")
	}
}

func TestSuperRefine_ArraySchema(t *testing.T) {
	// Test superRefine with array validation
	schema := Array(Int()).SuperRefine(func(value any, ctx *SuperRefineContext) {
		// Check for unique values
		val := reflect.ValueOf(value)
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			return
		}

		seen := make(map[int64]bool)
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i).Interface()
			var num int64
			switch n := elem.(type) {
			case int:
				num = int64(n)
			case int64:
				num = n
			case int32:
				num = int64(n)
			default:
				continue
			}
			if seen[num] {
				// Add error to the specific array index
				ctx.AddIssue([]any{i}, ErrCodeCustomValidation, "Duplicate value found")
				return
			}
			seen[num] = true
		}
	})

	// Valid: unique values
	err := schema.Validate([]int{1, 2, 3}, nil)
	if err != nil {
		t.Errorf("Expected no errors for unique array, got: %v", err)
	}

	// Invalid: duplicate values
	err = schema.Validate([]int{1, 2, 2}, nil)
	if err == nil {
		t.Error("Expected error for array with duplicates")
		return
	}
	found := false
	for _, e := range err.Errors {
		// Check if error is on index 2
		if len(e.Path) > 0 {
			if idx, ok := e.Path[len(e.Path)-1].(int); ok && idx == 2 {
				if e.Code == ErrCodeCustomValidation && e.Message == "Duplicate value found" {
					found = true
					break
				}
			}
		}
	}
	if !found {
		t.Error("Expected to find superRefine error on duplicate array element")
	}
}

func TestSuperRefine_IntSchema(t *testing.T) {
	// Test superRefine with custom error code
	schema := Int().SuperRefine(func(value any, ctx *SuperRefineContext) {
		num := value.(int64)
		if num%2 != 0 {
			ctx.AddIssue([]any{}, "odd_number", "Number must be even")
		}
	})

	// Valid: even number
	err := schema.Validate(int64(4), nil)
	if err != nil {
		t.Errorf("Expected no errors for even number, got: %v", err)
	}

	// Invalid: odd number
	err = schema.Validate(int64(5), nil)
	if err == nil {
		t.Error("Expected error for odd number")
		return
	}
	found := false
	for _, e := range err.Errors {
		if e.Code == "odd_number" && e.Message == "Number must be even" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find error with custom code for odd number")
	}
}

func TestSuperRefine_WithMeta(t *testing.T) {
	// Test superRefine with metadata
	schema := String().SuperRefine(func(value any, ctx *SuperRefineContext) {
		str := value.(string)
		if len(str) < 3 {
			meta := map[string]any{
				"minLength":    3,
				"actualLength": len(str),
			}
			ctx.AddIssueWithMeta([]any{}, ErrCodeTooSmall, "String too short", meta)
		}
	})

	err := schema.Validate("ab", nil)
	if err == nil {
		t.Error("Expected error for string too short")
		return
	}
	found := false
	for _, e := range err.Errors {
		if e.Code == ErrCodeTooSmall && e.Meta != nil {
			if e.Meta["minLength"] == 3 && e.Meta["actualLength"] == 2 {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected to find error with metadata")
	}
}

func TestSuperRefine_NestedPath(t *testing.T) {
	// Test superRefine with nested paths
	schema := Map(map[string]Schema{
		"user": Map(map[string]Schema{
			"profile": Map(map[string]Schema{
				"age": Int(),
			}),
		}),
	}).SuperRefine(func(value any, ctx *SuperRefineContext) {
		m := value.(map[string]any)
		if user, ok := m["user"].(map[string]any); ok {
			if profile, ok := user["profile"].(map[string]any); ok {
				if age, ok := profile["age"].(int64); ok {
					if age < 18 {
						// Add error to nested path
						ctx.AddIssue([]any{"user", "profile", "age"}, ErrCodeTooSmall, "Age must be at least 18")
					}
				}
			}
		}
	})

	err := schema.Validate(map[string]any{
		"user": map[string]any{
			"profile": map[string]any{
				"age": int64(15),
			},
		},
	}, nil)
	if err == nil {
		t.Error("Expected error for age < 18")
		return
	}
	found := false
	for _, e := range err.Errors {
		if len(e.Path) == 3 && e.Path[0] == "user" && e.Path[1] == "profile" && e.Path[2] == "age" {
			if e.Message == "Age must be at least 18" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Expected to find error on nested path")
	}
}
