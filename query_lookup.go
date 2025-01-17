package rdb

import (
	"database/sql"
	"fmt"
)

type lookupQuery[T any, K comparable, V any] struct {
	basicQuery[T]
	condition   Condition
	keyColumn   string
	valueColumn string
	reader      rowReader[T]
	keys        []K
}

/*
Output: Query (string), Values ([]any)

Note: Query could be blank string if invalid query parts
*/
func (q *lookupQuery[T, K, V]) Build() (string, []any) {
	// Check if table is blank
	if q.table == "" {
		return defaultQueryValues() // return empty query if blank table
	}

	// Check if blank column
	if q.keyColumn == "" || q.valueColumn == "" {
		return defaultQueryValues() // return empty query if blank column
	}

	// Build condition
	condition, values := q.condition.Build(q.object)

	// Build query
	query := "SELECT %s, %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.keyColumn, q.valueColumn, q.table, condition)

	return query, values
}

/*
Input: initialized DB connection

Output: map[key]value, where key in keys
*/
func (q *lookupQuery[T, K, V]) Lookup(dbc *sql.DB) (map[K]V, error) {
	if dbc == nil {
		return nil, errNoDBConnection
	}
	if q.reader == nil {
		return nil, errNoRowReader
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

	lookup := make(map[K]V)
	for rows.Next() {
		item, err := q.reader(rows)
		if err != nil {
			continue
		}
		key, err := getColumnValue[K](item, q.keyColumn)
		if err != nil {
			continue
		}
		value, err := getColumnValue[V](item, q.valueColumn)
		if err != nil {
			continue
		}
		lookup[key] = value
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lookup, nil
}

/*
Input: &struct &struct.KeyField, &struct.ValueField, []values, table

Output: &lookupQuery
*/
func NewLookupQuery[T any, K comparable, V any](object *T, key *K, value *V, keys []K, table string) *lookupQuery[T, K, V] {
	q := lookupQuery[T, K, V]{}
	q.initialize(object, table)
	columns := Columns(object, key, value)
	q.reader = Reader[T](columns)
	q.keyColumn = columns[0]
	q.valueColumn = columns[1]
	q.keys = keys
	q.condition = In(key, keys)
	return &q
}
