// Package query contains various types of SQL queries
package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/fn/lang"
	"github.com/roidaradal/fn/str"
	"github.com/roidaradal/rdb/internal/condition"
	"github.com/roidaradal/rdb/internal/rdb"
)

// Query object needs to implement the Build() method
// to output the query string and the parameter values
type Query interface {
	Build() (string, []any) // Return (query string, parameter values)
}

// Query, with table name
type baseQuery struct {
	table string
}

// Query, with table condition
type conditionQuery struct {
	baseQuery
	condition condition.Condition
}

// Initialize BaseQuery
func (q *baseQuery) initialize(table string) {
	q.table = str.WrapBackticks(table)
}

// Initialize ConditionQuery, with required condition
func (q *conditionQuery) initializeRequired(table string) {
	q.baseQuery.initialize(table)
	q.condition = condition.Missing{} // if condition not set later, defaults to false condtition
}

// Initialize ConditionQuery, with optional condition
func (q *conditionQuery) initializeOptional(table string) {
	q.baseQuery.initialize(table)
	q.condition = condition.MatchAll{} // if condition not set later, defaults to match all condition
}

// Check if table is set
func (q baseQuery) preBuildCheck() error {
	return lang.Ternary(q.table == "", errEmptyTable, nil)
}

// Check if table is set, and build the condition
func (q conditionQuery) preBuildCheck() (string, []any, error) {
	err := q.baseQuery.preBuildCheck()
	condition, values := q.condition.Build()
	return condition, values, err
}

// Set Query condition
func (q *conditionQuery) Where(queryCondition condition.Condition) {
	q.condition = queryCondition
}

// Before query, check the db connection and build the query
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

// Before SELECT query, check db connection, reader, and build the query
func preReadCheck[T any](q Query, dbc *sql.DB, reader rdb.RowReader[T]) (string, []any, error) {
	query, values, err := preQueryCheck(q, dbc)
	if err != nil {
		return query, values, err
	}
	if reader == nil {
		err = errNoReader
	}
	return query, values, err
}

// Build full query string
func ToString(q Query) string {
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

// Returns empty query and empty list of values
func emptyQueryValues() (string, []any) {
	return "", []any{}
}

// Read rows from query
func readRows[T any](dbc *sql.DB, query string, values []any, reader rdb.RowReader[T], task func(*T)) error {
	rows, err := dbc.Query(query, values...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		item, err := reader(rows)
		if err != nil {
			continue
		}
		task(item)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
