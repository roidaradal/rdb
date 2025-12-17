package query

import (
	"database/sql"
	"fmt"
)

// Count Query
type Count struct {
	conditionQuery
}

// Create new Count Query
func NewCount(table string) *Count {
	q := &Count{}
	q.initializeRequired(table)
	return q
}

// Build Count Query
func (q Count) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil {
		return emptyQueryValues()
	}
	query := "SELECT COUNT(*) FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)
	return query, values
}

// Execute CountQuery and get count
func (q Count) Count(dbc *sql.DB) (int, error) {
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

// Execute CountQuery and check if count > 0
func (q Count) Exists(dbc *sql.DB) (bool, error) {
	count, err := q.Count(dbc)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
