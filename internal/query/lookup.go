package query

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/row"
)

// T = object type, K = key type, V = value type
type LookupQuery[T any, K comparable, V any] struct {
	optionalConditionQuery
	keyColumn   string
	valueColumn string
	reader      row.RowReader[T]
}

// Initialize LookupQuery
func (q *LookupQuery[T, K, V]) Initialize(table string, keyFieldRef *K, valueFieldRef *V) {
	q.optionalConditionQuery.Initialize(table)
	columns := memo.GetColumns(keyFieldRef, valueFieldRef)
	if len(columns) == 2 {
		q.reader = row.Reader[T](columns...)
		q.keyColumn = columns[0]
		q.valueColumn = columns[1]
	}
}

// Build the LookupQuery
func (q LookupQuery[T, K, V]) Build() (string, []any) {
	condition, values, err := q.optionalConditionQuery.preBuildCheck()
	if err != nil || q.keyColumn == "" || q.valueColumn == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s, %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.keyColumn, q.valueColumn, q.table, condition)
	return query, values
}

// Execute the LookupQuery and get the map[K]V lookup
func (q LookupQuery[T, K, V]) Lookup(dbc *sql.DB) (map[K]V, error) {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}
	rows, err := dbc.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lookup := make(map[K]V)
	typeName := ""
	for rows.Next() {
		item, err := q.reader(rows)
		if err != nil {
			continue
		}
		if typeName == "" {
			typeName = dyn.TypeOf(item)
		}
		key, err1 := getColumnValue[K](item, typeName, q.keyColumn)
		value, err2 := getColumnValue[V](item, typeName, q.valueColumn)
		if err1 != nil || err2 != nil {
			continue
		}
		lookup[key] = value
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return lookup, nil
}
