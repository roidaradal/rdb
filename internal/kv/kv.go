package kv

import "github.com/roidaradal/rdb/internal/memo"

type Value struct {
	column string
	value  any
}

type List struct {
	column string
	values []any
}

func (v Value) Get() (string, any) {
	return v.column, v.value
}

func (l List) Get() (string, []any) {
	return l.column, l.values
}

func KeyValue[T any](key *T, value T) *Value {
	column := memo.GetColumn(key)
	if column == "" {
		return nil
	}
	return &Value{
		column: column,
		value:  value,
	}
}

func KeyList[T any](key *T, values []T) *List {
	column := memo.GetColumn(key)
	if column == "" {
		return nil
	}
	values2 := make([]any, len(values))
	for i, value := range values {
		values2[i] = value
	}
	return &List{
		column: column,
		values: values2,
	}
}
