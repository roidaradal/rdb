package rdb

import (
	"fmt"
	"slices"
	"strings"
)

type insertRowQuery[T any] struct {
	basicQuery[T]
	row map[string]any
}

/*
Output: Query (string), Values ([]any)

Note: Query could be blank string if invalid query parts
*/
func (q *insertRowQuery[T]) Build() (string, []any) {
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
	placeholders := strings.Join(slices.Repeat([]string{"?"}, count), ", ")
	query := "INSERT INTO %s (%s) VALUES (%s)"
	query = fmt.Sprintf(query, q.table, cols, placeholders)

	return query, values
}

/*
Input: row map[string]any
*/
func (q *insertRowQuery[T]) Row(row map[string]any) {
	q.row = row
}

/*
Input: &struct, table (string)

Note: Same &struct will be used for setting row later

Output: &insertRowQuery
*/
func NewInsertRowQuery[T any](object *T, table string) *insertRowQuery[T] {
	q := insertRowQuery[T]{}
	q.initialize(object, table)
	q.row = make(map[string]any)
	return &q
}
