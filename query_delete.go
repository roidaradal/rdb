package rdb

import "fmt"

type deleteQuery[T any] struct {
	conditionQuery[T]
}

/*
Output: Query (string), Values ([]any)

Note: Query could be blank string if invalid query parts
*/
func (q *deleteQuery[T]) Build() (string, []any) {
	// Check if table is blank
	if q.table == "" {
		return defaultQueryValues() // return empty query if blank table
	}

	// Build condition
	condition, values := q.condition.Build(q.object)

	// Build query
	query := "DELETE FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)

	return query, values
}

/*
Input: &struct, table (string)

Note: Same &struct will be used for setting conditions later

Output: &DeleteQuery
*/
func NewDeleteQuery[T any](object *T, table string) *deleteQuery[T] {
	q := deleteQuery[T]{}
	q.Initialize(object, table)
	return &q
}
