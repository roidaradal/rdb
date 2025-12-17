package rdb

import "github.com/roidaradal/rdb/internal/query"

type (
	Query         = query.Query         // Query interface
	ResultChecker = query.ResultChecker // Checks SQL result if condition is satisfied
	FieldUpdate   = query.FieldUpdate   // [OldValue, NewValue]
	FieldUpdates  = query.FieldUpdates  // {FieldName => [OldValue, NewValue]}
)

var (
	AssertNothing      = query.AssertNothing      // ResultChecker that does nothing
	AssertRowsAffected = query.AssertRowsAffected // Creates ResultChecker that asserts number of rows affected
	RowsAffected       = query.RowsAffected       // Get number of rows affected from SQL result (defaults to 0)
	LastInsertID       = query.LastInsertID       // Get last insert ID from SQL result (defaults to 0)
)

var (
	QueryString        = query.ToString      // Build full query string
	NewCountQuery      = query.NewCount      // Create new Count query
	NewDeleteQuery     = query.NewDelete     // Create new Delete Query
	NewInsertRowQuery  = query.NewInsertRow  // Create new InsertRow Query
	NewInsertRowsQuery = query.NewInsertRows // Create new InsertRows Query
)

// Create new Update Query
func NewUpdateQuery[T any](table string) *query.Update[T] {
	return query.NewUpdate[T](table)
}

// Add field=value update to Update Query
func Update[T any, V any](q *query.Update[T], fieldRef *V, value V) {
	query.AddUpdate(q, fieldRef, value)
}

// Create new Value Query
func NewValueQuery[T, V any](table string, fieldRef *V) *query.Value[T, V] {
	return query.NewValue[T](table, fieldRef)
}

// Create new Lookup Query
func NewLookupQuery[T any, K comparable, V any](table string, keyFieldRef *K, valueFieldRef *V) *query.Lookup[T, K, V] {
	return query.NewLookup[T](table, keyFieldRef, valueFieldRef)
}

// Create new DistinctValues Query
func NewDistinctValuesQuery[T, V any](table string, fieldRef *V) *query.DistinctValues[T, V] {
	return query.NewDistinctValues[T](table, fieldRef)
}

// Create new SelectRow Query
func NewSelectRowQuery[T any](table string, reader RowReader[T]) *query.SelectRow[T] {
	return query.NewSelectRow(table, reader)
}

// Create new SelectRow Query, using all columns
func NewFullSelectRowQuery[T any](table string, reader RowReader[T]) *query.SelectRow[T] {
	return query.NewFullSelectRow(table, reader)
}

// Create new SelectRows Query
func NewSelectRowsQuery[T any](table string, reader RowReader[T]) *query.SelectRows[T] {
	return query.NewSelectRows(table, reader)
}

// Create new SelectRows Query, using all columns
func NewFullSelectRowsQuery[T any](table string, reader RowReader[T]) *query.SelectRows[T] {
	return query.NewFullSelectRows(table, reader)
}

// Create new TopRow Query
func NewTopRowQuery[T any](table string, reader RowReader[T]) *query.TopRow[T] {
	return query.NewTopRow(table, reader)
}

// Create new TopValue Query
func NewTopValueQuery[T, V any](table string, fieldRef *V) *query.TopValue[T, V] {
	return query.NewTopValue[T](table, fieldRef)
}
