package query

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb/internal/rdb"
)

// Lookup Query, where T = object type, K = key type, V = value type
type Lookup[T any, K comparable, V any] struct {
	conditionQuery
	typeName    string
	keyColumn   string
	valueColumn string
	reader      rdb.RowReader[T]
}

// Create new Lookup Query
func NewLookup[T any, K comparable, V any](table string, keyFieldRef *K, valueFieldRef *V) *Lookup[T, K, V] {
	var t T
	q := &Lookup[T, K, V]{}
	q.initializeOptional(table)
	q.typeName = dyn.TypeOf(t)
	columns := rdb.GetColumnNames(keyFieldRef, valueFieldRef)
	if len(columns) == 2 {
		q.keyColumn = columns[0]
		q.valueColumn = columns[1]
		q.reader = rdb.NewReader[T](columns...)
	}
	return q
}

// Build Lookup Query
func (q Lookup[T, K, V]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.keyColumn == "" || q.valueColumn == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s, %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.keyColumn, q.valueColumn, q.table, condition)
	return query, values
}

// Execute Lookup Query and get map[K]V lookup
func (q Lookup[T, K, V]) Lookup(dbc *sql.DB) (map[K]V, error) {
	query, values, err := preReadCheck(q, dbc, q.reader)
	if err != nil {
		return nil, err
	}

	lookup := make(map[K]V)
	err = readRows(dbc, query, values, q.reader, func(item *T) {
		key, err1 := getTypedColumnValue[K](item, q.typeName, q.keyColumn)
		value, err2 := getTypedColumnValue[V](item, q.typeName, q.valueColumn)
		if err1 != nil || err2 != nil {
			return
		}
		lookup[key] = value
	})
	if err != nil {
		return nil, err
	}

	return lookup, nil
}
