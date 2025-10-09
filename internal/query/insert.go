package query

import (
	"fmt"
	"slices"
	"strings"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/str"
)

type InsertRowQuery struct {
	baseQuery
	row dict.Object
}

type InsertRowsQuery struct {
	baseQuery
	rows []dict.Object
}

// Initialize InsertRowQuery
func (q *InsertRowQuery) Initialize(table string) {
	q.baseQuery.Initialize(table)
	q.row = make(dict.Object)
}

// Initialize InsertRowsQuery
func (q *InsertRowsQuery) Initialize(table string) {
	q.baseQuery.Initialize(table)
	q.rows = make([]dict.Object, 0)
}

// Set insert row
func (q *InsertRowQuery) Row(row dict.Object) {
	q.row = row
}

// Set insert rows
func (q *InsertRowsQuery) Rows(rows []dict.Object) {
	q.rows = rows
}

// Build the InsertRowQuery
func (q InsertRowQuery) Build() (string, []any) {
	numColumns := len(q.row)
	err := q.baseQuery.preBuildCheck()
	if err != nil || numColumns == 0 {
		return emptyQueryValues()
	}
	columns, values := dict.Unzip(q.row)
	cols := strings.Join(columns, ", ")
	placeholders := str.Repeat(numColumns, "?", ", ")
	query := "INSERT INTO %s (%s) VALUES (%s)"
	query = fmt.Sprintf(query, q.table, cols, placeholders)
	return query, values
}

// Build the InsertRowsQuery
func (q InsertRowsQuery) Build() (string, []any) {
	numRows := len(q.rows)
	err := q.baseQuery.preBuildCheck()
	if err != nil || numRows == 0 {
		return emptyQueryValues()
	}
	row1 := q.rows[0]
	signature1 := columnSignature(row1)
	numColumns := len(row1)
	if numColumns == 0 {
		return emptyQueryValues()
	}
	values := make([]any, 0, numRows*numColumns)
	columns, values1 := dict.Unzip(row1)
	values = append(values, values1...)
	for i := 1; i < numRows; i++ {
		row := q.rows[i]
		// Ensure same column signature as first row
		if columnSignature(row) != signature1 {
			return emptyQueryValues()
		}
		// Follows row1's column order
		for _, column := range columns {
			values = append(values, row[column])
		}
	}
	cols := strings.Join(columns, ", ")
	placeholder := fmt.Sprintf("(%s)", str.Repeat(numColumns, "?", ", "))
	placeholders := str.Repeat(numRows, placeholder, ", ")
	query := "INSERT INTO %s (%s) VALUES %s"
	query = fmt.Sprintf(query, q.table, cols, placeholders)
	return query, values
}

// Join the sorted column names as the signature
func columnSignature(row dict.Object) string {
	columns := dict.Keys(row)
	slices.Sort(columns)
	return strings.Join(columns, "/")
}
