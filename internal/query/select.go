package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/roidaradal/rdb/internal/row"
)

type SelectRowQuery[T any] struct {
	conditionQuery
	reader  row.RowReader[T]
	columns []string
}

type SelectRowsQuery[T any] struct {
	conditionQuery
	reader  row.RowReader[T]
	columns []string
	limit   uint
	offset  uint
	order   string
}

type SelectAllRowsQuery[T any] struct {
	baseQuery
	reader  row.RowReader[T]
	columns []string
}

func (q *SelectRowQuery[T]) Initialize(table string, reader row.RowReader[T]) {
	q.conditionQuery.Initialize(table)
	q.reader = reader
	q.columns = make([]string, 0)
}

func (q *SelectRowsQuery[T]) Initialize(table string, reader row.RowReader[T]) {
	q.conditionQuery.Initialize(table)
	q.reader = reader
	q.columns = make([]string, 0)
}

func (q *SelectAllRowsQuery[T]) Initialize(table string, reader row.RowReader[T]) {
	q.baseQuery.Initialize(table)
	q.reader = reader
	q.columns = make([]string, 0)
}

func (q *SelectRowQuery[T]) Columns(columns []string) {
	q.columns = columns
}

func (q *SelectRowsQuery[T]) Columns(columns []string) {
	q.columns = columns
}

func (q *SelectAllRowsQuery[T]) Columns(columns []string) {
	q.columns = columns
}

func (q *SelectRowsQuery[T]) Limit(limit uint) {
	q.offset = 0
	q.limit = limit
}

func (q *SelectRowsQuery[T]) Page(number, batchSize uint) {
	q.offset = (number - 1) * batchSize
	q.limit = batchSize
}

func (q *SelectRowsQuery[T]) OrderAsc(column string) {
	q.order = fmt.Sprintf("%s ASC", column)
}

func (q *SelectRowsQuery[T]) OrderDesc(column string) {
	q.order = fmt.Sprintf("%s DESC", column)
}

func (q SelectRowQuery[T]) Build() (string, []any) {
	if q.table == "" || len(q.columns) == 0 {
		return defaultQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
	condition, values := q.condition.Build()
	query := "SELECT %s FROM %s WHERE %s LIMIT 1"
	query = fmt.Sprintf(query, columns, q.table, condition)
	return query, values
}

func (q SelectRowsQuery[T]) Build() (string, []any) {
	if q.table == "" || len(q.columns) == 0 {
		return defaultQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
	condition, values := q.condition.Build()
	query := "SELECT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, columns, q.table, condition)
	if q.order != "" {
		query = fmt.Sprintf("%s ORDER BY %s", query, q.order)
	}
	if q.limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d, %d", query, q.offset, q.limit)
	}
	return query, values
}

func (q SelectAllRowsQuery[T]) Build() (string, []any) {
	if q.table == "" || len(q.columns) == 0 {
		return defaultQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
	query := "SELECT %s FROM %s"
	query = fmt.Sprintf(query, columns, q.table)
	return query, []any{}
}

func (q SelectRowQuery[T]) QueryRow(dbc *sql.DB) (*T, error) {
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

func (q SelectRowsQuery[T]) Query(dbc *sql.DB) ([]T, error) {
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
	rows, err := dbc.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]T, 0)
	for rows.Next() {
		item, err := q.reader(rows)
		if err != nil {
			continue
		}
		items = append(items, *item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (q SelectAllRowsQuery[T]) Query(dbc *sql.DB) ([]T, error) {
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
	rows, err := dbc.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]T, 0)
	for rows.Next() {
		item, err := q.reader(rows)
		if err != nil {
			continue
		}
		items = append(items, *item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
