package ze

import (
	"slices"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb"
)

type Schema[T any] struct {
	Name       string
	Ref        *T
	Table      string
	Reader     rdb.RowReader[T]
	validators map[string]ValidatorFn
	fieldsInfo
}

type ValidatorFn = func(any) bool

// Create new Schema
func NewSchema[T any](structRef *T, table string) (*Schema[T], error) {
	// Add rdb type
	err := rdb.AddType(structRef)
	if err != nil {
		return nil, err
	}

	// Get required, editable fields and field transformers
	fields := getFieldsInfo(structRef)

	schema := &Schema[T]{
		Name:       dyn.TypeOf(structRef),
		Ref:        structRef,
		Table:      table,
		Reader:     rdb.FullReader(structRef),
		fieldsInfo: fields,
		validators: make(map[string]ValidatorFn),
	}
	return schema, nil
}

// Create new shared Schema (no table)
func NewSharedSchema[T any](structRef *T) (*Schema[T], error) {
	return NewSchema(structRef, "")
}

// Additional required field for schema
func AddRequiredField[T, V any](schema *Schema[T], fieldRef *V) {
	fieldName := rdb.Field(schema.Name, fieldRef)
	if !slices.Contains(schema.required, fieldName) {
		schema.required = append(schema.required, fieldName)
	}
}

// Additional editable field for schema
func AddEditableField[T, V any](schema *Schema[T], fieldRef *V) {
	fieldName := rdb.Field(schema.Name, fieldRef)
	if !slices.Contains(schema.editable, fieldName) {
		schema.editable = append(schema.editable, fieldName)
	}
}

// Additional transformer for schema field
func AddTransformer[T, V any](schema *Schema[T], fieldRef *V, transformKey string) {
	fieldName := rdb.Field(schema.Name, fieldRef)
	if transformFn, ok := transformers[transformKey]; ok {
		schema.transformers[fieldName] = transformFn
	}
}

// Additional transformer function for schema field
func AddTransformFn[T, V any](schema *Schema[T], fieldRef *V, transformFn TransformFn) {
	fieldName := rdb.Field(schema.Name, fieldRef)
	schema.transformers[fieldName] = transformFn
}

// Add custom validator for schema field
func AddValidator[T, V any](schema *Schema[T], fieldRef *V, validator ValidatorFn) {
	fieldName := rdb.Field(schema.Name, fieldRef)
	schema.validators[fieldName] = validator
}

// Create a new string validator function
func NewStringValidator(validator func(string) bool) ValidatorFn {
	return func(item any) bool {
		// Assume item is a string
		text, ok := item.(string)
		if !ok {
			return false
		}
		return validator(text)
	}
}
