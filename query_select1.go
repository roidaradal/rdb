package rdb

import (
	"database/sql"
	"fmt"
)

type selectRowQuery[T any] struct {
	selectQuery[T]
	reader rowReader[T]
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
Input: reader function
*/
func (q *selectRowQuery[T]) SetReader(reader rowReader[T]) {
	q.reader = reader
}

/*
Input: initialized DB connection

Constraint: Need to call SetReader() first, otherwise nothing happens

Output: &struct that contains reader data, error
*/
func (q *selectRowQuery[T]) Run(dbc *sql.DB) (*T, error) {
	query, values := q.Build()
	if query == "" {
		return nil, errEmptyQuery
	}
	if dbc == nil {
		return nil, errNoDBConnection
	}
	if q.reader == nil {
		return nil, errNoRowReader
	}
	row := dbc.QueryRow(query, values...)
	return q.reader(row)
}

/*
Input: &struct, table (string)

Note: Same &struct will be used for setting conditions later

Output: &SelectRowQuery
*/
func NewSelectRowQuery[T any](object *T, table string) *selectRowQuery[T] {
	q := selectRowQuery[T]{}
	q.initialize(object, table)
	q.columns = make([]string, 0)
	q.reader = nil
	return &q
}
