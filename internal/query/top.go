package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/row"
	"github.com/roidaradal/rdb/internal/types"
)

type TopRowQuery[T any] struct {
	conditionQuery
	reader  row.RowReader[T]
	columns []string
	order   string
}

type TopValueQuery[T any, V any] struct {
	conditionQuery
	reader row.RowReader[T]
	column string
	order  string
}

func (q *TopRowQuery[T]) Initialize(table string, reader row.RowReader[T]) {
	var t T
	q.conditionQuery.Initialize(table)
	q.columns = memo.ColumnsOf(t)
	q.reader = reader
}

func (q *TopValueQuery[T, V]) Initialize(table string, field *V) {
	q.conditionQuery.Initialize(table)
	q.column = memo.GetColumn(field)
	q.reader = row.Reader[T]([]string{q.column})
}

func (q *TopRowQuery[T]) OrderAsc(column string) {
	q.order = fmt.Sprintf("%s ASC", column)
}

func (q *TopRowQuery[T]) OrderDesc(column string) {
	q.order = fmt.Sprintf("%s DESC", column)
}

func (q *TopValueQuery[T, V]) OrderAsc(column string) {
	q.order = fmt.Sprintf("%s ASC", column)
}

func (q *TopValueQuery[T, V]) OrderDesc(column string) {
	q.order = fmt.Sprintf("%s DESC", column)
}

func (q TopRowQuery[T]) Build() (string, []any) {
	if q.table == "" || len(q.columns) == 0 {
		return defaultQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
	condition, values := q.condition.Build()
	query := "SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT 1"
	query = fmt.Sprintf(query, columns, q.table, condition, q.order)
	return query, values
}

func (q TopValueQuery[T, V]) Build() (string, []any) {
	if q.table == "" || q.column == "" {
		return defaultQueryValues()
	}
	condition, values := q.condition.Build()
	query := "SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT 1"
	query = fmt.Sprintf(query, q.column, q.table, condition, q.order)
	return query, values
}

func (q TopRowQuery[T]) QueryRow(dbc *sql.DB) (*T, error) {
	if dbc == nil {
		return nil, errNoDBConnection
	}
	if q.reader == nil {
		return nil, errNoReader
	}
	query, values := q.Build()
	if query == "" {
		return nil, errEmptyQuery
	}
	row := dbc.QueryRow(query, values...)
	return q.reader(row)
}

func (q TopValueQuery[T, V]) QueryValue(dbc *sql.DB) (V, error) {
	var v V
	if dbc == nil {
		return v, errNoDBConnection
	}
	if q.reader == nil {
		return v, errNoReader
	}
	query, values := q.Build()
	if query == "" {
		return v, errEmptyQuery
	}
	row := dbc.QueryRow(query, values...)
	item, err := q.reader(row)
	if err != nil {
		return v, err
	}
	var t T
	schema := types.NameOf(t)
	return getColumnValue[V](item, schema, q.column)
}
