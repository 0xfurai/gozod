package gozod

import (
	"fmt"
	"reflect"
	"strings"
)

// StructSchema validates struct values directly
type StructSchema struct {
	BaseSchema
	shape  map[string]Schema // Maps struct field names (or JSON tag names) to schemas
	strict bool              // If true, rejects unknown fields (default: false, allows extra fields)
}

// Struct creates a new struct schema
// The shape maps struct field names (or JSON tag names) to validation schemas
func Struct(shape Shape) *StructSchema {
	return &StructSchema{
		BaseSchema: BaseSchema{required: true},
		shape:      map[string]Schema(shape),
		strict:     false, // Extra fields allowed by default
	}
}

// Nilable allows null/nil values
func (s *StructSchema) Nilable() *StructSchema {
	s.nilable = true
	return s
}

// Strict rejects unknown fields
func (s *StructSchema) Strict() *StructSchema {
	s.strict = true
	return s
}

// Refine adds a custom validation function
// The function receives the value and returns (isValid, errorMessage)
// If isValid is false, validation fails with the provided errorMessage
// If errorMessage is empty, a default message will be used
func (s *StructSchema) Refine(validator RefineFunc) *StructSchema {
	s.BaseSchema.addRefinement(validator)
	return s
}

// SuperRefine adds a super refinement validation function
// Similar to Refine, but provides a context object for adding errors with custom paths and codes
// This allows for more fine-grained control over error reporting, similar to Zod's superRefine
func (s *StructSchema) SuperRefine(validator SuperRefineFunc) *StructSchema {
	s.BaseSchema.addSuperRefinement(validator)
	return s
}

// Validate validates a struct value against the schema
func (s *StructSchema) Validate(value any, path []any) *ValidationErrors {
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

	// Get reflect value, handling pointers
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			if !s.nilable {
				msg := s.getErrorMessage(path, ErrCodeRequired, "Required")
				errors.Add(path, ErrCodeRequired, msg)
				return errors
			}
			return nil
		}
		val = val.Elem()
	}

	// Must be a struct
	if val.Kind() != reflect.Struct {
		msg := s.getErrorMessage(path, ErrCodeInvalidType, fmt.Sprintf("Expected struct, got %T", value))
		errors.Add(path, ErrCodeInvalidType, msg)
		return errors
	}

	typ := val.Type()

	// Build a map of field names to their reflect.Field for quick lookup
	// This handles both struct field names and JSON tag names
	fieldMap := make(map[string]reflect.StructField)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}

		// Get the field name that will be used in the schema
		fieldName := getStructFieldName(field)
		fieldMap[fieldName] = field
	}

	// Validate each field in the shape
	for schemaFieldName, schema := range s.shape {
		fieldPath := PathAppend(path, schemaFieldName)

		// Find the struct field by schema field name
		structField, exists := fieldMap[schemaFieldName]
		if !exists {
			// Field not found in struct - this is a validation error
			// (unless it's nilable, but we still need to validate it)
			fieldErrors := schema.Validate(nil, fieldPath)
			if fieldErrors != nil {
				errors.Errors = append(errors.Errors, fieldErrors.Errors...)
			}
			continue
		}

		// Get the field value
		structFieldName := structField.Name
		fieldValue := val.FieldByName(structFieldName)
		if !fieldValue.IsValid() {
			// Field exists but can't be accessed
			fieldErrors := schema.Validate(nil, fieldPath)
			if fieldErrors != nil {
				errors.Errors = append(errors.Errors, fieldErrors.Errors...)
			}
			continue
		}

		// Get the interface value, handling pointers
		var fieldInterface any
		fieldType := structField.Type

		if fieldType.Kind() == reflect.Ptr {
			// Field is a pointer type
			if fieldValue.IsNil() {
				fieldInterface = nil
			} else {
				// Dereference the pointer
				fieldInterface = fieldValue.Elem().Interface()
			}
		} else {
			// Field is not a pointer
			if fieldValue.CanInterface() {
				fieldInterface = fieldValue.Interface()
			} else {
				// Field can't be interfaced, treat as nil
				fieldInterface = nil
			}
		}

		// Check for zero values with omitempty
		jsonTag := structField.Tag.Get("json")
		hasOmitempty := strings.Contains(jsonTag, "omitempty")
		if hasOmitempty && isEmptyValue(fieldInterface) {
			// For omitempty fields, pass nil to validation
			fieldInterface = nil
		}

		// Validate the field
		fieldErrors := schema.Validate(fieldInterface, fieldPath)
		if fieldErrors != nil {
			errors.Errors = append(errors.Errors, fieldErrors.Errors...)
		}
	}

	// Check for unknown fields if strict mode is enabled
	if s.strict {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if !field.IsExported() {
				continue
			}

			fieldName := getStructFieldName(field)
			if _, exists := s.shape[fieldName]; !exists {
				keyPath := PathAppend(path, fieldName)
				msg := s.getErrorMessage(keyPath, ErrCodeUnrecognizedKeys, fmt.Sprintf("Unrecognized field '%s'", fieldName))
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
func (s *StructSchema) CustomError(code, message string) *StructSchema {
	if s.BaseSchema.customErrors == nil {
		s.BaseSchema.customErrors = make(map[string]string)
	}
	s.BaseSchema.customErrors[code] = message
	return s
}

// SetErrorFormatter sets a custom error formatter function
func (s *StructSchema) SetErrorFormatter(formatter CustomErrorFunc) *StructSchema {
	s.BaseSchema.errorFormatter = formatter
	return s
}

// Type returns the schema type
func (s *StructSchema) Type() string {
	return "struct"
}

// getStructFieldName extracts the field name from a struct field
// It prioritizes JSON tags, then falls back to the struct field name
// Handles JSON tag options like "omitempty"
func getStructFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		// Extract the field name from JSON tag (before any comma)
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			return jsonTag[:idx]
		}
		return jsonTag
	}

	// No JSON tag, use struct field name (convert to camelCase for consistency)
	fieldName := field.Name
	if len(fieldName) > 0 {
		return strings.ToLower(fieldName[:1]) + fieldName[1:]
	}
	return fieldName
}
