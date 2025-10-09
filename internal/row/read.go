package row

import (
	"errors"
	"reflect"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/memo"
)

// Function that reads row values into object
type RowReader[T any] = func(rowScanner) (*T, error)

// Interface for *sql.Row and *sql.Rows
type rowScanner interface {
	Scan(...any) error
}

// Creates a RowReader for type T, with the given columns
func Reader[T any](columns ...string) RowReader[T] {
	return func(row rowScanner) (*T, error) {
		var x T
		if !dyn.IsStruct(x) {
			return nil, errors.New("type is not a struct")
		}
		typeName := dyn.TypeOf(x)
		numColumns := len(columns)
		fieldRefs := make([]any, 0, numColumns)
		for _, column := range columns {
			if column == "" {
				continue // skip blank columns
			}
			fieldRef, ok := getColumnField(&x, typeName, column)
			if !ok {
				continue // skip if column field not found
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

// Creates a RowReader for type T, using all columns
func FullReader[T any](ref *T) RowReader[T] {
	columns := memo.ColumnsOf(ref)
	return Reader[T](columns...)
}

// From the given object reference and type name, get the reference to the corresponding column field,
// Object is expected to be a struct pointer
func getColumnField(structRef any, typeName, column string) (any, bool) {
	structValue := reflect.ValueOf(structRef).Elem()
	fieldName := memo.GetFieldName(typeName, column)
	if fieldName == "" {
		return nil, false
	}
	fieldRef := structValue.FieldByName(fieldName).Addr().Interface()
	return fieldRef, true
}
