package query

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/row"
)

// T = object type, V = value type
type ValueQuery[T any, V any] struct {
	conditionQuery
	column string
	reader row.RowReader[T]
}

// Initialize ValueQuery
func (q *ValueQuery[T, V]) Initialize(table string, fieldRef *V) {
	q.conditionQuery.Initialize(table)
	q.column = memo.GetColumn(fieldRef)
	q.reader = row.Reader[T](q.column)
}

// Build the ValueQuery
func (q ValueQuery[T, V]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.column == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.column, q.table, condition)
	return query, values
}

// Execute the ValueQuery and get the column value
func (q ValueQuery[T, V]) QueryValue(dbc *sql.DB) (V, error) {
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
	var t T
	typeName := dyn.TypeOf(t)
	return getColumnValue[V](item, typeName, q.column)
}

// Get the column value from given structRef and typeName
func getColumnValue[V any](structRef any, typeName, column string) (V, error) {
	var v V
	value, ok := row.GetColumnValue(structRef, typeName, column)
	if !ok {
		return v, errNotFoundField
	}
	v, ok = value.(V)
	if !ok {
		return v, errFailedTypeAssertion
	}
	return v, nil
}
