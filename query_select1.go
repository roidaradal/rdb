package rdb

import "fmt"

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
Input: &struct, table (string)

Note: Same &struct will be used for setting conditions later

Output: &SelectRowQuery
*/
func NewSelectRowQuery[T any](object *T, table string) *selectRowQuery[T] {
	q := selectRowQuery[T]{}
	q.initialize(object, table)
	q.columns = make([]string, 0)
	return &q
}
