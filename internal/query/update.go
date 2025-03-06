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

func Update[T any](q *UpdateQuery, key *T, value T) {
	update := kv.KeyValue(key, value)
	q.updates = append(q.updates, update)
}

func (q *UpdateQuery) Initialize(table string) {
	q.conditionQuery.Initialize(table)
	q.updates = make([]*kv.Value, 0)
}

func (q UpdateQuery) Build() (string, []any) {
	updateCount := len(q.updates)
	if q.table == "" || updateCount == 0 {
		return defaultQueryValues()
	}
	condition, conditionValues := q.condition.Build()
	values := make([]any, 0, updateCount+len(conditionValues))
	updates := make([]string, updateCount)
	for i, pair := range q.updates {
		column, value := pair.Get()
		if column == "" {
			return defaultQueryValues()
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
