package query

import (
	"fmt"
	"strings"

	"github.com/roidaradal/rdb/internal/condition"
)

type Query interface {
	// Query string, Values []any
	Build() (string, []any)
}

func defaultQueryValues() (string, []any) {
	return "", []any{}
}

func QueryString(q Query) string {
	query, values := q.Build()
	query = strings.Replace(query, "?", "%v", -1)
	return fmt.Sprintf(query, values...)
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
