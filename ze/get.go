package ze

import (
	"github.com/roidaradal/fn"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/rdb"
)

// SelectRowQuery at schema.Table
func (s Schema[T]) Get(rq *Request, condition rdb.Condition) (*T, error) {
	return selectRowAt(rq, condition, s.Table, &s)
}

// SelectRowQuery at table
func (s Schema[T]) GetAt(rq *Request, condition rdb.Condition, table string) (*T, error) {
	return selectRowAt(rq, condition, table, &s)
}

// SelectRowQuery at schema.Table with pruning
func (s Schema[T]) GetOnly(rq *Request, condition rdb.Condition, fieldNames ...string) (*dict.Object, error) {
	item, err := selectRowAt(rq, condition, s.Table, &s)
	return prune(item, err, fieldNames...)
}

// SelectRowQuery at table with pruning
func (s Schema[T]) GetOnlyAt(rq *Request, condition rdb.Condition, table string, fieldNames ...string) (*dict.Object, error) {
	item, err := selectRowAt(rq, condition, table, &s)
	return prune(item, err, fieldNames...)
}

// SelectRowsQuery at schema.Table
func (s Schema[T]) GetRows(rq *Request, condition rdb.Condition) ([]*T, error) {
	return selectRowsAt(rq, condition, s.Table, &s)
}

// SelectRowsQuery at table
func (s Schema[T]) GetRowsAt(rq *Request, condition rdb.Condition, table string) ([]*T, error) {
	return selectRowsAt(rq, condition, table, &s)
}

// SelectRowsQuery at schema.Table with pruning
func (s Schema[T]) GetRowsOnly(rq *Request, condition rdb.Condition, fieldNames ...string) ([]*dict.Object, error) {
	items, err := selectRowsAt(rq, condition, s.Table, &s)
	return pruneRows(items, err)
}

// SelectRowsQuery at table with pruning
func (s Schema[T]) GetRowsOnlyAt(rq *Request, condition rdb.Condition, table string, fieldNames ...string) ([]*dict.Object, error) {
	items, err := selectRowsAt(rq, condition, table, &s)
	return pruneRows(items, err)
}

// SelectRowsQuery (all) at schema.Table
func (s Schema[T]) GetAllRows(rq *Request) ([]*T, error) {
	return selectRowsAt(rq, nil, s.Table, &s)
}

// SelectRowsQuery (all) at table
func (s Schema[T]) GetAllRowsAt(rq *Request, table string) ([]*T, error) {
	return selectRowsAt(rq, nil, table, &s)
}

// SelectRowsQuery (all) at schema.Table with pruning
func (s Schema[T]) GetAllRowsOnly(rq *Request, fieldNames ...string) ([]*dict.Object, error) {
	items, err := selectRowsAt(rq, nil, s.Table, &s)
	return pruneRows(items, err)
}

// SelectRowsQuery (all) at table with pruning
func (s Schema[T]) GetAllRowsOnlyAt(rq *Request, table string, fieldNames ...string) ([]*dict.Object, error) {
	items, err := selectRowsAt(rq, nil, table, &s)
	return pruneRows(items, err)
}

// Common: create and execute SelectRowQuery at given table
func selectRowAt[T any](rq *Request, condition rdb.Condition, table string, schema *Schema[T]) (*T, error) {
	// Check that condition is set
	if condition == nil {
		rq.AddLog("Condition is not set")
		rq.Status = Err500
		return nil, ErrMissingParams
	}

	// Build SelectRowQuery and execute
	q := rdb.NewFullSelectRowQuery(table, schema.Reader)
	q.Where(condition)
	item, err := q.QueryRow(rq.DB)
	if err != nil {
		rq.Status = Err500
		return nil, err
	}

	return item, nil
}

// Condition: create and execute SelectRowsQuery at given table
func selectRowsAt[T any](rq *Request, condition rdb.Condition, table string, schema *Schema[T]) ([]*T, error) {
	// Build SelectRowsQuery and execute
	q := rdb.NewFullSelectRowsQuery(table, schema.Reader)
	if condition != nil {
		q.Where(condition)
	}
	items, err := q.Query(rq.DB)
	if err != nil {
		rq.Status = Err500
		return nil, err
	}

	return items, nil
}

// Common: Prune item with given fieldNames
func prune[T any](item *T, err error, fieldNames ...string) (*dict.Object, error) {
	if err != nil {
		return nil, err
	}
	obj := dict.Prune(item, fieldNames...)
	return obj, nil
}

// Common: Prune items with given fieldNames
func pruneRows[T any](items []*T, err error, fieldNames ...string) ([]*dict.Object, error) {
	if err != nil {
		return nil, err
	}
	objs := fn.Map(items, func(item *T) *dict.Object {
		return dict.Prune(item, fieldNames...)
	})
	return objs, nil
}
