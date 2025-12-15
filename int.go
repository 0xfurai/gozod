package gozod

import (
	"fmt"
)

// IntSchema validates integer values
type IntSchema struct {
	BaseSchema
	min         *int64
	max         *int64
	positive    bool
	negative    bool
	nonNegative bool
	nonPositive bool
	multipleOf  *int64
}

// Int creates a new integer schema
func Int() *IntSchema {
	return &IntSchema{
		BaseSchema: BaseSchema{required: true},
	}
}

// Nilable allows null values for IntSchema
func (s *IntSchema) Nilable() *IntSchema {
	s.nilable = true
	return s
}

// Min sets the minimum value for IntSchema
func (s *IntSchema) Min(value int64) *IntSchema {
	s.min = &value
	return s
}

// Max sets the maximum value for IntSchema
func (s *IntSchema) Max(value int64) *IntSchema {
	s.max = &value
	return s
}

// Positive validates that the number is positive (> 0) for IntSchema
func (s *IntSchema) Positive() *IntSchema {
	s.positive = true
	return s
}

// Negative validates that the number is negative (< 0) for IntSchema
func (s *IntSchema) Negative() *IntSchema {
	s.negative = true
	return s
}

// NonNegative validates that the number is non-negative (>= 0) for IntSchema
func (s *IntSchema) NonNegative() *IntSchema {
	s.nonNegative = true
	return s
}

// NonPositive validates that the number is non-positive (<= 0) for IntSchema
func (s *IntSchema) NonPositive() *IntSchema {
	s.nonPositive = true
	return s
}

// MultipleOf validates that the number is a multiple of the given value for IntSchema
func (s *IntSchema) MultipleOf(value int64) *IntSchema {
	s.multipleOf = &value
	return s
}

// Refine adds a custom validation function
// The function receives the value and returns (isValid, errorMessage)
// If isValid is false, validation fails with the provided errorMessage
// If errorMessage is empty, a default message will be used
func (s *IntSchema) Refine(validator RefineFunc) *IntSchema {
	s.BaseSchema.addRefinement(validator)
	return s
}

// SuperRefine adds a super refinement validation function
// Similar to Refine, but provides a context object for adding errors with custom paths and codes
// This allows for more fine-grained control over error reporting, similar to Zod's superRefine
func (s *IntSchema) SuperRefine(validator SuperRefineFunc) *IntSchema {
	s.BaseSchema.addSuperRefinement(validator)
	return s
}

// Validate validates a value against the int schema
func (s *IntSchema) Validate(value any, path []any) *ValidationErrors {
	errors := &ValidationErrors{}

	// Handle nil/nilable
	if value == nil {
		if !s.nilable {
			msg := s.getErrorMessage(path, ErrCodeRequired, "Required")
			errors.Add(path, ErrCodeRequired, msg)
			return errors
		}
		return nil
	}

	// Convert to int64 for validation
	var num int64
	var isInt bool

	switch v := value.(type) {
	case int:
		num = int64(v)
		isInt = true
	case int8:
		num = int64(v)
		isInt = true
	case int16:
		num = int64(v)
		isInt = true
	case int32:
		num = int64(v)
		isInt = true
	case int64:
		num = v
		isInt = true
	case uint:
		num = int64(v)
		isInt = true
	case uint8:
		num = int64(v)
		isInt = true
	case uint16:
		num = int64(v)
		isInt = true
	case uint32:
		num = int64(v)
		isInt = true
	case uint64:
		// Check if uint64 can fit in int64
		if v > uint64(9223372036854775807) {
			msg := s.getErrorMessage(path, ErrCodeInvalidType, fmt.Sprintf("Expected integer, got %T", value))
			errors.Add(path, ErrCodeInvalidType, msg)
			return errors
		}
		num = int64(v)
		isInt = true
	case float32:
		// Reject floats
		msg := s.getErrorMessage(path, ErrCodeInvalidType, "Expected integer, got float")
		errors.Add(path, ErrCodeInvalidType, msg)
		return errors
	case float64:
		// Reject floats
		msg := s.getErrorMessage(path, ErrCodeInvalidType, "Expected integer, got float")
		errors.Add(path, ErrCodeInvalidType, msg)
		return errors
	default:
		msg := s.getErrorMessage(path, ErrCodeInvalidType, fmt.Sprintf("Expected integer, got %T", value))
		errors.Add(path, ErrCodeInvalidType, msg)
		return errors
	}

	if !isInt {
		msg := s.getErrorMessage(path, ErrCodeInvalidType, "Expected integer, got float")
		errors.Add(path, ErrCodeInvalidType, msg)
		return errors
	}

	// Min validation
	if s.min != nil && num < *s.min {
		msg := s.getErrorMessage(path, ErrCodeTooSmall, fmt.Sprintf("Number must be greater than or equal to %v, got %v", *s.min, num))
		errors.Add(path, ErrCodeTooSmall, msg)
	}

	// Max validation
	if s.max != nil && num > *s.max {
		msg := s.getErrorMessage(path, ErrCodeTooBig, fmt.Sprintf("Number must be less than or equal to %v, got %v", *s.max, num))
		errors.Add(path, ErrCodeTooBig, msg)
	}

	// Positive validation
	if s.positive && num <= 0 {
		msg := s.getErrorMessage(path, ErrCodeTooSmall, fmt.Sprintf("Number must be positive (> 0), got %v", num))
		errors.Add(path, ErrCodeTooSmall, msg)
	}

	// Negative validation
	if s.negative && num >= 0 {
		msg := s.getErrorMessage(path, ErrCodeTooBig, fmt.Sprintf("Number must be negative (< 0), got %v", num))
		errors.Add(path, ErrCodeTooBig, msg)
	}

	// NonNegative validation
	if s.nonNegative && num < 0 {
		msg := s.getErrorMessage(path, ErrCodeTooSmall, fmt.Sprintf("Number must be non-negative (>= 0), got %v", num))
		errors.Add(path, ErrCodeTooSmall, msg)
	}

	// NonPositive validation
	if s.nonPositive && num > 0 {
		msg := s.getErrorMessage(path, ErrCodeTooBig, fmt.Sprintf("Number must be non-positive (<= 0), got %v", num))
		errors.Add(path, ErrCodeTooBig, msg)
	}

	// MultipleOf validation
	if s.multipleOf != nil {
		if num%*s.multipleOf != 0 {
			msg := s.getErrorMessage(path, ErrCodeInvalidType, fmt.Sprintf("Number must be a multiple of %v, got %v", *s.multipleOf, num))
			errors.Add(path, ErrCodeInvalidType, msg)
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

// CustomError sets a custom error message for a specific error code for IntSchema
func (s *IntSchema) CustomError(code, message string) *IntSchema {
	if s.BaseSchema.customErrors == nil {
		s.BaseSchema.customErrors = make(map[string]string)
	}
	s.BaseSchema.customErrors[code] = message
	return s
}

// SetErrorFormatter sets a custom error formatter function for IntSchema
func (s *IntSchema) SetErrorFormatter(formatter CustomErrorFunc) *IntSchema {
	s.BaseSchema.errorFormatter = formatter
	return s
}

// Type returns the schema type for IntSchema
func (s *IntSchema) Type() string {
	return "int"
}
