package rdb

import (
	"database/sql"
	"fmt"
	"strings"
)

type selectQuery[T any] struct {
	conditionQuery[T]
	columns []string
	reader  RowReader[T]
	offset  uint
	limit   uint
	order   string
}

/*
Output: Query (string), Values ([]any)

Note: Query could be blank string if invalid query parts
*/
func (q *selectQuery[T]) Build() (string, []any) {
	return buildSelectQuery(q, true)
}

func buildSelectQuery[T any](q *selectQuery[T], includeLimit bool) (string, []any) {
	// Check if table is blank
	if q.table == "" {
		return defaultQueryValues() // return empty query if blank table
	}

	// Check if empty columns
	if len(q.columns) == 0 {
		return defaultQueryValues() // return empty query if empty columns
	}

	// Build columns
	columns := strings.Join(q.columns, ", ")

	// Build condition
	condition, values := q.condition.Build(q.object)

	// Build query
	query := "SELECT %s FROM %s WHERE %s"
	query = fmt.Sprintf(query, columns, q.table, condition)

	if q.order != "" {
		query = fmt.Sprintf("%s ORDER BY %s", query, q.order)
	}

	if includeLimit && q.limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d, %d", query, q.offset, q.limit)
	}

	return query, values
}

/*
Input: Columns []string

Note: Make sure corresponding Reader uses the same list of columns
*/
func (q *selectQuery[T]) Columns(columns []string) {
	q.columns = columns
}

/*
Input: limit uint
*/
func (q *selectQuery[T]) Limit(limit uint) {
	q.offset = 0
	q.limit = limit
}

/*
Input: page number, batch size uint
*/
func (q *selectQuery[T]) Page(number, batchSize uint) {
	q.offset = (number - 1) * batchSize
	q.limit = batchSize
}

/*
Input: column
*/
func (q *selectQuery[T]) OrderAsc(column string) {
	q.order = fmt.Sprintf("%s ASC", column)
}
func (q *selectQuery[T]) OrderDesc(column string) {
	q.order = fmt.Sprintf("%s DESC", column)
}

/*
Input: initialized DB connection

Output: list of structs that contain reader data, error
*/
func (q *selectQuery[T]) Query(dbc *sql.DB) ([]T, error) {
	if dbc == nil {
		return nil, errNoDBConnection
	}
	if q.reader == nil {
		return nil, errNoRowReader
	}
	query, values := q.Build()
	if query == "" {
		return nil, errEmptyQuery
	}

	rows, err := dbc.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]T, 0)
	for rows.Next() {
		item, err := q.reader(rows)
		if err != nil {
			continue
		}
		items = append(items, *item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

/*
Input: &struct, table (string), reader

Note: Same &struct will be used for setting conditions later

Output: &SelectQuery
*/
func NewSelectQuery[T any](object *T, table string, reader RowReader[T]) *selectQuery[T] {
	q := selectQuery[T]{}
	q.initialize(object, table)
	q.columns = make([]string, 0)
	q.reader = reader
	q.limit = 0 // default: no limit
	return &q
}
