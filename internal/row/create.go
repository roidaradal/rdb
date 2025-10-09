package row

import (
	"reflect"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/memo"
)

// Creates a CreateRowFn for the given columns,
// Given an object reference, convert into map[string]any using the columns as keys
func Creator(columns []string) memo.CreateRowFn {
	// Ref needs to be a struct pointer
	return func(structRef any) dict.Object {
		emptyRow := dict.Object{}
		if !dyn.IsStructPointer(structRef) {
			return emptyRow
		}
		typeName := dyn.TypeOf(structRef)
		numColumns := len(columns)
		row := make(dict.Object, numColumns)
		for _, column := range columns {
			value, ok := GetColumnValue(structRef, typeName, column)
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

// From the given object reference and type name, get the field value for the given column,
// Object is expected to be a struct pointer
func GetColumnValue(structRef any, typeName, column string) (any, bool) {
	structValue := reflect.ValueOf(structRef).Elem()
	fieldName := memo.GetFieldName(typeName, column)
	if fieldName == "" {
		return nil, false
	}
	fieldValue := structValue.FieldByName(fieldName).Interface()
	return fieldValue, true
}
