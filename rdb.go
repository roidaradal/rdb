package rdb

import (
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/query"
	"github.com/roidaradal/rdb/internal/row"
)

// Initialize rdb package
func Initialize() error {
	memo.Initialize()
	return nil
}

// Add a new type, extract its columns and fields, create its RowCreator
// StructRef is expected to be a struct pointer
func AddType(structRef any) error {
	err := memo.AddType(structRef)
	if err != nil {
		return err
	}
	typeName, allColumns := memo.TypeColumnsOf(structRef)
	memo.AddRowCreator(typeName, row.Creator(allColumns))
	return nil
}

// Get all column names of given item
func AllColumns(item any) []string {
	return memo.ColumnsOf(item)
}

// Get column name of given field pointer
func Column(fieldRef any) string {
	return memo.GetColumn(fieldRef)
}

// Get column names of given field pointers
func Columns(fieldRefs ...any) []string {
	return memo.GetColumns(fieldRefs...)
}

// Function that reads row values into object
type RowReader[T any] = row.RowReader[T]

// Creates a RowReader[T] with the given columns
func Reader[T any](columns ...string) RowReader[T] {
	return row.Reader[T](columns...)
}

// Creates a RowReader[T] using all columns
func FullReader[T any](structRef *T) RowReader[T] {
	return row.FullReader(structRef)
}

// Converts given object to map[string]any for row insertion
func ToRow[T any](structRef *T) dict.Object {
	typeName := dyn.TypeOf(structRef)
	rowFn, ok := memo.GetRowCreator(typeName)
	if !ok {
		return dict.Object{}
	}
	return rowFn(structRef)
}

// Function that checks SQL result if a condition has been satisfied
type QueryResultChecker = query.QueryResultChecker

// Gets the number of rows affected from SQL result (default: 0)
var RowsAffected = query.RowsAffected

// Gets the last insert ID (uint) from SQL result (default: 0)
var LastInsertID = query.LastInsertID

// QueryResultChecker that does nothing
var AssertNothing = query.AssertNothing

// Creates a QueryResultChecker that asserts the number of rows affected
var AssertRowsAffected = query.AssertRowsAffected

// Executes an SQL query
var Exec = query.Exec

// Executes an SQL query as part of a transaction,
// Applies Rollback on any errors
var ExecTx = query.ExecTx

// Rolls back the SQL transaction
var Rollback = query.Rollback
