package ze

import "github.com/roidaradal/rdb"

// CountQuery at schema.Table
func (s Schema[T]) Count(rq *Request, condition rdb.Condition) (int, error) {
	return countAt(rq, condition, s.Table)
}

// CountQuery at table
func (s Schema[T]) CountAt(rq *Request, condition rdb.Condition, table string) (int, error) {
	return countAt(rq, condition, table)
}

// Common: create and execute CountQuery at given table
func countAt(rq *Request, condition rdb.Condition, table string) (int, error) {
	// Build CountQuery and execute
	q := rdb.NewCountQuery(table)
	if condition != nil {
		q.Where(condition)
	}
	count, err := q.Count(rq.DB)
	if err != nil {
		rq.Status = Err500
		return 0, err
	}
	return count, nil
}
