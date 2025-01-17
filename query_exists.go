package rdb

import (
	"database/sql"
	"fmt"
)

type existsQuery[T any] struct {
	conditionQuery[T]
}

/*
Output: Query (string), Values ([]any)

Note: Query could be blank string if invalid query parts
*/
func (q *existsQuery[T]) Build() (string, []any) {
	// Check if table is blank
	if q.table == "" {
		return defaultQueryValues() // return empty query if blank table
	}

	// Build condition
	condition, values := q.condition.Build(q.object)

	// Build query
	query := "SELECT COUNT(*) FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)

	return query, values
}

/*
Input: initialized DB connection

Output: exists boolean, error
*/
func (q *existsQuery[T]) Exists(dbc *sql.DB) (bool, error) {
	if dbc == nil {
		return false, errNoDBConnection
	}
	query, values := q.Build()
	if query == "" {
		return false, errEmptyQuery
	}
	count := 0
	err := dbc.QueryRow(query, values...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

/*
Input: &struct, table (string)

Note: Same &struct will be used for setting conditions later

Output: &existsQuery
*/
func NewExistsQuery[T any](object *T, table string) *existsQuery[T] {
	q := existsQuery[T]{}
	q.initialize(object, table)
	return &q
}
