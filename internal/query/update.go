package query

import (
	"fmt"
	"strings"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/rdb"
)

type FieldUpdate [2]any                    // [OldValue, NewValue]
type FieldUpdates = map[string]FieldUpdate // {FieldName => [OldValue, NewValue]}

// Return OldValue, NewValue
func (f FieldUpdate) Tuple() (any, any) {
	return f[0], f[1]
}

// Update Query
type Update[T any] struct {
	conditionQuery
	typeName string
	updates  []*rdb.Value
}

// Create new Update Query
func NewUpdate[T any](table string) *Update[T] {
	var t T
	q := &Update[T]{}
	q.initializeRequired(table)
	q.typeName = dyn.TypeOf(t)
	q.updates = make([]*rdb.Value, 0)
	return q
}

// Add field=value update to Update Query
func AddUpdate[T, V any](q *Update[T], fieldRef *V, value V) {
	// Note: Cannot be method as generics are not supported in methods
	q.updates = append(q.updates, rdb.KeyValue(fieldRef, value))
}

// Add column=value update to Update Query
func (q *Update[T]) Update(fieldName string, value any) {
	q.updates = append(q.updates, rdb.ColumnValue(q.typeName, fieldName, value))
}

// Add list of column=value updates to Update Query
func (q *Update[T]) Updates(updates FieldUpdates) {
	for fieldName, update := range updates {
		_, newValue := update.Tuple()
		q.Update(fieldName, newValue)
	}
}

// Build Update Query
func (q Update[T]) Build() (string, []any) {
	numUpdates := len(q.updates)
	condition, conditionValues, err := q.conditionQuery.preBuildCheck()
	if err != nil || numUpdates == 0 {
		return emptyQueryValues()
	}
	values := make([]any, 0, numUpdates+len(conditionValues))
	updates := make([]string, numUpdates)
	for i, pair := range q.updates {
		if pair == nil {
			// if key-value pair is nil, return empty query
			return emptyQueryValues()
		}
		column, value := pair.Tuple()
		if column == "" {
			// if blank column found, return empty query
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
