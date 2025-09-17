package condition

import (
	"fmt"
	"strings"

	"github.com/roidaradal/rdb/internal/kv"
)

type Condition interface {
	// Condition: string, Values []any
	Build() (string, []any)
}

type None struct{}

func (c None) Build() (string, []any) {
	return defaultConditionValues()
}

type MatchAll struct{}

func (c MatchAll) Build() (string, []any) {
	return matchAllConditionValues()
}

type Value struct {
	pair *kv.Value
	op   string
}

func (c Value) Build() (string, []any) {
	column, value := c.pair.Get()
	if column == "" {
		return defaultConditionValues()
	} else {
		return singleConditionValues(column, c.op, value)
	}
}

func NewValue[T any](key *T, value T, op string) *Value {
	return &Value{
		pair: kv.KeyValue(key, value),
		op:   op,
	}
}

type List struct {
	pair   *kv.List
	listOp string
	soloOp string
}

func (c List) Build() (string, []any) {
	column, values := c.pair.Get()
	numValues := len(values)
	if column == "" || numValues == 0 {
		return defaultConditionValues()
	} else if numValues == 1 {
		return singleConditionValues(column, c.soloOp, values[0])
	} else {
		return listCondition(column, c.listOp, numValues), values
	}
}

func NewList[T any](key *T, values []T, listOp, soloOp string) *List {
	return &List{
		pair:   kv.KeyList(key, values),
		listOp: listOp,
		soloOp: soloOp,
	}
}

type Multi struct {
	conditions []Condition
	op         string
}

func (c Multi) Build() (string, []any) {
	numConditions := len(c.conditions)
	switch numConditions {
	case 0:
		return defaultConditionValues()
	case 1:
		return c.conditions[0].Build()
	default:
		conditions := make([]string, numConditions)
		allValues := make([]any, 0)
		for i, cond := range c.conditions {
			condition, values := cond.Build()
			if condition == defaultCondition {
				return defaultConditionValues()
			}
			conditions[i] = condition
			allValues = append(allValues, values...)
		}
		glue := fmt.Sprintf(" %s ", c.op)
		fullCondition := fmt.Sprintf("(%s)", strings.Join(conditions, glue)) // wrap in parentheses
		return fullCondition, allValues
	}
}

func NewMulti(op string, conditions ...Condition) *Multi {
	return &Multi{
		conditions: conditions,
		op:         op,
	}
}
