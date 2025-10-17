package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/roidaradal/rdb/internal/row"
)

type SelectRowQuery[T any] struct {
	conditionQuery
	columns []string
	reader  row.RowReader[T]
}

type SelectRowsQuery[T any] struct {
	optionalConditionQuery
	columns []string
	reader  row.RowReader[T]
	limit   uint
	offset  uint
	order   string
}

// Initialize SelectRowQuery
func (q *SelectRowQuery[T]) Initialize(table string, reader row.RowReader[T]) {
	q.conditionQuery.Initialize(table)
	q.reader = reader
	q.columns = make([]string, 0)
}

// Initialize SelectRowsQuery
func (q *SelectRowsQuery[T]) Initialize(table string, reader row.RowReader[T]) {
	q.optionalConditionQuery.Initialize(table)
	q.reader = reader
	q.columns = make([]string, 0)
}

// Set SelectRowQuery Columns
func (q *SelectRowQuery[T]) Columns(columns []string) {
	q.columns = columns
}

// Set SelectRowsQuery Columns
func (q *SelectRowsQuery[T]) Columns(columns []string) {
	q.columns = columns
}

// Set SelectRowsQuery Limit
func (q *SelectRowsQuery[T]) Limit(limit uint) {
	q.offset = 0
	q.limit = limit
}

// Set SelectRowsQuery Page number
func (q *SelectRowsQuery[T]) Page(number, batchSize uint) {
	q.offset = (number - 1) * batchSize
	q.limit = batchSize
}

// Set SelectRowsQuery column order (Ascending)
func (q *SelectRowsQuery[T]) OrderAsc(column string) {
	q.order = fmt.Sprintf("%s ASC", column)
}

// Set SelectRowsQuery column order (Descending)
func (q *SelectRowsQuery[T]) OrderDesc(column string) {
	q.order = fmt.Sprintf("%s DESC", column)
}

// Build the SelectRowQuery
func (q SelectRowQuery[T]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || len(q.columns) == 0 {
		return emptyQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
	query := "SELECT %s FROM %s WHERE %s LIMIT 1"
	query = fmt.Sprintf(query, columns, q.table, condition)
	return query, values
}

// Build the SelectRowsQuery
func (q SelectRowsQuery[T]) Build() (string, []any) {
	condition, values, err := q.optionalConditionQuery.preBuildCheck()
	if err != nil || len(q.columns) == 0 {
		return emptyQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
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

// Execute the SelectRowQuery and get the object
func (q SelectRowQuery[T]) QueryRow(dbc *sql.DB) (*T, error) {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}
	row := dbc.QueryRow(query, values...)
	return q.reader(row)
}

// Execute the SelectRowsQuery and get the list of objects
func (q SelectRowsQuery[T]) Query(dbc *sql.DB) ([]*T, error) {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}
	rows, err := dbc.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*T, 0)
	for rows.Next() {
		item, err := q.reader(rows)
		if err != nil {
			continue
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
