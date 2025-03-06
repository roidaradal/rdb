package rdb

import (
	"github.com/roidaradal/rdb/internal/condition"
	"github.com/roidaradal/rdb/internal/op"
)

type Condition = condition.Condition

func Equal[T any](key *T, value T) *condition.Value {
	return condition.NewValue(key, value, op.Equal)
}

func NotEqual[T any](key *T, value T) *condition.Value {
	return condition.NewValue(key, value, op.NotEqual)
}

func Greater[T any](key *T, value T) *condition.Value {
	return condition.NewValue(key, value, op.Greater)
}

func GreaterEqual[T any](key *T, value T) *condition.Value {
	return condition.NewValue(key, value, op.GreaterEqual)
}

func Less[T any](key *T, value T) *condition.Value {
	return condition.NewValue(key, value, op.Less)
}

func LessEqual[T any](key *T, value T) *condition.Value {
	return condition.NewValue(key, value, op.LessEqual)
}

func In[T any](key *T, values []T) *condition.List {
	return condition.NewList(key, values, op.In, op.Equal)
}

func NotIn[T any](key *T, values []T) *condition.List {
	return condition.NewList(key, values, op.NotIn, op.NotEqual)
}

func And(conditions ...Condition) *condition.Multi {
	return condition.NewMulti(op.And, conditions...)
}

func Or(conditions ...Condition) *condition.Multi {
	return condition.NewMulti(op.Or, conditions...)
}
