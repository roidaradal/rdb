package query

import "fmt"

type DeleteQuery struct {
	conditionQuery
}

func (q DeleteQuery) Build() (string, []any) {
	if q.table == "" {
		return defaultQueryValues()
	}
	condition, values := q.condition.Build()
	query := "DELETE FROM %s WHERE %s"
	query = fmt.Sprintf(query, q.table, condition)
	return query, values
}
