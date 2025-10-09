package rdb

import "github.com/roidaradal/rdb/internal/condition"

// Condition interface
type Condition = condition.Condition

// Create Equal condition
func Equal[T any](fieldRef *T, value T) *condition.Value {
	return condition.NewValue(fieldRef, value, condition.Equal)
}

// Create NotEqual condition
func NotEqual[T any](fieldRef *T, value T) *condition.Value {
	return condition.NewValue(fieldRef, value, condition.NotEqual)
}

// Create Prefix condition
func Prefix(fieldRef *string, value string) *condition.Value {
	return condition.NewValue(fieldRef, value, condition.Prefix)
}

// Create Suffix condition
func Suffix(fieldRef *string, value string) *condition.Value {
	return condition.NewValue(fieldRef, value, condition.Suffix)
}

// Create Substring condition
func Substring(fieldRef *string, value string) *condition.Value {
	return condition.NewValue(fieldRef, value, condition.Substring)
}

// Create Greater condition
func Greater[T any](fieldRef *T, value T) *condition.Value {
	return condition.NewValue(fieldRef, value, condition.Greater)
}

// Create GreaterEqual condition
func GreaterEqual[T any](fieldRef *T, value T) *condition.Value {
	return condition.NewValue(fieldRef, value, condition.GreaterEqual)
}

// Create Less condition
func Less[T any](fieldRef *T, value T) *condition.Value {
	return condition.NewValue(fieldRef, value, condition.Less)
}

// Create LessEqual condition
func LessEqual[T any](fieldRef *T, value T) *condition.Value {
	return condition.NewValue(fieldRef, value, condition.LessEqual)
}

// Create In condition
func In[T any](fieldRef *T, values []T) *condition.List {
	return condition.NewList(fieldRef, values, condition.In, condition.Equal)
}

// Create NotIn condition
func NotIn[T any](fieldRef *T, values []T) *condition.List {
	return condition.NewList(fieldRef, values, condition.NotIn, condition.NotEqual)
}

// Create And condition
func And(conditions ...Condition) *condition.Multi {
	return condition.NewMulti(condition.And, conditions...)
}

// Create Or condition
func Or(conditions ...Condition) *condition.Multi {
	return condition.NewMulti(condition.Or, conditions...)
}

// Create NoCondition (match all)
func NoCondition() *condition.MatchAll {
	return &condition.MatchAll{}
}
