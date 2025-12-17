package query

import (
	"fmt"
	"slices"
	"strings"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/str"
)

// InsertRow Query
type InsertRow struct {
	baseQuery
	row dict.Object
}

// InsertRows Query
type InsertRows struct {
	baseQuery
	rows []dict.Object
}

// Create new InsertRow Query
func NewInsertRow(table string) *InsertRow {
	q := &InsertRow{}
	q.baseQuery.initialize(table)
	q.row = make(dict.Object)
	return q
}

// Create new InsertRows Query
func NewInsertRows(table string) *InsertRows {
	q := &InsertRows{}
	q.baseQuery.initialize(table)
	q.rows = make([]dict.Object, 0)
	return q
}

// Set InsertRow Query's row
func (q *InsertRow) Row(row dict.Object) {
	q.row = row
}

// Set InsertRows Query's rows
func (q *InsertRows) Rows(rows []dict.Object) {
	q.rows = rows
}

// Build InsertRow Query
func (q InsertRow) Build() (string, []any) {
	numColumns := len(q.row)
	err := q.baseQuery.preBuildCheck()
	if err != nil || numColumns == 0 {
		return emptyQueryValues()
	}
	columnList, values := dict.Unzip(q.row)
	columns := strings.Join(columnList, ", ")
	placeholders := str.Repeat(numColumns, "?", ", ")
	query := "INSERT INTO %s (%s) VALUES (%s)"
	query = fmt.Sprintf(query, q.table, columns, placeholders)
	return query, values
}

// Build InsertRows Query
func (q InsertRows) Build() (string, []any) {
	numRows := len(q.rows)
	err := q.baseQuery.preBuildCheck()
	if err != nil || numRows == 0 {
		return emptyQueryValues()
	}
	// Fix the order and number of columns based on first row
	row1 := q.rows[0]
	fixedOrder := columnOrder(row1)
	numColumns := len(row1)
	if numColumns == 0 {
		return emptyQueryValues()
	}
	values := make([]any, 0, numRows*numColumns)
	columnList, values1 := dict.Unzip(row1)
	values = append(values, values1...)
	for _, row := range q.rows[1:] {
		// Ensure same column order as first row
		if columnOrder(row) != fixedOrder {
			return emptyQueryValues()
		}
		// Follow row1's column order
		for _, column := range columnList {
			values = append(values, row[column])
		}
	}
	columns := strings.Join(columnList, ", ")
	placeholder := fmt.Sprintf("(%s)", str.Repeat(numColumns, "?", ", "))
	placeholders := str.Repeat(numRows, placeholder, ", ")
	query := "INSERT INTO %s (%s) VALUES %s"
	query = fmt.Sprintf(query, q.table, columns, placeholders)
	return query, values
}

// Join sorted column names
func columnOrder(row dict.Object) string {
	columns := dict.Keys(row)
	slices.Sort(columns)
	return strings.Join(columns, "/")
}
