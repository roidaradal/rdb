package rdb

type BuildableQuery interface {
	// Output: Query, Values
	Build() (string, []any)
}

type Query[T any] struct {
	object *T
	table  string
}

type ConditionQuery[T any] struct {
	Query[T]
	condition Condition
}

/******************************** QUERY METHODS ********************************/

func (q *Query[T]) Initialize(object *T, table string) {
	q.object = object
	q.table = table
}

/*************************** CONDITION QUERY METHODS ***************************/

func (q *ConditionQuery[T]) Initialize(object *T, table string) {
	q.Query.Initialize(object, table)
	q.condition = &noCondition{}
}

func (q *ConditionQuery[T]) Where(condition Condition) {
	q.condition = condition
}

/*************************** PRIVATE FUNCTIONS ***************************/

func defaultQueryValues() (string, []any) {
	return "", []any{}
}
