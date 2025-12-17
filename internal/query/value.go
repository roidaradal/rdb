package query

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/rdb"
)

// Value Query, where T = object type, V = value type
type Value[T, V any] struct {
	conditionQuery
	typeName   string
	columnName string
	reader     rdb.RowReader[T]
}

// Create new Value Query
func NewValue[T, V any](table string, fieldRef *V) *Value[T, V] {
	var t T
	q := &Value[T, V]{}
	q.initializeRequired(table)
	q.typeName = dyn.TypeOf(t)
	q.columnName = rdb.GetColumnName(fieldRef)
	q.reader = rdb.NewReader[T](q.columnName)
	return q
}

// Build Value Query
func (q Value[T, V]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.columnName == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.columnName, q.table, condition)
	return query, values
}

// Execute Value Query and get column value
func (q Value[T, V]) QueryValue(dbc *sql.DB) (V, error) {
	var v V
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return v, err
	}
	row := dbc.QueryRow(query, values...)
	item, err := q.reader(row)
	if err != nil {
		return v, err
	}
	return getTypedColumnValue[V](item, q.typeName, q.columnName)
}

// Get column value from given struct pointer and type name.
// Type coerce the value into the correct type
func getTypedColumnValue[V any](structRef any, typeName, columnName string) (V, error) {
	var v V
	value, ok := rdb.GetStructColumnValue(structRef, typeName, columnName)
	if !ok {
		return v, errNotFoundField
	}
	v, ok = value.(V)
	if !ok {
		return v, errFailedTypeAssertion
	}
	return v, nil
}
