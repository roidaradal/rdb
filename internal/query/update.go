package query

import (
	"fmt"
	"strings"

	"github.com/roidaradal/rdb/internal/kv"
)

type UpdateQuery struct {
	conditionQuery
	updates []*kv.Value
}

// Add query update (key = value),
// Cannot be a method as generics are not supported in methods
func Update[T any](q *UpdateQuery, fieldRef *T, value T) {
	q.updates = append(q.updates, kv.KeyValue(fieldRef, value))
}

// Initialize UpdateQuery
func (q *UpdateQuery) Initialize(table string) {
	q.conditionQuery.Initialize(table)
	q.updates = make([]*kv.Value, 0)
}

// Build UpdateQuery
func (q UpdateQuery) Build() (string, []any) {
	numUpdates := len(q.updates)
	condition, conditionValues, err := q.conditionQuery.preBuildCheck()
	if err != nil || numUpdates == 0 {
		return emptyQueryValues()
	}
	values := make([]any, 0, numUpdates+len(conditionValues))
	updates := make([]string, numUpdates)
	for i, pair := range q.updates {
		if pair == nil {
			// if kv pair is nil, return empty query
			return emptyQueryValues()
		}
		column, value := pair.Tuple()
		if column == "" {
			// if blank column is found, return empty query
			return emptyQueryValues()
		}
		updates[i] = fmt.Sprintf("%s = ?", column)
		values = append(values, value)
	}
	values = append(values, conditionValues...)
	update := strings.Join(updates, ", ")
	query := "UPDATE %s SET %s WHERE %s"
	query = fmt.Sprintf(query, q.table, update, condition)
	return query, values
}
