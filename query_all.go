package rdb

import (
	"slices"
	"strings"
)

type buildableQuery interface {
	// Output: Query, Values
	Build() (string, []any)
}

type basicQuery[T any] struct {
	object *T
	table  string
}

type conditionQuery[T any] struct {
	basicQuery[T]
	condition Condition
}

/******************************** QUERY METHODS ********************************/

func (q *basicQuery[T]) initialize(object *T, table string) {
	q.object = object
	q.table = table
}

/*************************** CONDITION QUERY METHODS ***************************/

func (q *conditionQuery[T]) initialize(object *T, table string) {
	q.basicQuery.initialize(object, table)
	q.condition = &noCondition{}
}

func (q *conditionQuery[T]) Where(condition Condition) {
	q.condition = condition
}

/*************************** PRIVATE FUNCTIONS ***************************/

func defaultQueryValues() (string, []any) {
	return "", []any{}
}

func repeatString(repeat int, item, glue string) string {
	return strings.Join(slices.Repeat([]string{item}, repeat), glue)
}
