package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/fn/str"
	"github.com/roidaradal/rdb/internal/condition"
	"github.com/roidaradal/rdb/internal/row"
)

// Query object only needs to implement the Build() method
type Query interface {
	// query string, values []any
	Build() (string, []any)
}

// Returns empty query, empty list for values
func emptyQueryValues() (string, []any) {
	return "", []any{}
}

// Base query type, with table
type baseQuery struct {
	table string
}

// Initialize BaseQuery
func (q *baseQuery) Initialize(table string) {
	q.table = str.WrapBackticks(table)
}

// Checks if the table is empty
func (q baseQuery) preBuildCheck() error {
	if q.table == "" {
		return errEmptyTable
	}
	return nil
}

// Condition query type, with table and required condition
type conditionQuery struct {
	baseQuery
	condition condition.Condition
}

// Initialize ConditionQuery
func (q *conditionQuery) Initialize(table string) {
	q.baseQuery.Initialize(table)
	q.condition = condition.None{} // if condition is not set later, defaults to false condition
}

// Set condition for ConditionQuery
func (q *conditionQuery) Where(queryCondition condition.Condition) {
	q.condition = queryCondition
}

// Calls the baseQuery.preBuildCheck and builds the condition
func (q conditionQuery) preBuildCheck() (string, []any, error) {
	err := q.baseQuery.preBuildCheck()
	condition, values := q.condition.Build()
	return condition, values, err
}

//

// OptionalCondition query type, with table and optional condition
type optionalConditionQuery struct {
	baseQuery
	condition condition.Condition
}

// Initialize OptionalConditionQuery
func (q *optionalConditionQuery) Initialize(table string) {
	q.baseQuery.Initialize(table)
	q.condition = condition.MatchAll{} // if condition is not set later, defaults to match all condition
}

// Set condition for OptionalConditionQuery
func (q *optionalConditionQuery) Where(queryCondition condition.Condition) {
	q.condition = queryCondition
}

// Calls the baseQuery.preBuildCheck and builds the condition
func (q optionalConditionQuery) preBuildCheck() (string, []any, error) {
	err := q.baseQuery.preBuildCheck()
	condition, values := q.condition.Build()
	return condition, values, err
}

// Checks the db connection and builds the query
func preQueryCheck(q Query, dbc *sql.DB) (string, []any, error) {
	var err error = nil
	query, values := q.Build()
	if dbc == nil {
		err = errNoDBConnection
	} else if query == "" {
		err = errEmptyQuery
	}
	return query, values, err
}

// Checks the db connection and reader, and builds the query
func preReadCheck[T any](q Query, dbc *sql.DB, reader row.RowReader[T]) (string, []any, error) {
	query, values, err := preQueryCheck(q, dbc)
	if err != nil {
		return query, values, err
	}
	if reader == nil {
		err = errNoReader
	}
	return query, values, err
}

// Builds the full query as a string
func QueryString(q Query) string {
	query, rawValues := q.Build()
	values := make([]any, len(rawValues))
	for i, value := range rawValues {
		typeName := fmt.Sprintf("%T", value)
		if strings.HasPrefix(typeName, "*") {
			values[i] = dyn.Deref(value)
		} else {
			values[i] = value
		}
	}
	// Replace all placeholders with %v format for values
	query = strings.Replace(query, "?", "%v", -1)
	return fmt.Sprintf(query, values...)
}
