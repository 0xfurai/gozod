package gozod

import (
	"fmt"
)

// FloatSchema validates float values
type FloatSchema struct {
	BaseSchema
	min         *float64
	max         *float64
	positive    bool
	negative    bool
	nonNegative bool
	nonPositive bool
	multipleOf  *float64
}

// Float creates a new float schema
func Float() *FloatSchema {
	return &FloatSchema{
		BaseSchema: BaseSchema{required: true},
	}
}

// Nilable allows null values for FloatSchema
func (s *FloatSchema) Nilable() *FloatSchema {
	s.nilable = true
	return s
}

// Min sets the minimum value for FloatSchema
func (s *FloatSchema) Min(value float64) *FloatSchema {
	s.min = &value
	return s
}

// Max sets the maximum value for FloatSchema
func (s *FloatSchema) Max(value float64) *FloatSchema {
	s.max = &value
	return s
}

// Positive validates that the number is positive (> 0) for FloatSchema
func (s *FloatSchema) Positive() *FloatSchema {
	s.positive = true
	return s
}

// Negative validates that the number is negative (< 0) for FloatSchema
func (s *FloatSchema) Negative() *FloatSchema {
	s.negative = true
	return s
}

// NonNegative validates that the number is non-negative (>= 0) for FloatSchema
func (s *FloatSchema) NonNegative() *FloatSchema {
	s.nonNegative = true
	return s
}

// NonPositive validates that the number is non-positive (<= 0) for FloatSchema
func (s *FloatSchema) NonPositive() *FloatSchema {
	s.nonPositive = true
	return s
}

// MultipleOf validates that the number is a multiple of the given value for FloatSchema
func (s *FloatSchema) MultipleOf(value float64) *FloatSchema {
	s.multipleOf = &value
	return s
}

// Refine adds a custom validation function
// The function receives the value and returns (isValid, errorMessage)
// If isValid is false, validation fails with the provided errorMessage
// If errorMessage is empty, a default message will be used
func (s *FloatSchema) Refine(validator RefineFunc) *FloatSchema {
	s.BaseSchema.addRefinement(validator)
	return s
}

// SuperRefine adds a super refinement validation function
// Similar to Refine, but provides a context object for adding errors with custom paths and codes
// This allows for more fine-grained control over error reporting, similar to Zod's superRefine
func (s *FloatSchema) SuperRefine(validator SuperRefineFunc) *FloatSchema {
	s.BaseSchema.addSuperRefinement(validator)
	return s
}

// Validate validates a value against the float schema
func (s *FloatSchema) Validate(value any, path []any) *ValidationErrors {
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

	// Convert to float64 for validation
	var num float64
	var isFloat bool

	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		// Reject integers
		msg := s.getErrorMessage(path, ErrCodeInvalidType, "Expected float, got integer")
		errors.Add(path, ErrCodeInvalidType, msg)
		return errors
	case float32:
		num = float64(v)
		isFloat = true
	case float64:
		num = v
		isFloat = true
	default:
		msg := s.getErrorMessage(path, ErrCodeInvalidType, fmt.Sprintf("Expected float, got %T", value))
		errors.Add(path, ErrCodeInvalidType, msg)
		return errors
	}

	if !isFloat {
		msg := s.getErrorMessage(path, ErrCodeInvalidType, "Expected float, got integer")
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
		remainder := num / *s.multipleOf
		// Check if remainder is close to an integer (handling floating point precision)
		if remainder-float64(int64(remainder)) > 0.0001 && remainder-float64(int64(remainder)) < 0.9999 {
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

// CustomError sets a custom error message for a specific error code for FloatSchema
func (s *FloatSchema) CustomError(code, message string) *FloatSchema {
	if s.BaseSchema.customErrors == nil {
		s.BaseSchema.customErrors = make(map[string]string)
	}
	s.BaseSchema.customErrors[code] = message
	return s
}

// SetErrorFormatter sets a custom error formatter function for FloatSchema
func (s *FloatSchema) SetErrorFormatter(formatter CustomErrorFunc) *FloatSchema {
	s.BaseSchema.errorFormatter = formatter
	return s
}

// Type returns the schema type for FloatSchema
func (s *FloatSchema) Type() string {
	return "float"
}
