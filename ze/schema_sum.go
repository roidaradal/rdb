package ze

import "github.com/roidaradal/rdb"

// SumQuery at schema.Table
func (s Schema[T]) Sum(rq *Request, columns []string, reader rdb.RowReader[T], condition rdb.Condition) (*T, error) {
	return sumAt(rq, columns, reader, condition, s.Table)
}

// SumQuery at table
func (s Schema[T]) SumAt(rq *Request, columns []string, reader rdb.RowReader[T], condition rdb.Condition, table string) (*T, error) {
	return sumAt(rq, columns, reader, condition, table)
}

// Common: create and execute SumQuery at given table
func sumAt[T any](rq *Request, columns []string, reader rdb.RowReader[T], condition rdb.Condition, table string) (*T, error) {
	// Build SumQuery and execute
	q := rdb.NewSumQuery(table, reader)
	q.Columns(columns)
	if condition != nil {
		q.Where(condition)
	}
	sum, err := q.Sum(rq.DB)
	if err != nil {
		rq.Status = Err500
		return nil, err
	}
	return sum, nil
}
