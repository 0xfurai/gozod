package gozod

// Schema is the base interface for all schemas
type Schema interface {
	Validate(value any, path []any) *ValidationErrors
	Type() string
}

// Shape defines the structure of a struct schema
// Maps struct field names (or JSON tag names) to validation schemas
type Shape map[string]Schema

// CustomErrorFunc is a function type for custom error formatting
// The path parameter is now a []any array
type CustomErrorFunc func(path []any, code, defaultMessage string) string

// RefineFunc is a function type for custom validation refinement
// Returns true if validation passes, false otherwise
// The second return value is the error message (optional, can be empty string)
type RefineFunc func(value any) (bool, string)

// SuperRefineContext provides methods to add validation errors with custom paths and codes
// Similar to Zod's superRefine context, allowing fine-grained control over error reporting
type SuperRefineContext struct {
	errors    *ValidationErrors
	basePath  []any
	formatter CustomErrorFunc
}

// AddIssue adds a validation error with a custom path, code, and message
// The path is relative to the base path of the schema being validated
func (ctx *SuperRefineContext) AddIssue(path []any, code, message string) {
	fullPath := ctx.basePath
	for _, part := range path {
		fullPath = PathAppend(fullPath, part)
	}
	ctx.errors.Add(fullPath, code, message)
}

// AddIssueWithMeta adds a validation error with metadata
func (ctx *SuperRefineContext) AddIssueWithMeta(path []any, code, message string, meta map[string]any) {
	fullPath := ctx.basePath
	for _, part := range path {
		fullPath = PathAppend(fullPath, part)
	}
	ctx.errors.AddWithMeta(fullPath, code, message, meta)
}

// SuperRefineFunc is a function type for super refinement validation
// Provides access to a context object for adding errors with custom paths and codes
type SuperRefineFunc func(value any, ctx *SuperRefineContext)

// getErrorMessage returns the custom error message if set, otherwise returns the default
func (b *BaseSchema) getErrorMessage(path []any, code, defaultMessage string) string {
	if b.errorFormatter != nil {
		return b.errorFormatter(path, code, defaultMessage)
	}
	if b.customErrors != nil {
		if customMsg, ok := b.customErrors[code]; ok {
			return customMsg
		}
	}
	return defaultMessage
}

// addRefinement adds a refinement function to the schema
// This is a helper method to avoid code duplication across schema types
func (b *BaseSchema) addRefinement(validator RefineFunc) {
	if b.refinements == nil {
		b.refinements = make([]RefineFunc, 0)
	}
	b.refinements = append(b.refinements, validator)
}

// applyRefinements applies all refine functions to the value
// Returns errors if any refinement fails
func (b *BaseSchema) applyRefinements(value any, path []any, errors *ValidationErrors) {
	if len(b.refinements) == 0 {
		return
	}

	for _, refine := range b.refinements {
		valid, message := refine(value)
		if !valid {
			if message == "" {
				message = b.getErrorMessage(path, ErrCodeCustomValidation, "Custom validation failed")
			} else {
				// Use the custom message but still use the error code
				message = b.getErrorMessage(path, ErrCodeCustomValidation, message)
			}
			errors.Add(path, ErrCodeCustomValidation, message)
		}
	}
}

// addSuperRefinement adds a super refinement function to the schema
func (b *BaseSchema) addSuperRefinement(validator SuperRefineFunc) {
	if b.superRefinements == nil {
		b.superRefinements = make([]SuperRefineFunc, 0)
	}
	b.superRefinements = append(b.superRefinements, validator)
}

// applySuperRefinements applies all super refine functions to the value
func (b *BaseSchema) applySuperRefinements(value any, path []any, errors *ValidationErrors) {
	if len(b.superRefinements) == 0 {
		return
	}

	for _, superRefine := range b.superRefinements {
		ctx := &SuperRefineContext{
			errors:    errors,
			basePath:  path,
			formatter: b.errorFormatter,
		}
		superRefine(value, ctx)
	}
}

// BaseSchema provides common functionality for all schemas
type BaseSchema struct {
	required         bool
	nilable          bool
	customErrors     map[string]string // Map of error code to custom message
	errorFormatter   func(path []any, code, defaultMessage string) string
	refinements      []RefineFunc      // Custom validation refinements
	superRefinements []SuperRefineFunc // Super refinement validations
}
