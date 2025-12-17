package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/rdb"
)

// TopRow Query
type TopRow[T any] struct {
	conditionQuery
	columns []string
	limit   uint
	order   string
	reader  rdb.RowReader[T]
}

// TopValue Query
type TopValue[T, V any] struct {
	conditionQuery
	typeName   string
	columnName string
	limit      uint
	order      string
	reader     rdb.RowReader[T]
}

// Create new TopRow Query
func NewTopRow[T any](table string, reader rdb.RowReader[T]) *TopRow[T] {
	var t T
	q := &TopRow[T]{}
	q.initializeRequired(table)
	q.limit = 1
	q.columns = rdb.ColumnsOf(t)
	q.reader = reader
	return q
}

// Create new TopValue Query
func NewTopValue[T, V any](table string, fieldRef *V) *TopValue[T, V] {
	var t T
	q := &TopValue[T, V]{}
	q.initializeRequired(table)
	q.limit = 1
	q.typeName = dyn.TypeOf(t)
	q.columnName = rdb.GetColumnName(fieldRef)
	q.reader = rdb.NewReader[T](q.columnName)
	return q
}

// Set TopRow limit
func (q *TopRow[T]) Limit(limit uint) {
	q.limit = max(1, limit)
}

// Set TopValue limit
func (q *TopValue[T, V]) Limit(limit uint) {
	q.limit = max(1, limit)
}

// Set TopRow order (ascending)
func (q *TopRow[T]) OrderAsc(column string) {
	q.order = fmt.Sprintf("%s ASC", column)
}

// Set TopRow order (descending)
func (q *TopRow[T]) OrderDesc(column string) {
	q.order = fmt.Sprintf("%s DESC", column)
}

// Set TopValue order (ascending)
func (q *TopValue[T, V]) OrderAsc(column string) {
	q.order = fmt.Sprintf("%s ASC", column)
}

// Set TopValue order (descending)
func (q *TopValue[T, V]) OrderDesc(column string) {
	q.order = fmt.Sprintf("%s DESC", column)
}

// Build TopRow Query
func (q TopRow[T]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || len(q.columns) == 0 || q.order == "" {
		return emptyQueryValues()
	}
	columns := strings.Join(q.columns, ", ")
	query := "SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT %d"
	query = fmt.Sprintf(query, columns, q.table, condition, q.order, q.limit)
	return query, values
}

// Build TopValue Query
func (q TopValue[T, V]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.columnName == "" || q.order == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT %d"
	query = fmt.Sprintf(query, q.columnName, q.table, condition, q.order, q.limit)
	return query, values
}

// Execute TopRow Query and get top row object
func (q TopRow[T]) QueryRow(dbc *sql.DB) (*T, error) {
	q.limit = 1 // override limit = 1
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}
	row := dbc.QueryRow(query, values...)
	return q.reader(row)
}

// Execute TopRow Query and get top N row objects
func (q TopRow[T]) QueryRows(dbc *sql.DB) ([]*T, error) {
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

// Execute TopValue Query and get top value
func (q TopValue[T, V]) QueryValue(dbc *sql.DB) (V, error) {
	var v V
	q.limit = 1 // override limit = 1
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

// Execute TopValue Query and get top values
func (q TopValue[T, V]) QueryValues(dbc *sql.DB) ([]V, error) {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}

	topValues := make([]V, 0)
	err = readRows(dbc, query, values, q.reader, func(item *T) {
		value, err := getTypedColumnValue[V](item, q.typeName, q.columnName)
		if err != nil {
			return
		}
		topValues = append(topValues, value)
	})
	if err != nil {
		return nil, err
	}

	return topValues, nil
}
