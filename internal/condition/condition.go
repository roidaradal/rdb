// Package condition contains types and functions to build query conditions
package condition

import (
	"fmt"
	"strings"

	"github.com/roidaradal/rdb/internal/rdb"
)

// Condition object needs to implement the Build() method
// to output the condition string and the parameter values
type Condition interface {
	Build() (string, []any) // Return (condition string, parameter values)
}

// Missing Condition; default for UPDATE, DELETE to ensure condition is set
// Equivalent to 'WHERE false'
type Missing struct{}

// MatchAll Condition; default for SELECT (no condition)
// Equivalent to 'WHERE true'
type MatchAll struct{}

// Value Condition, uses KeyValue (one value)
type Value struct {
	pair     *rdb.Value
	operator string
}

// List Condition, uses KeyList (multiple values)
type List struct {
	pair         *rdb.List
	listOperator string
	soloOperator string
}

// Multi Condition for joining multiple conditions through AND, OR
type Multi struct {
	conditions []Condition // Multiple conditions
	operator   string      // Join operator
}

// Build Missing condition
func (c Missing) Build() (string, []any) {
	return falseConditionValues()
}

// Build MatchAll condition
func (c MatchAll) Build() (string, []any) {
	return trueConditionValues()
}

// Build Value condition
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
	return soloConditionValues(column, c.operator, value)
}

// Build List condition
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
		return soloConditionValues(column, c.soloOperator, values[0])
	} else {
		return listCondition(column, c.listOperator, numValues), values
	}
}

// Build Multi condition
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
				// If any condition fails, return false condition immediately
				return falseConditionValues()
			}
			conditions[i] = conditionString
			allValues = append(allValues, values...)
		}
		// Join by operator and wrap in parentheses
		glue := fmt.Sprintf(" %s ", c.operator)
		fullCondition := fmt.Sprintf("(%s)", strings.Join(conditions, glue))
		return fullCondition, allValues
	}
}

// Create new Value condition
func NewValue[T any](fieldRef *T, value T, operator string) *Value {
	return &Value{rdb.KeyValue(fieldRef, value), operator}
}

// Create new List condition
func NewList[T any](fieldRef *T, values []T, listOperator, soloOperator string) *List {
	return &List{rdb.KeyList(fieldRef, values), listOperator, soloOperator}
}

// Create new Multi condition
func NewMulti(operator string, conditions ...Condition) *Multi {
	return &Multi{conditions, operator}
}
