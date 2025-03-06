package row

import (
	"reflect"

	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/types"
)

func Creator(columns []string) types.RowFn {
	return func(x any) map[string]any {
		emptyMap := map[string]any{}
		if !types.IsStructPointer(x) {
			return emptyMap
		}
		schema := types.NameOf(x)
		numColumns := len(columns)
		row := make(map[string]any, numColumns)
		for _, column := range columns {
			value, ok := FindColumnValue(x, schema, column)
			if !ok {
				continue
			}
			row[column] = value
		}
		if len(row) != numColumns {
			return emptyMap
		}
		return row
	}
}

func FindColumnValue(x any, schema, column string) (any, bool) {
	structValue := reflect.ValueOf(x).Elem()
	fieldName := memo.GetFieldName(schema, column)
	if fieldName == "" {
		return nil, false
	}
	fieldValue := structValue.FieldByName(fieldName).Interface()
	return fieldValue, true
}
