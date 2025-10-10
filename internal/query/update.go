package query

import (
	"fmt"
	"strings"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/kv"
)

type FieldUpdates = map[string]FieldUpdate
type FieldUpdate [2]any // [Old, New]

// Returns Old, New
func (f FieldUpdate) Tuple() (any, any) {
	return f[0], f[1]
}

type UpdateQuery[T any] struct {
	conditionQuery
	typeName string
	updates  []*kv.Value
}

// Add query update (key = value),
// Cannot be a method as generics are not supported in methods
func Update[T any, V any](q *UpdateQuery[T], fieldRef *V, value V) {
	q.updates = append(q.updates, kv.KeyValue(fieldRef, value))
}

// Initialize UpdateQuery
func (q *UpdateQuery[T]) Initialize(table string) {
	var t T
	q.conditionQuery.Initialize(table)
	q.typeName = dyn.TypeOf(t)
	q.updates = make([]*kv.Value, 0)
}

// Add one column=value update
func (q *UpdateQuery[T]) Update(fieldName string, value any) {
	q.updates = append(q.updates, kv.ColumnValue(q.typeName, fieldName, value))
}

// Add list of column=value updates
func (q *UpdateQuery[T]) Updates(updates FieldUpdates) {
	for fieldName, update := range updates {
		_, value := update.Tuple()
		q.Update(fieldName, value)
	}
}

// Build UpdateQuery
func (q UpdateQuery[T]) Build() (string, []any) {
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
