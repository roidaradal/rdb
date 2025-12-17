package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/roidaradal/rdb/internal/rdb"
)

// SelectRow Query
type SelectRow[T any] struct {
	conditionQuery
	columns []string
	reader  rdb.RowReader[T]
}

// SelectRows Query
type SelectRows[T any] struct {
	conditionQuery
	columns []string
	reader  rdb.RowReader[T]
	limit   uint
	offset  uint
	order   string
}

// Create new SelectRow Query
func NewSelectRow[T any](table string, reader rdb.RowReader[T]) *SelectRow[T] {
	q := &SelectRow[T]{}
	q.initializeRequired(table)
	q.reader = reader
	q.columns = make([]string, 0)
	return q
}

// Create new SelectRow Query, using all columns
func NewFullSelectRow[T any](table string, reader rdb.RowReader[T]) *SelectRow[T] {
	var t T
	q := NewSelectRow(table, reader)
	q.Columns(rdb.ColumnsOf(t))
	return q
}

// Create new SelectRows Query
func NewSelectRows[T any](table string, reader rdb.RowReader[T]) *SelectRows[T] {
	q := &SelectRows[T]{}
	q.initializeOptional(table)
	q.reader = reader
	q.columns = make([]string, 0)
	return q
}

// Create new SelectRows Query, using all columns
func NewFullSelectRows[T any](table string, reader rdb.RowReader[T]) *SelectRows[T] {
	var t T
	q := NewSelectRows(table, reader)
	q.Columns(rdb.ColumnsOf(t))
	return q
}

// Set SelectRow columns
func (q *SelectRow[T]) Columns(columns []string) {
	q.columns = columns
}

// Set SelectRows columns
func (q *SelectRows[T]) Columns(columns []string) {
	q.columns = columns
}

// Set SelectRows Limit
func (q *SelectRows[T]) Limit(limit uint) {
	q.offset = 0
	q.limit = limit
}

// Set SelectRows Page number
func (q *SelectRows[T]) Page(number, batchSize uint) {
	q.offset = (number - 1) * batchSize
	q.limit = batchSize
}

// Set SelectRows column order (ascending)
func (q *SelectRows[T]) OrderAsc(column string) {
	q.order = fmt.Sprintf("%s ASC", column)
}

// Set SelectRows column order (descending)
func (q *SelectRows[T]) OrderDesc(column string) {
	q.order = fmt.Sprintf("%s DESC", column)
}

// Build SelectRow Query
func (q SelectRow[T]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || len(q.columns) == 0 {
		return emptyQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
	query := "SELECT %s FROM %s WHERE %s LIMIT 1"
	query = fmt.Sprintf(query, columns, q.table, condition)
	return query, values
}

// Build SelectRows Query
func (q SelectRows[T]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
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

// Execute SelectRow Query and get the row object
func (q SelectRow[T]) QueryRow(dbc *sql.DB) (*T, error) {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}
	row := dbc.QueryRow(query, values...)
	return q.reader(row)
}

// Execute SelectRows Query and get list of objects
func (q SelectRows[T]) Query(dbc *sql.DB) ([]*T, error) {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}

	items := make([]*T, 0)
	err = readRows(dbc, query, values, q.reader, func(item *T) {
		items = append(items, item)
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}
