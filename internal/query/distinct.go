package query

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/row"
)

// T = object type, V = value type
type DistinctValuesQuery[T any, V any] struct {
	optionalConditionQuery
	column string
	reader row.RowReader[T]
}

// Initialize DistinctValuesQuery
func (q *DistinctValuesQuery[T, V]) Initialize(table string, fieldRef *V) {
	q.optionalConditionQuery.Initialize(table)
	q.column = memo.GetColumn(fieldRef)
	q.reader = row.Reader[T](q.column)
}

// Build the DistinctValuesQuery
func (q DistinctValuesQuery[T, V]) Build() (string, []any) {
	condition, values, err := q.optionalConditionQuery.preBuildCheck()
	if err != nil || q.column == "" {
		return emptyQueryValues()
	}
	query := "SELECT DISTINCT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.column, q.table, condition)
	return query, values
}

// Execute the DistinctValuesQuery and get list of distinct values
func (q DistinctValuesQuery[T, V]) Query(dbc *sql.DB) ([]V, error) {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}
	rows, err := dbc.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	distinct := make([]V, 0)
	typeName := ""
	for rows.Next() {
		item, err := q.reader(rows)
		if err != nil {
			continue
		}
		if typeName == "" {
			typeName = dyn.TypeOf(item)
		}
		value, err := getColumnValue[V](item, typeName, q.column)
		if err != nil {
			continue
		}
		distinct = append(distinct, value)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return distinct, nil
}
