package rdb

import (
	"database/sql"
	"fmt"
)

type selectRowQuery[T any] struct {
	selectQuery[T]
}

/*
Output: Query (string), Values ([]any)

Note: Query could be blank string if invalid query parts
*/
func (q *selectRowQuery[T]) Build() (string, []any) {
	query, values := q.selectQuery.Build()
	if query != "" {
		query = fmt.Sprintf("%s LIMIT 1", query)
	}
	return query, values
}

/*
Input: initialized DB connection

Output: &struct that contains reader data, error
*/
func (q *selectRowQuery[T]) QueryRow(dbc *sql.DB) (*T, error) {
	if dbc == nil {
		return nil, errNoDBConnection
	}
	if q.reader == nil {
		return nil, errNoRowReader
	}
	query, values := q.Build()
	if query == "" {
		return nil, errEmptyQuery
	}
	row := dbc.QueryRow(query, values...)
	return q.reader(row)
}

/*
Input: &struct, table (string)

Note: Same &struct will be used for setting conditions later

Output: &SelectRowQuery
*/
func NewSelectRowQuery[T any](object *T, table string, reader rowReader[T]) *selectRowQuery[T] {
	q := selectRowQuery[T]{
		selectQuery: *NewSelectQuery(object, table, reader),
	}
	return &q
}
