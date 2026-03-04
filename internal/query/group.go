package query

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/rdb/internal/rdb"
)

type Number interface {
	~int | ~uint | ~int64 | ~float32 | ~float64
}

// Group Count Query
type GroupCount[K comparable] struct {
	conditionQuery
	groupColumn string
}

// Group Sum Query
type GroupSum[K comparable, V Number] struct {
	conditionQuery
	groupColumn string
	sumColumn   string
}

// Create new GroupCount Query
func NewGroupCount[K comparable](table string, groupFieldRef *K) *GroupCount[K] {
	q := &GroupCount[K]{}
	q.initializeOptional(table)
	q.groupColumn = rdb.GetColumnName(groupFieldRef)
	return q
}

// Create new GroupSum Query
func NewGroupSum[K comparable, V Number](table string, groupFieldRef *K, sumFieldRef *V) *GroupSum[K, V] {
	q := &GroupSum[K, V]{}
	q.initializeOptional(table)
	columns := rdb.GetColumnNames(groupFieldRef, sumFieldRef)
	q.groupColumn = columns[0]
	q.sumColumn = columns[1]
	return q
}

// Build GroupCount Query
func (q GroupCount[K]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.groupColumn == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s, COUNT(*) FROM %s WHERE %s GROUP BY %s"
	query = fmt.Sprintf(query, q.groupColumn, q.table, condition, q.groupColumn)
	return query, values
}

// Build GroupSum Query
func (q GroupSum[K, V]) Build() (string, []any) {
	condition, values, err := q.conditionQuery.preBuildCheck()
	if err != nil || q.groupColumn == "" || q.sumColumn == "" {
		return emptyQueryValues()
	}
	query := "SELECT %s, SUM(%s) FROM %s WHERE %s GROUP BY %s"
	query = fmt.Sprintf(query, q.groupColumn, q.sumColumn, q.table, condition, q.groupColumn)
	return query, values
}

// Execute GroupCountQuery and get map[group]count
func (q GroupCount[K]) GroupCount(dbc *sql.DB) (map[K]int, error) {
	query, values, err := preQueryCheck(q, dbc)
	if err != nil {
		return nil, err
	}

	rows, err := dbc.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[K]int)
	for rows.Next() {
		var key K
		var count int
		err = rows.Scan(&key, &count)
		if err != nil {
			continue
		}
		counts[key] = count
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}

// Execute GroupSumQuery and get map[group]sum
func (q GroupSum[K, V]) GroupSum(dbc *sql.DB) (map[K]V, error) {
	query, values, err := preQueryCheck(q, dbc)
	if err != nil {
		return nil, err
	}

	rows, err := dbc.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sums := make(map[K]V)
	for rows.Next() {
		var key K
		var sum V
		err = rows.Scan(&key, &sum)
		if err != nil {
			continue
		}
		sums[key] = sum
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return sums, nil
}
