package query

import (
	"database/sql"
	"fmt"
)

type CountQuery struct {
	conditionQuery
}

// Build the CountQuery
func (q CountQuery) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil {
		return emptyQueryValues()
	}
	query := "SELECT COUNT(*) FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)
	return query, values
}

// Execute the CountQuery and get the count
func (q CountQuery) Count(dbc *sql.DB) (int, error) {
	query, values, err := preQueryCheck(q, dbc)
	if err != nil {
		return 0, err
	}
	count := 0
	err = dbc.QueryRow(query, values...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Execute the CountQuery and check if count > 0
func (q CountQuery) Exists(dbc *sql.DB) (bool, error) {
	count, err := q.Count(dbc)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
