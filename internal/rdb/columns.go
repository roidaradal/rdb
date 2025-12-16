package rdb

import (
	"fmt"
	"reflect"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/fn/str"
)

const (
	columnTag    string = "col" // struct tag used to defined column name
	skipTagValue string = "-"   // set col:"-" to skip column
)

type columnsInfo struct {
	columns       []string       // list of column names
	columnFields  dict.StringMap // {ColumnName => FieldName}
	fieldColumns  dict.StringMap // {FieldName => ColumnName}
	addressColumn dict.StringMap // {FieldAddress => ColumnName}
}

// From given struct reference and type name, get field value for given column
// StructRef is assumed to be validated as struct pointer
func GetStructColumnValue(structRef any, typeName, columnName string) (any, bool) {
	structValue := reflect.ValueOf(structRef).Elem() // dereference pointer
	fieldName := GetColumnFieldName(typeName, columnName)
	if fieldName == "" {
		return nil, false
	}
	fieldValue := structValue.FieldByName(fieldName).Interface()
	return fieldValue, true
}

// Get all columns and field names from given struct pointer
func readStructColumns(structRef any) *columnsInfo {
	result := &columnsInfo{
		columns:       make([]string, 0),
		columnFields:  make(dict.StringMap),
		fieldColumns:  make(dict.StringMap),
		addressColumn: make(dict.StringMap),
	}

	if !dyn.IsStructPointer(structRef) {
		return result
	}

	structValue := reflect.ValueOf(structRef).Elem() // dereference pointer
	structType := structValue.Type()
	numFields := structType.NumField()

	for idx := range numFields {
		structField := structType.Field(idx)
		fieldName := structField.Name
		if structField.Anonymous {
			// Embedded struct, get columns using recursion
			innerStructRef := structValue.FieldByName(fieldName).Addr().Interface()
			inner := readStructColumns(innerStructRef)
			result.columns = append(result.columns, inner.columns...)
			result.addressColumn = dict.Update(result.addressColumn, inner.addressColumn)
			result.columnFields = dict.Update(result.columnFields, inner.columnFields)
			result.fieldColumns = dict.Update(result.fieldColumns, inner.fieldColumns)
		} else {
			// Normal field
			column := extractColumnName(structField)
			if column == "" {
				continue // skip blank columns
			}
			column = str.WrapBackticks(column)
			fieldAddress := fmt.Sprintf("%v", structValue.Field(idx).Addr().Interface())
			result.columns = append(result.columns, column)
			result.addressColumn[fieldAddress] = column
			result.columnFields[column] = fieldName
			result.fieldColumns[fieldName] = column
		}
	}
	return result
}

// Extract custom column name from struct tag if not skipped,
// Defaults to field name
func extractColumnName(structField reflect.StructField) string {
	column := structField.Tag.Get(columnTag)
	if column == skipTagValue {
		// no column if skipped
		return ""
	} else if column != "" {
		// use custom column name if not blank
		return column
	}
	// default: field name
	return structField.Name
}
