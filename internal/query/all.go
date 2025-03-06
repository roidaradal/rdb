package query

import (
	"fmt"

	"github.com/roidaradal/rdb/internal/condition"
)

type Query interface {
	// Query string, Values []any
	Build() (string, []any)
}

func defaultQueryValues() (string, []any) {
	return "", []any{}
}

type baseQuery struct {
	table string
}

func (q *baseQuery) Initialize(table string) {
	q.table = fmt.Sprintf("`%s`", table) // wrap in backticks
}

type conditionQuery struct {
	baseQuery
	condition condition.Condition
}

func (q *conditionQuery) Initialize(table string) {
	q.baseQuery.Initialize(table)
	q.condition = condition.None{}
}

func (q *conditionQuery) Where(queryCondition condition.Condition) {
	q.condition = queryCondition
}
