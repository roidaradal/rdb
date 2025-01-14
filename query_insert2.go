package rdb

import (
	"database/sql"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type insertRowsQuery struct {
	table string
	rows  []map[string]any
}

/*
Output: Query (string), Values ([]any)

Note: Query could be blank string if invalid query parts
*/
func (q *insertRowsQuery) Build() (string, []any) {
	// Check if table is blank
	if q.table == "" {
		return defaultQueryValues() // return empty query if blank table
	}

	// Check if empty rows
	numRows := len(q.rows)
	if numRows == 0 {
		return defaultQueryValues() // return empty query if empty rows
	}

	// Get first row's column signature
	signature := columnSignature(q.rows[0])
	numColumns := len(q.rows[0])

	// Build columns from the first row
	columns := make([]string, 0, numColumns)
	values := make([]any, 0, numRows*numColumns)
	for column, value := range q.rows[0] {
		columns = append(columns, column)
		values = append(values, value)
	}

	// Add other rows
	for i := 1; i < numRows; i++ {
		row := q.rows[i]
		currSignature := columnSignature(row)
		if currSignature != signature {
			return defaultQueryValues() // return empty query if column signature mismatch
		}
		// Process row in columns order
		for _, column := range columns {
			value := row[column]
			values = append(values, value)
		}
	}

	// Build query
	cols := strings.Join(columns, ", ")
	placeholder := fmt.Sprintf("(%s)", repeatString(numColumns, "?", ", "))
	placeholders := repeatString(numRows, placeholder, ", ")
	query := "INSERT INTO %s (%s) VALUES %s"
	query = fmt.Sprintf(query, q.table, cols, placeholders)

	return query, values
}

/*
Input: rows []map[string]any

Constraint: Expects rows to have the same set of columns, column set will be based on the first row
*/
func (q *insertRowsQuery) Rows(rows []map[string]any) {
	q.rows = rows
}

/*
Input: initialized DB connection

Output: *sql.Result, error
*/
func (q *insertRowsQuery) Exec(dbc *sql.DB) (*sql.Result, error) {
	return prepareAndExec(q, dbc)
}

/*
Input: &struct, table (string)

Output: &insertRowsQuery
*/
func NewInsertRowsQuery(table string) *insertRowsQuery {
	q := insertRowsQuery{}
	q.table = table
	q.rows = make([]map[string]any, 0)
	return &q
}

/******************************** PRIVATE FUNCTIONS ********************************/

func columnSignature(row map[string]any) string {
	keys := slices.Collect(maps.Keys(row))
	slices.SortFunc(keys, func(key1, key2 string) int {
		return strings.Compare(key1, key2)
	})
	return strings.Join(keys, "/")
}
