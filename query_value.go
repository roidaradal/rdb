package rdb

import (
	"database/sql"
	"fmt"
)

type valueQuery[T any, V any] struct {
	conditionQuery[T]
	column string
	reader RowReader[T]
}

/*
Output: Query (string), Values ([]any)

Note: Query could be blank string if invalid query parts
*/
func (q *valueQuery[T, V]) Build() (string, []any) {
	// Check if table is blank
	if q.table == "" {
		return defaultQueryValues() // return empty query if blank table
	}

	// Check if blank column
	if q.column == "" {
		return defaultQueryValues() // return empty query if blank column
	}

	// Build condition
	condition, values := q.condition.Build(q.object)

	// Build query
	query := "SELECT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.column, q.table, condition)

	return query, values
}

/*
Input: initialized DB connection

Output: value, error
*/
func (q *valueQuery[T, V]) QueryValue(dbc *sql.DB) (V, error) {
	var v V
	if dbc == nil {
		return v, errNoDBConnection
	}
	if q.reader == nil {
		return v, errNoRowReader
	}
	query, values := q.Build()
	if query == "" {
		return v, errEmptyQuery
	}
	row := dbc.QueryRow(query, values...)
	item, err := q.reader(row)
	if err != nil {
		return v, err
	}
	return getColumnValue[V](item, q.column)
}

/*
Input: &struct, &struct.Field, table

Note: Same &struct will be used for setting conditions later

Output: &valueQuery
*/
func NewValueQuery[T any, V any](object *T, field *V, table string) *valueQuery[T, V] {
	q := valueQuery[T, V]{}
	q.initialize(object, table)
	q.column = Column(object, field)
	q.reader = Reader[T]([]string{q.column})
	return &q
}
