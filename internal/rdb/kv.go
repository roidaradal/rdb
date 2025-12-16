package rdb

import "github.com/roidaradal/fn/list"

// Key-Value pair; key = column
type Value struct {
	column string
	value  any
}

// Key-Values pair; key = column, values = list
type List struct {
	column string
	values []any
}

// Return Value's column and value
func (v Value) Tuple() (string, any) {
	return v.column, v.value
}

// Return List's column and list
func (l List) Tuple() (string, []any) {
	return l.column, l.values
}

// Create new KeyValue pair
func KeyValue[T any](key *T, value T) *Value {
	column := GetColumnName(key)
	if column == "" {
		return nil
	}
	return &Value{column, value}
}

// Create new KeyList pair
func KeyList[T any](key *T, values []T) *List {
	column := GetColumnName(key)
	if column == "" {
		return nil
	}
	values2 := list.ToAny(values)
	return &List{column, values2}
}

// Create new KeyValue pair, get column from fieldName
func ColumnValue(typeName, fieldName string, value any) *Value {
	column := getFieldColumnName(typeName, fieldName)
	if column == "" {
		return nil
	}
	return &Value{column, value}
}
