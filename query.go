package rdb

import (
	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/query"
	"github.com/roidaradal/rdb/internal/row"
)

type Query = query.Query

func NewCountQuery(table string) *query.CountQuery {
	q := query.CountQuery{}
	q.Initialize(table)
	return &q
}

func NewDeleteQuery(table string) *query.DeleteQuery {
	q := query.DeleteQuery{}
	q.Initialize(table)
	return &q
}

func NewUpdateQuery(table string) *query.UpdateQuery {
	q := query.UpdateQuery{}
	q.Initialize(table)
	return &q
}

func Update[T any](q *query.UpdateQuery, key *T, value T) {
	query.Update(q, key, value)
}

func NewInsertRowQuery(table string) *query.InsertRowQuery {
	q := query.InsertRowQuery{}
	q.Initialize(table)
	return &q
}

func NewInsertRowsQuery(table string) *query.InsertRowsQuery {
	q := query.InsertRowsQuery{}
	q.Initialize(table)
	return &q
}

func NewValueQuery[T any, V any](table string, field *V) *query.ValueQuery[T, V] {
	q := query.ValueQuery[T, V]{}
	q.Initialize(table, field)
	return &q
}

func NewLookupQuery[T any, K comparable, V any](table string, key *K, value *V) *query.LookupQuery[T, K, V] {
	q := query.LookupQuery[T, K, V]{}
	q.Initialize(table, key, value)
	return &q
}

func NewSelectRowQuery[T any](table string, reader row.RowReader[T]) *query.SelectRowQuery[T] {
	q := query.SelectRowQuery[T]{}
	q.Initialize(table, reader)
	return &q
}

func NewFullSelectRowQuery[T any](table string, reader row.RowReader[T]) *query.SelectRowQuery[T] {
	var t T
	q := query.SelectRowQuery[T]{}
	q.Initialize(table, reader)
	q.Columns(memo.ColumnsOf(t))
	return &q
}

func NewSelectRowsQuery[T any](table string, reader row.RowReader[T]) *query.SelectRowsQuery[T] {
	q := query.SelectRowsQuery[T]{}
	q.Initialize(table, reader)
	return &q
}

func NewFullSelectRowsQuery[T any](table string, reader row.RowReader[T]) *query.SelectRowsQuery[T] {
	var t T
	q := query.SelectRowsQuery[T]{}
	q.Initialize(table, reader)
	q.Columns(memo.ColumnsOf(t))
	return &q
}
