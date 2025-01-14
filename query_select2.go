package rdb

import (
	"database/sql"
	"fmt"
	"strings"
)

type selectQuery[T any] struct {
	conditionQuery[T]
	columns []string
	reader  rowReader[T]
}

/*
Output: Query (string), Values ([]any)

Note: Query could be blank string if invalid query parts
*/
func (q *selectQuery[T]) Build() (string, []any) {
	// Check if table is blank
	if q.table == "" {
		return defaultQueryValues() // return empty query if blank table
	}

	// Check if empty columns
	if len(q.columns) == 0 {
		return defaultQueryValues() // return empty query if empty columns
	}

	// Build columns
	columns := strings.Join(q.columns, ", ")

	// Build condition
	condition, values := q.condition.Build(q.object)

	// Build query
	query := "SELECT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, columns, q.table, condition)

	return query, values
}

/*
Input: Columns []string

Note: Make sure corresponding Reader uses the same list of columns
*/
func (q *selectQuery[T]) Columns(columns []string) {
	q.columns = columns
}

/*
Input: reader function
*/
func (q *selectQuery[T]) SetReader(reader rowReader[T]) {
	q.reader = reader
}

/*
Input: initialized DB connection

Constraint: Need to call SetReader() first, otherwise nothing happens

Output: list of structs that contain reader data, error
*/
func (q *selectQuery[T]) Query(dbc *sql.DB) ([]T, error) {
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

	rows, err := dbc.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]T, 0)
	for rows.Next() {
		item, err := q.reader(rows)
		if err != nil {
			continue
		}
		items = append(items, *item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

/*
Input: &struct, table (string)

Note: Same &struct will be used for setting conditions later

Output: &SelectQuery
*/
func NewSelectQuery[T any](object *T, table string) *selectQuery[T] {
	q := selectQuery[T]{}
	q.initialize(object, table)
	q.columns = make([]string, 0)
	q.reader = nil
	return &q
}
