package rdb

import (
	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/query"
)

// Query interface
type Query = query.Query

// Builds the full query as string
var QueryString = query.QueryString

// Creates a new CountQuery
func NewCountQuery(table string) *query.CountQuery {
	q := &query.CountQuery{}
	q.Initialize(table)
	return q
}

// Creates a new DeleteQuery
func NewDeleteQuery(table string) *query.DeleteQuery {
	q := &query.DeleteQuery{}
	q.Initialize(table)
	return q
}

// Creates a new UpdateQuery
func NewUpdateQuery(table string) *query.UpdateQuery {
	q := &query.UpdateQuery{}
	q.Initialize(table)
	return q
}

// Add key = value update to UpdateQuery
func Update[T any](q *query.UpdateQuery, fieldRef *T, value T) {
	query.Update(q, fieldRef, value)
}

// Creates a new InsertRowQuery
func NewInsertRowQuery(table string) *query.InsertRowQuery {
	q := &query.InsertRowQuery{}
	q.Initialize(table)
	return q
}

// Creates a new InsertRowsQuery
func NewInsertRowsQuery(table string) *query.InsertRowsQuery {
	q := &query.InsertRowsQuery{}
	q.Initialize(table)
	return q

}

// Creates a new ValueQuery
func NewValueQuery[T any, V any](table string, fieldRef *V) *query.ValueQuery[T, V] {
	q := &query.ValueQuery[T, V]{}
	q.Initialize(table, fieldRef)
	return q
}

// Creates a new LookupQuery
func NewLookupQuery[T any, K comparable, V any](table string, keyFieldRef *K, valueFieldRef *V) *query.LookupQuery[T, K, V] {
	q := &query.LookupQuery[T, K, V]{}
	q.Initialize(table, keyFieldRef, valueFieldRef)
	return q
}

// Creates a new DistinctValuesQuery
func NewDistinctValuesQuery[T any, V any](table string, fieldRef *V) *query.DistinctValuesQuery[T, V] {
	q := &query.DistinctValuesQuery[T, V]{}
	q.Initialize(table, fieldRef)
	return q
}

// Creates a new SelectRowQuery with selected columns (set later)
func NewSelectRowQuery[T any](table string, reader RowReader[T]) *query.SelectRowQuery[T] {
	q := &query.SelectRowQuery[T]{}
	q.Initialize(table, reader)
	return q
}

// Creates a new SelectRowQuery that uses all columns
func NewFullSelectRowQuery[T any](table string, reader RowReader[T]) *query.SelectRowQuery[T] {
	var t T
	q := &query.SelectRowQuery[T]{}
	q.Initialize(table, reader)
	q.Columns(memo.ColumnsOf(t))
	return q
}

// Creates a new SelectRowsQuery with selected columns (set later)
func NewSelectRowsQuery[T any](table string, reader RowReader[T]) *query.SelectRowsQuery[T] {
	q := &query.SelectRowsQuery[T]{}
	q.Initialize(table, reader)
	return q
}

// Creates a new SelectRowsQuery that uses all columns
func NewFullSelectRowsQuery[T any](table string, reader RowReader[T]) *query.SelectRowsQuery[T] {
	var t T
	q := &query.SelectRowsQuery[T]{}
	q.Initialize(table, reader)
	q.Columns(memo.ColumnsOf(t))
	return q
}

// Creates a new TopRowQuery
func NewTopRowQuery[T any](table string, reader RowReader[T]) *query.TopRowQuery[T] {
	q := &query.TopRowQuery[T]{}
	q.Initialize(table, reader)
	return q
}

// Creates a new TopValueQuery
func NewTopValueQuery[T any, V any](table string, fieldRef *V) *query.TopValueQuery[T, V] {
	q := &query.TopValueQuery[T, V]{}
	q.Initialize(table, fieldRef)
	return q
}
