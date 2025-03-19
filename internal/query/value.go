package query

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/row"
	"github.com/roidaradal/rdb/internal/types"
)

type ValueQuery[T any, V any] struct {
	conditionQuery
	column string
	reader row.RowReader[T]
}

func (q *ValueQuery[T, V]) Initialize(table string, field *V) {
	q.conditionQuery.Initialize(table)
	q.column = memo.GetColumn(field)
	q.reader = row.Reader[T]([]string{q.column})
}

func (q ValueQuery[T, V]) Build() (string, []any) {
	if q.table == "" || q.column == "" {
		return defaultQueryValues()
	}
	condition, values := q.condition.Build()
	query := "SELECT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.column, q.table, condition)
	return query, values
}

func (q ValueQuery[T, V]) QueryValue(dbc *sql.DB) (V, error) {
	var v V
	if dbc == nil {
		return v, errNoDBConnection
	}
	if q.reader == nil {
		return v, errNoReader
	}
	query, values := q.Build()
	if query == "" {
		return v, errEmptyQuery
	}
	row := dbc.QueryRow(query, values...)
	item, err := q.reader(row)
	if err != nil {
		return v, err
	}
	var t T
	schema := types.NameOf(t)
	return getColumnValue[V](item, schema, q.column)
}

func getColumnValue[V any](x any, schema, column string) (V, error) {
	var v V
	value, ok := row.FindColumnValue(x, schema, column)
	if !ok {
		return v, errFieldNotFound
	}
	v, ok = value.(V)
	if !ok {
		return v, errTypeAssertion
	}
	return v, nil
}
