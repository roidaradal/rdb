package query

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/rdb"
)

// DistinctValues Query, where T = object type, V = value type
type DistinctValues[T, V any] struct {
	conditionQuery
	typeName   string
	columnName string
	reader     rdb.RowReader[T]
}

// Create new DistinctValues Query
func NewDistinctValues[T, V any](table string, fieldRef *V) *DistinctValues[T, V] {
	var t T
	q := &DistinctValues[T, V]{}
	q.initializeOptional(table)
	q.typeName = dyn.TypeOf(t)
	q.columnName = rdb.GetColumnName(fieldRef)
	q.reader = rdb.NewReader[T](q.columnName)
	return q
}

// Build DistinctValues Query
func (q DistinctValues[T, V]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.columnName == "" {
		return emptyQueryValues()
	}
	query := "SELECT DISTINCT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.columnName, q.table, condition)
	return query, values
}

// Execute DistinctValues Query and get list of distinct values
func (q DistinctValues[T, V]) Query(dbc *sql.DB) ([]V, error) {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}

	distinct := make([]V, 0)
	err = readRows(dbc, query, values, q.reader, func(item *T) {
		value, err := getTypedColumnValue[V](item, q.typeName, q.columnName)
		if err != nil {
			return
		}
		distinct = append(distinct, value)
	})
	if err != nil {
		return nil, err
	}

	return distinct, nil
}
