package rdb

import (
	"errors"
	"reflect"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
)

// Interface for *sql.Row and *sql.Rows
type rowScanner interface {
	Scan(...any) error
}

// Function that reads row values into struct
type RowReader[T any] = func(rowScanner) (*T, error)

// Converts structRef to map[string]any for row insertion
type createRowFn func(any) dict.Object

// Creates a CreateRowFn for given columns:
// Given struct reference, convert into map[string]any using columns as keys
func newRowCreator(typeName string, columns []string) createRowFn {
	return func(structRef any) dict.Object {
		emptyRow := dict.Object{}
		if !dyn.IsStructPointer(structRef) {
			return emptyRow
		}
		numColumns := len(columns)
		row := make(dict.Object, numColumns)
		for _, column := range columns {
			value, ok := GetStructColumnValue(structRef, typeName, column)
			if !ok {
				continue // skip if column value not found
			}
			row[column] = value
		}
		if len(row) != numColumns {
			return emptyRow // return empty if some columns failed
		}
		return row
	}
}

// Convert given struct to map[string]any for row insertion
func ToRow[T any](structRef *T) dict.Object {
	typeName := dyn.TypeOf(structRef)
	rowFn, ok := rowCreator[typeName]
	if !ok {
		return dict.Object{}
	}
	return rowFn(structRef)
}

// Creates a RowReader for type T, using all columns
func FullReader[T any](structRef *T) RowReader[T] {
	columns := ColumnsOf(structRef)
	return NewReader[T](columns...)
}

// Create new RowReader for type T, with given columns
func NewReader[T any](columns ...string) RowReader[T] {
	return func(row rowScanner) (*T, error) {
		var x T
		if !dyn.IsStruct(x) {
			return nil, errors.New("type is not struct")
		}
		typeName := dyn.TypeOf(x)
		numColumns := len(columns)
		fieldRefs := make([]any, 0, numColumns)
		for _, column := range columns {
			if column == "" {
				continue // skip blank columns
			}
			fieldRef, ok := getColumnFieldRef(&x, typeName, column)
			if !ok {
				continue // skip if column's field not found
			}
			fieldRefs = append(fieldRefs, fieldRef)
		}
		if len(fieldRefs) != numColumns {
			// return nil if some columns failed
			return nil, errors.New("incomplete fields")
		}
		err := row.Scan(fieldRefs...)
		return &x, err
	}
}

// Get reference to corresponding column field of given type and struct reference
func getColumnFieldRef(structRef any, typeName, columnName string) (any, bool) {
	structValue := reflect.ValueOf(structRef).Elem() // dereference pointer
	fieldName := GetColumnFieldName(typeName, columnName)
	if fieldName == "" {
		return nil, false
	}
	fieldRef := structValue.FieldByName(fieldName).Addr().Interface()
	return fieldRef, true
}
