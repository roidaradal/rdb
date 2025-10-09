package condition

import (
	"fmt"
	"strings"

	"github.com/roidaradal/rdb/internal/kv"
)

// Condition object only needs to implement the Build() method
type Condition interface {
	// condition string, values []any
	Build() (string, []any)
}

// No Condition; used as default for UPDATE, DELETE
type None struct{}

// Implement Condition interface for condition.None
func (c None) Build() (string, []any) {
	return falseConditionValues()
}

// Match All Condition; used as default for SELECT
type MatchAll struct{}

// Implement Condition interface for condition.MatchAll
func (c MatchAll) Build() (string, []any) {
	return trueConditionValues()
}

// Value Condition; uses kv.Value (one value)
type Value struct {
	pair *kv.Value
	op   string
}

// Implement Condition interface for condition.Value
func (c Value) Build() (string, []any) {
	if c.pair == nil {
		// no pair = false condition
		return falseConditionValues()
	}
	column, value := c.pair.Tuple()
	if column == "" {
		// no column = false condition
		return falseConditionValues()
	}
	return soloConditionValues(column, c.op, value)
}

// Creates a new condition.Value
func NewValue[T any](fieldRef *T, value T, op string) *Value {
	return &Value{kv.KeyValue(fieldRef, value), op}
}

// List Condition; uses kv.List (multiple values)
type List struct {
	pair   *kv.List
	listOp string
	soloOp string
}

// Implement Condition interface for condition.List
func (c List) Build() (string, []any) {
	if c.pair == nil {
		// no pair = false condition
		return falseConditionValues()
	}
	column, values := c.pair.Tuple()
	numValues := len(values)
	if column == "" || numValues == 0 {
		// no column or no values = false condition
		return falseConditionValues()
	} else if numValues == 1 {
		return soloConditionValues(column, c.soloOp, values[0])
	} else {
		return listCondition(column, c.listOp, numValues), values
	}
}

// Creates a new condition.List
func NewList[T any](fieldRef *T, values []T, listOp, soloOp string) *List {
	return &List{kv.KeyList(fieldRef, values), listOp, soloOp}
}

// Multi Condition; multiple conditions (AND, OR)
type Multi struct {
	conditions []Condition
	op         string
}

// Implement Condition interface for condition.Multi
func (c Multi) Build() (string, []any) {
	numConditions := len(c.conditions)
	switch numConditions {
	case 0:
		// no conditions = false condition
		return falseConditionValues()
	case 1:
		// one condition = only build that one
		return c.conditions[0].Build()
	default:
		conditions := make([]string, numConditions)
		allValues := make([]any, 0)
		for i, condition := range c.conditions {
			conditionString, values := condition.Build()
			if conditionString == falseCondition {
				// If any of the conditions failed, return false condition
				return falseConditionValues()
			}
			conditions[i] = conditionString
			allValues = append(allValues, values...)
		}
		glue := fmt.Sprintf(" %s ", c.op)
		fullCondition := fmt.Sprintf("(%s)", strings.Join(conditions, glue)) // join by operator and wrap in parentheses
		return fullCondition, allValues
	}
}

// Creates a new condition.Multi
func NewMulti(op string, conditions ...Condition) *Multi {
	return &Multi{conditions, op}
}
