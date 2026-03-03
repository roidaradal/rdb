package query

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/rdb/internal/rdb"
)

// Group Count Query
type GroupCount[K comparable] struct {
	conditionQuery
	groupColumn string
}

// Create new GroupCount Query
func NewGroupCount[K comparable](table string, groupFieldRef *K) *GroupCount[K] {
	q := &GroupCount[K]{}
	q.initializeOptional(table)
	q.groupColumn = rdb.GetColumnName(groupFieldRef)
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
