package query

import (
	"database/sql"
	"fmt"
)

type CountQuery struct {
	conditionQuery
}

func (q CountQuery) Build() (string, []any) {
	if q.table == "" {
		return defaultQueryValues()
	}
	condition, values := q.condition.Build()
	query := "SELECT COUNT(*) FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)
	return query, values
}

func (q CountQuery) Count(dbc *sql.DB) (int, error) {
	if dbc == nil {
		return 0, errNoDBConnection
	}
	query, values := q.Build()
	if query == "" {
		return 0, errEmptyQuery
	}
	count := 0
	err := dbc.QueryRow(query, values...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q CountQuery) Exists(dbc *sql.DB) (bool, error) {
	count, err := q.Count(dbc)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
