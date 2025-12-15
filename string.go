package gozod

import (
	"fmt"
	"regexp"
	"strings"
)

// StringSchema validates string values
type StringSchema struct {
	BaseSchema
	minLength    *int
	maxLength    *int
	email        bool
	url          bool
	regex        *regexp.Regexp
	regexMessage string
	oneOf        []string
	notOneOf     []string
	startsWith   *string
	endsWith     *string
	includes     *string
}

// String creates a new string schema
func String() *StringSchema {
	return &StringSchema{
		BaseSchema: BaseSchema{required: true},
	}
}

// Nilable allows null values
func (s *StringSchema) Nilable() *StringSchema {
	s.nilable = true
	return s
}

// Min sets the minimum length
func (s *StringSchema) Min(length int) *StringSchema {
	s.minLength = &length
	return s
}

// Max sets the maximum length
func (s *StringSchema) Max(length int) *StringSchema {
	s.maxLength = &length
	return s
}

// Email validates email format
func (s *StringSchema) Email() *StringSchema {
	s.email = true
	return s
}

// URL validates URL format
func (s *StringSchema) URL() *StringSchema {
	s.url = true
	return s
}

// Regex validates against a regular expression
func (s *StringSchema) Regex(pattern string, message ...string) *StringSchema {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("invalid regex pattern: %s", err))
	}
	s.regex = regex
	if len(message) > 0 {
		s.regexMessage = message[0]
	}
	return s
}

// OneOf validates that the value is one of the provided options
func (s *StringSchema) OneOf(options ...string) *StringSchema {
	s.oneOf = options
	return s
}

// NotOneOf validates that the value is not one of the provided options
func (s *StringSchema) NotOneOf(options ...string) *StringSchema {
	s.notOneOf = options
	return s
}

// StartsWith validates that the string starts with the given prefix
func (s *StringSchema) StartsWith(prefix string) *StringSchema {
	s.startsWith = &prefix
	return s
}

// EndsWith validates that the string ends with the given suffix
func (s *StringSchema) EndsWith(suffix string) *StringSchema {
	s.endsWith = &suffix
	return s
}

// Includes validates that the string includes the given substring
func (s *StringSchema) Includes(substring string) *StringSchema {
	s.includes = &substring
	return s
}

// Refine adds a custom validation function
// The function receives the value and returns (isValid, errorMessage)
// If isValid is false, validation fails with the provided errorMessage
// If errorMessage is empty, a default message will be used
func (s *StringSchema) Refine(validator RefineFunc) *StringSchema {
	s.BaseSchema.addRefinement(validator)
	return s
}

// SuperRefine adds a super refinement validation function
// Similar to Refine, but provides a context object for adding errors with custom paths and codes
// This allows for more fine-grained control over error reporting, similar to Zod's superRefine
func (s *StringSchema) SuperRefine(validator SuperRefineFunc) *StringSchema {
	s.BaseSchema.addSuperRefinement(validator)
	return s
}

// Validate validates a value against the string schema
func (s *StringSchema) Validate(value any, path []any) *ValidationErrors {
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
	str, ok := value.(string)
	if !ok {
		msg := s.getErrorMessage(path, ErrCodeInvalidType, fmt.Sprintf("Expected string, got %T", value))
		errors.Add(path, ErrCodeInvalidType, msg)
		return errors
	}

	// Length validations
	if s.minLength != nil && len(str) < *s.minLength {
		msg := s.getErrorMessage(path, ErrCodeTooSmall, fmt.Sprintf("String must be at least %d character(s) long, got %d", *s.minLength, len(str)))
		errors.Add(path, ErrCodeTooSmall, msg)
	}

	if s.maxLength != nil && len(str) > *s.maxLength {
		msg := s.getErrorMessage(path, ErrCodeTooBig, fmt.Sprintf("String must be at most %d character(s) long, got %d", *s.maxLength, len(str)))
		errors.Add(path, ErrCodeTooBig, msg)
	}

	// Email validation
	if s.email {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(str) {
			msg := s.getErrorMessage(path, ErrCodeInvalidString, "Invalid email format")
			errors.Add(path, ErrCodeInvalidString, msg)
		}
	}

	// URL validation
	if s.url {
		urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
		if !urlRegex.MatchString(str) {
			msg := s.getErrorMessage(path, ErrCodeInvalidString, "Invalid URL format")
			errors.Add(path, ErrCodeInvalidString, msg)
		}
	}

	// Regex validation
	if s.regex != nil && !s.regex.MatchString(str) {
		message := s.regexMessage
		if message == "" {
			message = s.getErrorMessage(path, ErrCodeInvalidString, "String does not match required pattern")
		}
		errors.Add(path, ErrCodeInvalidString, message)
	}

	// OneOf validation
	if len(s.oneOf) > 0 {
		found := false
		for _, option := range s.oneOf {
			if str == option {
				found = true
				break
			}
		}
		if !found {
			msg := s.getErrorMessage(path, ErrCodeInvalidEnumValue, fmt.Sprintf("String must be one of: %s", strings.Join(s.oneOf, ", ")))
			errors.Add(path, ErrCodeInvalidEnumValue, msg)
		}
	}

	// NotOneOf validation
	if len(s.notOneOf) > 0 {
		for _, option := range s.notOneOf {
			if str == option {
				msg := s.getErrorMessage(path, ErrCodeInvalidEnumValue, fmt.Sprintf("String must not be one of: %s", strings.Join(s.notOneOf, ", ")))
				errors.Add(path, ErrCodeInvalidEnumValue, msg)
				break
			}
		}
	}

	// StartsWith validation
	if s.startsWith != nil && !strings.HasPrefix(str, *s.startsWith) {
		msg := s.getErrorMessage(path, ErrCodeInvalidString, fmt.Sprintf("String must start with '%s'", *s.startsWith))
		errors.Add(path, ErrCodeInvalidString, msg)
	}

	// EndsWith validation
	if s.endsWith != nil && !strings.HasSuffix(str, *s.endsWith) {
		msg := s.getErrorMessage(path, ErrCodeInvalidString, fmt.Sprintf("String must end with '%s'", *s.endsWith))
		errors.Add(path, ErrCodeInvalidString, msg)
	}

	// Includes validation
	if s.includes != nil && !strings.Contains(str, *s.includes) {
		msg := s.getErrorMessage(path, ErrCodeInvalidString, fmt.Sprintf("String must include '%s'", *s.includes))
		errors.Add(path, ErrCodeInvalidString, msg)
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
// This can be called on any schema type to customize error messages
func (s *StringSchema) CustomError(code, message string) *StringSchema {
	if s.BaseSchema.customErrors == nil {
		s.BaseSchema.customErrors = make(map[string]string)
	}
	s.BaseSchema.customErrors[code] = message
	return s
}

// SetErrorFormatter sets a custom error formatter function
// The formatter receives (path, code, defaultMessage) and returns the formatted message
func (s *StringSchema) SetErrorFormatter(formatter CustomErrorFunc) *StringSchema {
	s.BaseSchema.errorFormatter = formatter
	return s
}

// Type returns the schema type
func (s *StringSchema) Type() string {
	return "string"
}
