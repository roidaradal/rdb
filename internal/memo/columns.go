package memo

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/fn/str"
)

const (
	columnTag    string = "col"
	skipTagValue string = "-"
)

type columnsResult struct {
	columns      []string
	columnFields dict.StringMap
	fieldColumns dict.StringMap
	addressOf    dict.StringMap
}

// Add a new type, extract its columns and fields,
// StructRef is expected to be a struct pointer
func AddType(structRef any) error {
	if !dyn.IsStructPointer(structRef) {
		return errors.New("type is not a struct pointer")
	}
	typeName := dyn.TypeOf(structRef)

	result := getAllColumns(structRef)
	allColumns[typeName] = result.columns
	typeColumnFields[typeName] = result.columnFields
	typeFieldColumns[typeName] = result.fieldColumns
	columnAddress = dict.Update(columnAddress, result.addressOf)
	return nil
}

// Get all columns and field names from the given struct pointer
func getAllColumns(structRef any) *columnsResult {
	result := &columnsResult{
		columns:      make([]string, 0),
		columnFields: make(dict.StringMap),
		fieldColumns: make(dict.StringMap),
		addressOf:    make(dict.StringMap),
	}

	if !dyn.IsStructPointer(structRef) {
		return result
	}

	structValue := reflect.ValueOf(structRef).Elem() // dereference pointer
	structType := structValue.Type()
	numFields := structType.NumField()

	for i := range numFields {
		structField := structType.Field(i)
		fieldName := structField.Name
		if structField.Anonymous {
			// Embedded struct, get columns using recursion
			embeddedStructRef := structValue.FieldByName(fieldName).Addr().Interface()
			embedded := getAllColumns(embeddedStructRef)
			result.columns = append(result.columns, embedded.columns...)
			result.addressOf = dict.Update(result.addressOf, embedded.addressOf)
			result.columnFields = dict.Update(result.columnFields, embedded.columnFields)
			result.fieldColumns = dict.Update(result.fieldColumns, embedded.fieldColumns)
		} else {
			// Normal field
			column := getColumnName(structField)
			if column == "" {
				continue // skip blank columns
			}
			column = str.WrapBackticks(column)
			fieldAddress := getFieldAddress(structValue, i)
			result.columns = append(result.columns, column)
			result.addressOf[fieldAddress] = column
			result.columnFields[column] = fieldName
			result.fieldColumns[fieldName] = column
		}
	}
	return result
}

// Extract custom column name from struct tag if not skipped,
// Defaults to field name
func getColumnName(structField reflect.StructField) string {
	column := structField.Tag.Get(columnTag)
	if column == skipTagValue {
		return ""
	} else if column != "" {
		return column
	}
	return structField.Name
}

// Get string address of given struct's index-th field
func getFieldAddress(structValue reflect.Value, fieldIndex int) string {
	return fmt.Sprintf("%v", structValue.Field(fieldIndex).Addr().Interface())
}
