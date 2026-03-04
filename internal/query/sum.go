package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/roidaradal/fn/list"
	"github.com/roidaradal/rdb/internal/rdb"
)

// Sum Query
type SumQuery[T any] struct {
	conditionQuery
	columns []string
	reader  rdb.RowReader[T]
}

// Create new SumQuery
func NewSum[T any](table string, reader rdb.RowReader[T]) *SumQuery[T] {
	q := &SumQuery[T]{}
	q.initializeOptional(table)
	q.reader = reader
	q.columns = make([]string, 0)
	return q
}

// Set Sum columns
func (q *SumQuery[T]) Columns(columns []string) {
	q.columns = columns
}

// Build Sum Query
func (q SumQuery[T]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || len(q.columns) == 0 {
		return emptyQueryValues()
	}
	sumColumns := list.Map(q.columns, func(column string) string {
		return fmt.Sprintf("SUM(%s)", column)
	})
	columns := strings.Join(sumColumns, ", ")
	query := "SELECT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, columns, q.table, condition)
	return query, values
}

// Execute Sum Query and get sum object
func (q SumQuery[T]) Sum(dbc *sql.DB) (*T, error) {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}
	row := dbc.QueryRow(query, values...)
	return q.reader(row)
}
