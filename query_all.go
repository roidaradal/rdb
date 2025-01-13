package rdb

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

func (q *basicQuery[T]) Initialize(object *T, table string) {
	q.object = object
	q.table = table
}

/*************************** CONDITION QUERY METHODS ***************************/

func (q *conditionQuery[T]) Initialize(object *T, table string) {
	q.basicQuery.Initialize(object, table)
	q.condition = &noCondition{}
}

func (q *conditionQuery[T]) Where(condition Condition) {
	q.condition = condition
}

/*************************** PRIVATE FUNCTIONS ***************************/

func defaultQueryValues() (string, []any) {
	return "", []any{}
}
