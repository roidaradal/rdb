package kv

import "github.com/roidaradal/rdb/internal/memo"

// Key-Value pair, key = column
type Value struct {
	column string
	value  any
}

// Key-Values pair, key = column, values = list
type List struct {
	column string
	values []any
}

// Return column and value for destructuring
func (v Value) Tuple() (string, any) {
	return v.column, v.value
}

// Return column and values for destructuring
func (l List) Tuple() (string, []any) {
	return l.column, l.values
}

// Creates a new KeyValue pair, from fieldName => column
func ColumnValue(typeName, fieldName string, value any) *Value {
	column := memo.GetColumnName(typeName, fieldName)
	if column == "" {
		return nil
	}
	return &Value{column, value}
}

// Creates a new KeyValue pair
func KeyValue[T any](key *T, value T) *Value {
	column := memo.GetColumn(key)
	if column == "" {
		return nil
	}
	return &Value{column, value}
}

// Creates a new KeyList pair
func KeyList[T any](key *T, values []T) *List {
	column := memo.GetColumn(key)
	if column == "" {
		return nil
	}
	values2 := make([]any, len(values))
	for i, value := range values {
		values2[i] = value
	}
	return &List{column, values2}
}
