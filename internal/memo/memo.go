package memo

import (
	"github.com/roidaradal/fn"
	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
)

var (
	allColumns       dict.StringListMap        // {TypeName => []ColumnNames}
	columnAddress    dict.StringMap            // {FieldAddress => ColumnName}
	typeColumnFields map[string]dict.StringMap // {TypeName => {ColumnName => FieldName}}
	typeFieldColumns map[string]dict.StringMap // {TypeName => {FieldName => ColumnName}}
	rowCreator       map[string]CreateRowFn    // {TypeName => CreateRowFn}
)

// Converts object to map[string]any for row insertion
type CreateRowFn func(any) dict.Object

// Initialize the memo data structures
func Initialize() {
	allColumns = make(dict.StringListMap)
	columnAddress = make(dict.StringMap)
	typeColumnFields = make(map[string]dict.StringMap)
	typeFieldColumns = make(map[string]dict.StringMap)
	rowCreator = make(map[string]CreateRowFn)
}

// Adds a createRowFn for the given typeName
func AddRowCreator(typeName string, createRowFn CreateRowFn) {
	rowCreator[typeName] = createRowFn
}

// Get the createRowFn associated with given typeName
func GetRowCreator(typeName string) (CreateRowFn, bool) {
	rowFn, ok := rowCreator[typeName]
	return rowFn, ok
}

// Get type name and all column names of given object
func TypeColumnsOf(item any) (string, []string) {
	typeName := dyn.TypeOf(item)
	columns, ok := allColumns[typeName]
	if !ok {
		columns = []string{}
	}
	return typeName, columns
}

// Get all column names of given item
func ColumnsOf(item any) []string {
	_, columns := TypeColumnsOf(item)
	return columns
}

// Get column name of given field reference,
// Field must be from the singleton object
func GetColumn(fieldRef any) string {
	address := dyn.AddressOf(fieldRef)
	return columnAddress[address]
}

// Get column names of given field references,
// Fields must be from the singleton object
func GetColumns(fieldRefs ...any) []string {
	numFields := len(fieldRefs)
	columns := make([]string, 0, numFields)
	for _, fieldRef := range fieldRefs {
		column := GetColumn(fieldRef)
		if column != "" {
			columns = append(columns, column)
		}
	}
	if len(columns) != numFields {
		// reset to empty list if not all columns were found
		columns = []string{}
	}
	return columns
}

// Get field names of given field references,
// Fields must be from the singleton object
func GetFields(typeName string, fieldRefs ...any) []string {
	columns := GetColumns(fieldRefs...)
	fields := fn.Map(columns, func(column string) string {
		return GetFieldName(typeName, column)
	})
	fields = fn.Filter(fields, check.NotEmptyString)
	return fields
}

// Get field name for given type name's column
func GetFieldName(typeName, column string) string {
	if dict.NoKey(typeColumnFields, typeName) {
		return ""
	}
	return typeColumnFields[typeName][column]
}

// Get column name for given type name's field
func GetColumnName(typeName, fieldName string) string {
	if dict.NoKey(typeFieldColumns, typeName) {
		return ""
	}
	return typeFieldColumns[typeName][fieldName]
}
