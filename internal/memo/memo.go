package memo

import "github.com/roidaradal/rdb/internal/types"

var AllColumns map[string][]string           // {SchemaName => []ColumnNames}
var ColumnAddress map[string]string          // {FieldAddress => ColumnName}
var ColumnField map[string]map[string]string // {SchemaName => {ColumnName => FieldName}}
var Instance map[string]any                  // {SchemaName => *T singleton}
var RowCreator map[string]types.RowFn        // {SchemaName => RowCreatorFn}

func Initialize() {
	AllColumns = make(map[string][]string)
	ColumnAddress = make(map[string]string)
	ColumnField = make(map[string]map[string]string)
	Instance = make(map[string]any)
	RowCreator = make(map[string]types.RowFn)
}

func UpdateColumnAddress(columns map[string]string) {
	for k, v := range columns {
		ColumnAddress[k] = v
	}
}

func InstanceOf[T any](t T) *T {
	name := types.NameOf(t)
	instance, ok := Instance[name]
	if !ok {
		return nil
	}
	return instance.(*T)
}

func ColumnsOf(schema any) []string {
	name := types.NameOf(schema)
	columns, ok := AllColumns[name]
	if !ok {
		return []string{}
	}
	return columns
}

func GetColumn(field any) string {
	address := types.AddressOf(field)
	return ColumnAddress[address]
}

func GetColumns(fields ...any) []string {
	numFields := len(fields)
	columns := make([]string, 0, numFields)
	for _, field := range fields {
		column := GetColumn(field)
		if column != "" {
			columns = append(columns, column)
		}
	}
	if len(columns) != numFields {
		return []string{}
	}
	return columns
}

func GetFieldName(schema, column string) string {
	if _, ok := ColumnField[schema]; !ok {
		return ""
	}
	return ColumnField[schema][column]
}
