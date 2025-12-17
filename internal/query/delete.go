package query

import "fmt"

// Delete Query
type Delete struct {
	conditionQuery
}

// Create new Delete Query
func NewDelete(table string) *Delete {
	q := &Delete{}
	q.initializeRequired(table)
	return q
}

// Build Delete Query
func (q Delete) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil {
		return emptyQueryValues()
	}
	query := "DELETE FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)
	return query, values
}
