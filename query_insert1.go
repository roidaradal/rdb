package rdb

import (
	"database/sql"
	"fmt"
	"strings"
)

type insertRowQuery struct {
	table string
	row   map[string]any
}

/*
Output: Query (string), Values ([]any)

Note: Query could be blank string if invalid query parts
*/
func (q *insertRowQuery) Build() (string, []any) {
	// Check if table is blank
	if q.table == "" {
		return defaultQueryValues() // return empty query if blank table
	}

	// Check if empty row
	count := len(q.row)
	if count == 0 {
		return defaultQueryValues() // return empty query if empty row
	}

	// Build columns
	columns := make([]string, 0, count)
	values := make([]any, 0, count)
	for column, value := range q.row {
		columns = append(columns, column)
		values = append(values, value)
	}

	// Build query
	cols := strings.Join(columns, ", ")
	placeholders := repeatString(count, "?", ", ")
	query := "INSERT INTO %s (%s) VALUES (%s)"
	query = fmt.Sprintf(query, q.table, cols, placeholders)

	return query, values
}

/*
Input: row map[string]any
*/
func (q *insertRowQuery) Row(row map[string]any) {
	q.row = row
}

/*
Input: initialized DB connection

Output: *sql.Result, error
*/
func (q *insertRowQuery) Exec(dbc *sql.DB) (*sql.Result, error) {
	return prepareAndExec(q, dbc)
}

/*
Input: table (string)

Output: &insertRowQuery
*/
func NewInsertRowQuery(table string) *insertRowQuery {
	q := insertRowQuery{}
	q.table = table
	q.row = make(map[string]any)
	return &q
}
