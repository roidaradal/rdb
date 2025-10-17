package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/row"
)

type TopRowQuery[T any] struct {
	conditionQuery
	columns []string
	order   string
	reader  row.RowReader[T]
}

// T = object type, V = value type
type TopValueQuery[T any, V any] struct {
	conditionQuery
	typeName string
	column   string
	order    string
	reader   row.RowReader[T]
}

// Initialize TopRowQuery
func (q *TopRowQuery[T]) Initialize(table string, reader row.RowReader[T]) {
	var t T
	q.conditionQuery.Initialize(table)
	q.columns = memo.ColumnsOf(t)
	q.reader = reader
}

// Initialize TopValueQuery
func (q *TopValueQuery[T, V]) Initialize(table string, fieldRef *V) {
	var t T
	q.conditionQuery.Initialize(table)
	q.typeName = dyn.TypeOf(t)
	q.column = memo.GetColumn(fieldRef)
	q.reader = row.Reader[T](q.column)
}

// Set TopRowQuery order (Ascending)
func (q *TopRowQuery[T]) OrderAsc(column string) {
	q.order = fmt.Sprintf("%s ASC", column)
}

// Set TopRowQuery order (Descending)
func (q *TopRowQuery[T]) OrderDesc(column string) {
	q.order = fmt.Sprintf("%s DESC", column)
}

// Set TopValueQuery order (Ascending)
func (q *TopValueQuery[T, V]) OrderAsc(column string) {
	q.order = fmt.Sprintf("%s ASC", column)
}

// Set TopValueQuery order (Descending)
func (q *TopValueQuery[T, V]) OrderDesc(column string) {
	q.order = fmt.Sprintf("%s DESC", column)
}

// Build the TopRowQuery
func (q TopRowQuery[T]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || len(q.columns) == 0 || q.order == "" {
		return emptyQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
	query := "SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT 1"
	query = fmt.Sprintf(query, columns, q.table, condition, q.order)
	return query, values
}

// Build the TopValueQuery
func (q TopValueQuery[T, V]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.column == "" || q.order == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT 1"
	query = fmt.Sprintf(query, q.column, q.table, condition, q.order)
	return query, values
}

// Execute the TopRowQuery and  get the top object
func (q TopRowQuery[T]) QueryRow(dbc *sql.DB) (*T, error) {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}
	row := dbc.QueryRow(query, values...)
	return q.reader(row)
}

// Execute the TopValueQuery and get the top value
func (q TopValueQuery[T, V]) QueryValue(dbc *sql.DB) (V, error) {
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
	return getColumnValue[V](item, q.typeName, q.column)
}
