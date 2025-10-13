package ze

import (
	"slices"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb"
	"github.com/roidaradal/rdb/internal/memo"
)

type Schema[T any] struct {
	Name         string
	Ref          *T
	Table        string
	Reader       rdb.RowReader[T]
	required     []string
	editable     []string
	transformers map[string]memo.TransformFn
}

// Creates a new schema
func NewSchema[T any](structRef *T, table string) (*Schema[T], error) {
	// Add rdb type
	err := rdb.AddType(structRef)
	if err != nil {
		return nil, err
	}

	// Get required, editable fields and field transformers
	fields := memo.GetFieldsInfo(structRef)

	schema := &Schema[T]{
		Name:         dyn.TypeOf(structRef),
		Ref:          structRef,
		Table:        table,
		Reader:       rdb.FullReader(structRef),
		required:     fields.Required,
		editable:     fields.Editable,
		transformers: fields.Transformers,
	}
	return schema, nil
}

// Creates a new shared schema (no table)
func NewSharedSchema[T any](structRef *T) (*Schema[T], error) {
	return NewSchema(structRef, "")
}

// Additional required field for schema
func AddRequiredField[T any, V any](schema *Schema[T], fieldRef *V) {
	typeName := schema.Name
	column := rdb.Column(fieldRef)
	fieldName := memo.GetFieldName(typeName, column)
	if !slices.Contains(schema.required, fieldName) {
		schema.required = append(schema.required, fieldName)
	}
}

// Additional editable field for schema
func AddEditableField[T any, V any](schema *Schema[T], fieldRef *V) {
	typeName := schema.Name
	column := rdb.Column(fieldRef)
	fieldName := memo.GetFieldName(typeName, column)
	if !slices.Contains(schema.editable, fieldName) {
		schema.editable = append(schema.editable, fieldName)
	}
}

// Additional transformer for schema field
func AddTransformer[T any, V any](schema *Schema[T], fieldRef *V, transformKey string) {
	typeName := schema.Name
	column := rdb.Column(fieldRef)
	fieldName := memo.GetFieldName(typeName, column)
	if transform, ok := memo.Transformer(transformKey); ok {
		schema.transformers[fieldName] = transform
	}
}
