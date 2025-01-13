package rdb

import (
	"fmt"
	"slices"
	"strings"
)

const defaultCondition string = "false"

type condition interface {
	// Input: &struct, same struct used for Field=&struct.Field
	// Output: condition (string), values ([]any)
	Build(any) (string, []any)
}

type noCondition struct{}

type kvCondition struct {
	pair     *kvc
	operator string
}

type klCondition struct {
	pair         *klc
	listOperator string
	soloOperator string
}

type conditionSet struct {
	conditions []condition
	operator   string
}

/******************************** BUILD METHODS ********************************/

/*
Input: &struct, same struct used for Field=&struct.Field

Note: Ignores the struct input, just passed in to meet the Build() interface

Output: "false", empty list of values
*/
func (c *noCondition) Build(t any) (string, []any) {
	return defaultConditionValues()
}

/*
Input: &struct, same struct used for Field=&struct.Field

Output:  "<column> <operator> ?", []any{value}
*/
func (c *kvCondition) Build(t any) (string, []any) {
	column, value := c.pair.Build(t)
	if column == "" {
		return defaultConditionValues()
	} else {
		return valueCondition(column, c.operator), []any{value}
	}
}

/*
Input: &struct, same struct used for Field=&struct.Field

Output:

- if no conditions: "false", []any{}

- if 1 condition: "<column> <soloOperator> ?", []any{value}

- if multiple conditions: "<column> <listOperator> (?, ?, ...)", values ([]any)
*/
func (c *klCondition) Build(t any) (string, []any) {
	column, values := c.pair.Build(t)
	numValues := len(values)
	if column == "" {
		return defaultConditionValues()
	} else if numValues == 0 {
		return defaultConditionValues()
	} else if numValues == 1 {
		return valueCondition(column, c.soloOperator), values
	} else {
		return listCondition(column, c.listOperator, numValues), values
	}
}

/*
Input: &struct, same struct used for Field=&struct.Field

Output:

- if no conditions: "false", []any{}

- if 1 condition: condition, values ([]any)

- if multiple conditions: "<condition> <operator> <condition> ...", values ([]any)
*/
func (cs *conditionSet) Build(t any) (string, []any) {
	numConditions := len(cs.conditions)
	if numConditions == 0 {
		return defaultConditionValues()
	} else if numConditions == 1 {
		return cs.conditions[0].Build(t)
	} else {
		conditions := make([]string, numConditions)
		allValues := make([]any, 0)
		for i, cond := range cs.conditions {
			c, values := cond.Build(t)
			if c == defaultCondition {
				// if any condition fails, return default condition (false)
				return defaultConditionValues()
			}
			conditions[i] = c
			allValues = append(allValues, values...)
		}
		glue := fmt.Sprintf(" %s ", cs.operator)
		allCondition := fmt.Sprintf("(%s)", strings.Join(conditions, glue)) // group by parentheses
		return allCondition, allValues
	}
}

/******************************** PRIVATE FUNCTIONS ********************************/

// Default condition response: "false", []any{}
func defaultConditionValues() (string, []any) {
	return defaultCondition, []any{}
}

// Used for single value conditions
func valueCondition(column, operator string) string {
	return fmt.Sprintf("%s %s ?", column, operator)
}

// Used for list values conditions
func listCondition(column, operator string, numValues int) string {
	placeholders := strings.Join(slices.Repeat([]string{"?"}, numValues), ", ")
	return fmt.Sprintf("%s %s (%s)", column, operator, placeholders)
}
