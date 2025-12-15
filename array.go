package gozod

import (
	"fmt"
	"reflect"
)

// ArraySchema validates array/slice values
type ArraySchema struct {
	BaseSchema
	elementSchema Schema
	minLength     *int
	maxLength     *int
	nonEmpty      bool
}

// Array creates a new array schema
func Array(elementSchema Schema) *ArraySchema {
	return &ArraySchema{
		BaseSchema:    BaseSchema{required: true},
		elementSchema: elementSchema,
	}
}

// Nilable allows null values
func (s *ArraySchema) Nilable() *ArraySchema {
	s.nilable = true
	return s
}

// Min sets the minimum length
func (s *ArraySchema) Min(length int) *ArraySchema {
	s.minLength = &length
	return s
}

// Max sets the maximum length
func (s *ArraySchema) Max(length int) *ArraySchema {
	s.maxLength = &length
	return s
}

// NonEmpty validates that the array is not empty
func (s *ArraySchema) NonEmpty() *ArraySchema {
	s.nonEmpty = true
	return s
}

// Refine adds a custom validation function
// The function receives the value and returns (isValid, errorMessage)
// If isValid is false, validation fails with the provided errorMessage
// If errorMessage is empty, a default message will be used
func (s *ArraySchema) Refine(validator RefineFunc) *ArraySchema {
	s.BaseSchema.addRefinement(validator)
	return s
}

// SuperRefine adds a super refinement validation function
// Similar to Refine, but provides a context object for adding errors with custom paths and codes
// This allows for more fine-grained control over error reporting, similar to Zod's superRefine
func (s *ArraySchema) SuperRefine(validator SuperRefineFunc) *ArraySchema {
	s.BaseSchema.addSuperRefinement(validator)
	return s
}

// Validate validates a value against the array schema
func (s *ArraySchema) Validate(value any, path []any) *ValidationErrors {
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

	// Convert to slice
	var slice []any
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		slice = make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			slice[i] = val.Index(i).Interface()
		}
	default:
		msg := s.getErrorMessage(path, ErrCodeInvalidType, fmt.Sprintf("Expected array, got %T", value))
		errors.Add(path, ErrCodeInvalidType, msg)
		return errors
	}

	// Length validations
	if s.nonEmpty && len(slice) == 0 {
		msg := s.getErrorMessage(path, ErrCodeTooSmall, "Array must not be empty")
		errors.Add(path, ErrCodeTooSmall, msg)
	}

	if s.minLength != nil && len(slice) < *s.minLength {
		msg := s.getErrorMessage(path, ErrCodeTooSmall, fmt.Sprintf("Array must have at least %d element(s), got %d", *s.minLength, len(slice)))
		errors.Add(path, ErrCodeTooSmall, msg)
	}

	if s.maxLength != nil && len(slice) > *s.maxLength {
		msg := s.getErrorMessage(path, ErrCodeTooBig, fmt.Sprintf("Array must have at most %d element(s), got %d", *s.maxLength, len(slice)))
		errors.Add(path, ErrCodeTooBig, msg)
	}

	// Validate each element
	for i, element := range slice {
		elementPath := PathAppend(path, i)
		elementErrors := s.elementSchema.Validate(element, elementPath)
		if elementErrors != nil {
			errors.Errors = append(errors.Errors, elementErrors.Errors...)
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
func (s *ArraySchema) CustomError(code, message string) *ArraySchema {
	if s.BaseSchema.customErrors == nil {
		s.BaseSchema.customErrors = make(map[string]string)
	}
	s.BaseSchema.customErrors[code] = message
	return s
}

// SetErrorFormatter sets a custom error formatter function
func (s *ArraySchema) SetErrorFormatter(formatter CustomErrorFunc) *ArraySchema {
	s.BaseSchema.errorFormatter = formatter
	return s
}

// Type returns the schema type
func (s *ArraySchema) Type() string {
	return "array"
}
