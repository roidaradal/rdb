package query

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/roidaradal/rdb/internal/op"
)

type InsertRowQuery struct {
	baseQuery
	row map[string]any
}

type InsertRowsQuery struct {
	baseQuery
	rows []map[string]any
}

func (q *InsertRowQuery) Initialize(table string) {
	q.baseQuery.Initialize(table)
	q.row = make(map[string]any)
}

func (q *InsertRowsQuery) Initialize(table string) {
	q.baseQuery.Initialize(table)
	q.rows = make([]map[string]any, 0)
}

func (q *InsertRowQuery) Row(row map[string]any) {
	q.row = row
}

func (q *InsertRowsQuery) Rows(rows []map[string]any) {
	q.rows = rows
}

func (q InsertRowQuery) Build() (string, []any) {
	count := len(q.row)
	if q.table == "" || count == 0 {
		return defaultQueryValues()
	}
	columns := make([]string, 0, count)
	values := make([]any, 0, count)
	for column, value := range q.row {
		columns = append(columns, column)
		values = append(values, value)
	}
	cols := strings.Join(columns, ", ")
	placeholders := op.RepeatString(count, "?", ", ")
	query := "INSERT INTO %s (%s) VALUES (%s)"
	query = fmt.Sprintf(query, q.table, cols, placeholders)
	return query, values
}

func (q InsertRowsQuery) Build() (string, []any) {
	numRows := len(q.rows)
	if q.table == "" || numRows == 0 {
		return defaultQueryValues()
	}
	firstSignature := columnSignature(q.rows[0])
	numColumns := len(q.rows[0])
	columns := make([]string, 0, numColumns)
	values := make([]any, 0, numRows*numColumns)
	for column, value := range q.rows[0] {
		columns = append(columns, column)
		values = append(values, value)
	}
	for i := 1; i < numRows; i++ {
		row := q.rows[i]
		if columnSignature(row) != firstSignature {
			return defaultQueryValues()
		}
		for _, column := range columns {
			values = append(values, row[column])
		}
	}
	cols := strings.Join(columns, ", ")
	placeholder := fmt.Sprintf("(%s)", op.RepeatString(numColumns, "?", ", "))
	placeholders := op.RepeatString(numRows, placeholder, ", ")
	query := "INSERT INTO %s (%s) VALUES %s"
	query = fmt.Sprintf(query, q.table, cols, placeholders)
	return query, values
}

func columnSignature(row map[string]any) string {
	keys := slices.Collect(maps.Keys(row))
	slices.SortFunc(keys, strings.Compare)
	return strings.Join(keys, "/")
}
