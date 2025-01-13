package rdb

import "fmt"

type DeleteQuery[T any] struct {
	ConditionQuery[T]
}

/*
Output: Query (string), Values ([]any)
*/
func (q *DeleteQuery[T]) Build() (string, []any) {
	// Build condition
	condition, values := q.condition.Build(q.object)

	// Build query
	query := "DELETE FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)

	return query, values
}

/*
Input: &struct, table (string)

Note: Same &struct will be used for setting columns and conditions later

Output: &DeleteQuery, or nil if table is blank
*/
func NewDeleteQuery[T any](object *T, table string) *DeleteQuery[T] {
	if table == "" {
		return nil
	}
	q := DeleteQuery[T]{}
	q.Initialize(object, table)
	return &q
}
