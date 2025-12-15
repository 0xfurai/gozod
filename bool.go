package gozod

import "fmt"

// BoolSchema validates boolean values
type BoolSchema struct {
	BaseSchema
}

// Bool creates a new boolean schema
func Bool() *BoolSchema {
	return &BoolSchema{
		BaseSchema: BaseSchema{required: true},
	}
}

// Nilable allows null values
func (s *BoolSchema) Nilable() *BoolSchema {
	s.nilable = true
	return s
}

// Refine adds a custom validation function
// The function receives the value and returns (isValid, errorMessage)
// If isValid is false, validation fails with the provided errorMessage
// If errorMessage is empty, a default message will be used
func (s *BoolSchema) Refine(validator RefineFunc) *BoolSchema {
	s.BaseSchema.addRefinement(validator)
	return s
}

// SuperRefine adds a super refinement validation function
// Similar to Refine, but provides a context object for adding errors with custom paths and codes
// This allows for more fine-grained control over error reporting, similar to Zod's superRefine
func (s *BoolSchema) SuperRefine(validator SuperRefineFunc) *BoolSchema {
	s.BaseSchema.addSuperRefinement(validator)
	return s
}

// Validate validates a value against the boolean schema
func (s *BoolSchema) Validate(value any, path []any) *ValidationErrors {
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

	// Type check
	_, ok := value.(bool)
	if !ok {
		msg := s.getErrorMessage(path, ErrCodeInvalidType, fmt.Sprintf("Expected boolean, got %T", value))
		errors.Add(path, ErrCodeInvalidType, msg)
		return errors
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
func (s *BoolSchema) CustomError(code, message string) *BoolSchema {
	if s.BaseSchema.customErrors == nil {
		s.BaseSchema.customErrors = make(map[string]string)
	}
	s.BaseSchema.customErrors[code] = message
	return s
}

// SetErrorFormatter sets a custom error formatter function
func (s *BoolSchema) SetErrorFormatter(formatter CustomErrorFunc) *BoolSchema {
	s.BaseSchema.errorFormatter = formatter
	return s
}

// Type returns the schema type
func (s *BoolSchema) Type() string {
	return "bool"
}
