package query

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/rdb/internal/condition"
	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/row"
	"github.com/roidaradal/rdb/internal/types"
)

type DistinctValuesQuery[T any, V any] struct {
	conditionQuery
	reader row.RowReader[T]
	column string
}

func (q *DistinctValuesQuery[T, V]) Initialize(table string, field *V) {
	q.conditionQuery.Initialize(table)
	q.condition = condition.MatchAll{}
	q.column = memo.GetColumn(field)
	q.reader = row.Reader[T]([]string{q.column})
}

func (q DistinctValuesQuery[T, V]) Build() (string, []any) {
	if q.table == "" || q.column == "" {
		return defaultQueryValues()
	}
	condition, values := q.condition.Build()
	query := "SELECT DISTINCT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.column, q.table, condition)
	return query, values
}

func (q *DistinctValuesQuery[T, V]) Query(dbc *sql.DB) ([]V, error) {
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

	distinct := make([]V, 0)
	var schema string = ""
	for rows.Next() {
		item, err := q.reader(rows)
		if err != nil {
			continue
		}
		if schema == "" {
			schema = types.NameOf(item)
		}
		value, err := getColumnValue[V](item, schema, q.column)
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
