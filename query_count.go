package rdb

import (
	"database/sql"
	"fmt"
)

type countQuery[T any] struct {
	conditionQuery[T]
}

func (q *countQuery[T]) Build() (string, []any) {
	// Check if table is blank
	if q.table == "" {
		return defaultQueryValues()
	}
	// Build condition
	condition, values := q.condition.Build(q.object)

	// Build query
	query := "SELECT COUNT(*) FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)

	return query, values
}

func (q *countQuery[T]) Count(dbc *sql.DB) (int, error) {
	if dbc == nil {
		return 0, errNoDBConnection
	}
	query, values := q.Build()
	if query == "" {
		return 0, errEmptyQuery
	}
	count := 0
	err := dbc.QueryRow(query, values...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func NewCountQuery[T any](object *T, table string) *countQuery[T] {
	q := countQuery[T]{}
	q.initialize(object, table)
	return &q
}
