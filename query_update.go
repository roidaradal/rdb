package rdb

import (
	"fmt"
	"strings"
)

type updateQuery[T any] struct {
	conditionQuery[T]
	updates []*kvc
}

/*
Output: Query (string), Values ([]any)

Note: Query could be blank string if invalid query parts
*/
func (q *updateQuery[T]) Build() (string, []any) {
	// Check if table is blank
	if q.table == "" {
		return defaultQueryValues() // return empty query if blank table
	}

	// Check if empty updates
	updateCount := len(q.updates)
	if updateCount == 0 {
		return defaultQueryValues() // return empty query if nothing to update
	}

	// Build condition
	condition, conditionValues := q.condition.Build(q.object)
	conditionCount := len(conditionValues)

	// Build update
	values := make([]any, 0, updateCount+conditionCount)
	updates := make([]string, updateCount)
	for i, kv := range q.updates {
		column, value := kv.Build(q.object)
		if column == "" {
			return defaultQueryValues() // return empty query if blank column
		}
		updates[i] = fmt.Sprintf("%s = ?", column)
		values = append(values, value)
	}
	// Add condition values to end
	values = append(values, conditionValues...)

	// Build query
	update := strings.Join(updates, ", ")
	query := "UPDATE %s SET %s WHERE %s"
	query = fmt.Sprintf(query, q.table, update, condition)

	return query, values
}

/*
Input: &struct, table (string)

Note: Same &struct will be used for setting updates and conditions later

Output: &UpdateQuery
*/
func NewUpdateQuery[T any](object *T, table string) *updateQuery[T] {
	q := updateQuery[T]{}
	q.Initialize(object, table)
	q.updates = make([]*kvc, 0)
	return &q
}

/*
Input: &UpdateQuery, &struct.Field, value
*/
func Update[T any, U any](query *updateQuery[T], key *U, value U) {
	update := keyValue(key, value)
	query.updates = append(query.updates, update)
}
