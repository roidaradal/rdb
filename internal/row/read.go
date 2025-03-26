package row

import (
	"reflect"

	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/types"
)

type RowScanner interface {
	Scan(...any) error
}

type RowReader[T any] = func(RowScanner) (*T, error)

func Reader[T any](columns []string) RowReader[T] {
	return func(row RowScanner) (*T, error) {
		var x T
		if !types.IsStruct(x) {
			return nil, errNotStruct
		}
		schema := types.NameOf(x)
		numColumns := len(columns)
		fields := make([]any, 0, numColumns)
		for _, column := range columns {
			field, ok := findColumnField(&x, schema, column)
			if !ok {
				continue
			}
			fields = append(fields, field)
		}
		if len(fields) != numColumns {
			return nil, errIncompleteFields
		}
		err := row.Scan(fields...)
		return &x, err
	}
}

func FullReader[T any](schema *T) RowReader[T] {
	columns := memo.ColumnsOf(schema)
	return Reader[T](columns)
}

func findColumnField(x any, schema, column string) (any, bool) {
	structValue := reflect.ValueOf(x).Elem()
	fieldName := memo.GetFieldName(schema, column)
	if fieldName == "" {
		return nil, false
	}
	fieldAddress := structValue.FieldByName(fieldName).Addr().Interface()
	return fieldAddress, true
}
