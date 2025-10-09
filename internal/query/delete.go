package query

import "fmt"

type DeleteQuery struct {
	conditionQuery
}

// Build DeleteQuery
func (q DeleteQuery) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil {
		return emptyQueryValues()
	}
	query := "DELETE FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)
	return query, values
}
