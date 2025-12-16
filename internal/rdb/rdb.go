// Package rdb is the internal database of rdb; it contains saved types, columns, and fields information
package rdb

import (
	"errors"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/fn/list"
	"github.com/roidaradal/fn/str"
)

var (
	addressColumn    dict.StringMap            // {FieldAddress => ColumnName}
	typeColumns      dict.StringListMap        // {TypeName => []ColumnNames}
	typeColumnFields map[string]dict.StringMap // {TypeName => {ColumnName => FieldName}}
	typeFieldColumns map[string]dict.StringMap // {TypeName => {FieldName => ColumnName}}
	rowCreator       map[string]createRowFn    // {TypeName => CreateRowFn}
)

// Initialize memo data structures
func Initialize() {
	addressColumn = make(dict.StringMap)
	typeColumns = make(dict.StringListMap)
	typeColumnFields = make(map[string]dict.StringMap)
	typeFieldColumns = make(map[string]dict.StringMap)
	rowCreator = make(map[string]createRowFn)
}

// Add new type, extract its columns and fields.
// Return type name, list of columns.
// StructRef is expected to be a struct pointer
func AddType(structRef any) error {
	if !dyn.IsStructPointer(structRef) {
		return errors.New("type is not a struct pointer")
	}
	typeName := dyn.TypeOf(structRef)

	result := readStructColumns(structRef)
	addressColumn = dict.Update(addressColumn, result.addressColumn)
	typeColumns[typeName] = result.columns
	typeColumnFields[typeName] = result.columnFields
	typeFieldColumns[typeName] = result.fieldColumns
	rowCreator[typeName] = newRowCreator(typeName, result.columns)
	return nil
}

// Get all column names of given item's type
func ColumnsOf(item any) []string {
	typeName := dyn.TypeOf(item)
	return dict.DefaultGet(typeColumns, typeName, []string{})
}

// Get column name of given field reference,
// Field must be from the singleton object
func GetColumnName(fieldRef any) string {
	address := dyn.AddressOf(fieldRef)
	return addressColumn[address]
}

// Get column names of given field references,
// Fields must be from the singleton object
func GetColumnNames(fieldRefs ...any) []string {
	columns := list.Filter(list.Map(fieldRefs, GetColumnName), str.NotEmpty)
	if len(columns) != len(fieldRefs) {
		// return empty list if not all columns found
		return []string{}
	}
	return columns
}

// Get field name of given field reference,
// Field must be from the singleton object
func GetFieldName(typeName string, fieldRef any) string {
	return GetColumnFieldName(typeName, GetColumnName(fieldRef))
}

// Get field names of given field references,
// Fields must be from the singleton object
func GetFieldNames(typeName string, fieldRefs ...any) []string {
	fields := list.Map(fieldRefs, func(fieldRef any) string {
		return GetFieldName(typeName, fieldRef)
	})
	fields = list.Filter(fields, str.NotEmpty)
	if len(fields) != len(fieldRefs) {
		// return empty list if not all fields found
		return []string{}
	}
	return fields
}

// Get field name for given type name's column
func GetColumnFieldName(typeName, columnName string) string {
	if dict.NoKey(typeColumnFields, typeName) {
		return ""
	}
	return typeColumnFields[typeName][columnName]
}

// Get column name for given type name's field
func getFieldColumnName(typeName, fieldName string) string {
	if dict.NoKey(typeFieldColumns, typeName) {
		return ""
	}
	return typeFieldColumns[typeName][fieldName]
}
