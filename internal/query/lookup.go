package query

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/rdb/internal/condition"
	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/row"
	"github.com/roidaradal/rdb/internal/types"
)

type LookupQuery[T any, K comparable, V any] struct {
	conditionQuery
	keyColumn   string
	valueColumn string
	reader      row.RowReader[T]
}

func (q *LookupQuery[T, K, V]) Initialize(table string, key *K, value *V) {
	q.conditionQuery.Initialize(table)
	q.condition = condition.MatchAll{}
	columns := memo.GetColumns(key, value)
	q.reader = row.Reader[T](columns)
	q.keyColumn = columns[0]
	q.valueColumn = columns[1]
}

func (q LookupQuery[T, K, V]) Build() (string, []any) {
	if q.table == "" || q.keyColumn == "" || q.valueColumn == "" {
		return defaultQueryValues()
	}
	condition, values := q.condition.Build()
	query := "SELECT %s, %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.keyColumn, q.valueColumn, q.table, condition)
	return query, values
}

func (q LookupQuery[T, K, V]) Lookup(dbc *sql.DB) (map[K]V, error) {
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

	lookup := make(map[K]V)
	var schema string = ""
	for rows.Next() {
		item, err := q.reader(rows)
		if err != nil {
			continue
		}
		if schema == "" {
			schema = types.NameOf(item)
		}
		key, err := getColumnValue[K](item, schema, q.keyColumn)
		if err != nil {
			continue
		}
		value, err := getColumnValue[V](item, schema, q.valueColumn)
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
