package gozod

import (
	"fmt"
	"reflect"
)

// MapSchema validates object/map values
type MapSchema struct {
	BaseSchema
	shape  map[string]Schema
	strict bool // If true, rejects unknown keys (default: false, allows extra keys)
}

// Map creates a new object/map schema
func Map(shape map[string]Schema) *MapSchema {
	return &MapSchema{
		BaseSchema: BaseSchema{required: true},
		shape:      shape,
		strict:     false, // Extra keys allowed by default
	}
}

// Nilable allows null values
func (s *MapSchema) Nilable() *MapSchema {
	s.nilable = true
	return s
}

// Strict rejects unknown keys
func (s *MapSchema) Strict() *MapSchema {
	s.strict = true
	return s
}

// Refine adds a custom validation function
// The function receives the value and returns (isValid, errorMessage)
// If isValid is false, validation fails with the provided errorMessage
// If errorMessage is empty, a default message will be used
func (s *MapSchema) Refine(validator RefineFunc) *MapSchema {
	s.BaseSchema.addRefinement(validator)
	return s
}

// SuperRefine adds a super refinement validation function
// Similar to Refine, but provides a context object for adding errors with custom paths and codes
// This allows for more fine-grained control over error reporting, similar to Zod's superRefine
func (s *MapSchema) SuperRefine(validator SuperRefineFunc) *MapSchema {
	s.BaseSchema.addSuperRefinement(validator)
	return s
}

// Validate validates a value against the object schema
func (s *MapSchema) Validate(value any, path []any) *ValidationErrors {
	errors := &ValidationErrors{}

	// Handle nil/nilable
	// Only nilable allows explicit nil values
	if value == nil {
		if !s.nilable {
			msg := s.getErrorMessage(path, ErrCodeRequired, "Required")
			errors.Add(path, ErrCodeRequired, msg)
			return errors
		}
		return nil
	}

	// Convert to map[string]any
	var obj map[string]any
	val := reflect.ValueOf(value)

	if val.Kind() != reflect.Map {
		msg := s.getErrorMessage(path, ErrCodeInvalidType, fmt.Sprintf("Expected map, got %T", value))
		errors.Add(path, ErrCodeInvalidType, msg)
		return errors
	}

	obj = make(map[string]any)
	for _, key := range val.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())
		obj[keyStr] = val.MapIndex(key).Interface()
	}

	// Validate each field in the shape
	for fieldName, schema := range s.shape {
		fieldPath := PathAppend(path, fieldName)

		fieldValue, exists := obj[fieldName]

		// If field is missing, pass nil to validation (will fail if required, pass if nilable)
		if !exists {
			fieldValue = nil
		}

		fieldErrors := schema.Validate(fieldValue, fieldPath)
		if fieldErrors != nil {
			errors.Errors = append(errors.Errors, fieldErrors.Errors...)
		}
	}

	// Check for unknown keys if strict mode is enabled
	if s.strict {
		for key := range obj {
			if _, exists := s.shape[key]; !exists {
				keyPath := PathAppend(path, key)
				msg := s.getErrorMessage(keyPath, ErrCodeUnrecognizedKeys, fmt.Sprintf("Unrecognized key '%s'", key))
				errors.Add(keyPath, ErrCodeUnrecognizedKeys, msg)
			}
		}
	}

	// Apply custom refinements (only if type check passed)
	s.applyRefinements(value, path, errors)

	// Apply super refinements (only if type check passed)
	s.applySuperRefinements(value, path, errors)

	if len(errors.Errors) == 0 {
		return nil
	}
	return errors
}

// CustomError sets a custom error message for a specific error code
func (s *MapSchema) CustomError(code, message string) *MapSchema {
	if s.BaseSchema.customErrors == nil {
		s.BaseSchema.customErrors = make(map[string]string)
	}
	s.BaseSchema.customErrors[code] = message
	return s
}

// SetErrorFormatter sets a custom error formatter function
func (s *MapSchema) SetErrorFormatter(formatter CustomErrorFunc) *MapSchema {
	s.BaseSchema.errorFormatter = formatter
	return s
}

// Type returns the schema type
func (s *MapSchema) Type() string {
	return "object"
}

// isEmptyValue checks if a value is empty (zero value)
func isEmptyValue(v any) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.String:
		return val.Len() == 0
	case reflect.Slice, reflect.Map:
		return val.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return val.IsNil()
	default:
		return val.IsZero()
	}
}
