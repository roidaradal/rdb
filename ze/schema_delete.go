package ze

import (
	"database/sql"

	"github.com/roidaradal/fn/fail"
	"github.com/roidaradal/rdb"
)

// DeleteQuery at schema.Table
func (s Schema[T]) Delete(rq *Request, condition rdb.Condition) error {
	_, err := deleteAt(rq, condition, s.Name, s.Table, false)
	return err
}

// DeleteQuery at table
func (s Schema[T]) DeleteAt(rq *Request, condition rdb.Condition, table string) error {
	_, err := deleteAt(rq, condition, s.Name, table, false)
	return err
}

// DeleteQuery transaction at schema.Table
func (s Schema[T]) DeleteTx(rqtx *Request, condition rdb.Condition) error {
	_, err := deleteAt(rqtx, condition, s.Name, s.Table, true)
	return err
}

// DeleteQuery transaction at table
func (s Schema[T]) DeleteTxAt(rqtx *Request, condition rdb.Condition, table string) error {
	_, err := deleteAt(rqtx, condition, s.Name, table, true)
	return err
}

// DeleteQuery at schema.Table, return rowsAffected
func (s Schema[T]) CountDelete(rq *Request, condition rdb.Condition) (int, error) {
	return deleteAt(rq, condition, s.Name, s.Table, false)
}

// DeleteQuery at table, return rowsAffected
func (s Schema[T]) CountDeleteAt(rq *Request, condition rdb.Condition, table string) (int, error) {
	return deleteAt(rq, condition, s.Name, table, false)
}

// DeleteQuery transaction at schema.Table, return rowsAffected
func (s Schema[T]) CountDeleteTx(rqtx *Request, condition rdb.Condition) (int, error) {
	return deleteAt(rqtx, condition, s.Name, s.Table, true)
}

// DeleteQuery transaction at table, return rowsAffected
func (s Schema[T]) CountDeleteTxAt(rqtx *Request, condition rdb.Condition, table string) (int, error) {
	return deleteAt(rqtx, condition, s.Name, table, true)
}

// Common: create and execute DeleteQuery at given table using condition
func deleteAt(rq *Request, condition rdb.Condition, name, table string, isTx bool) (int, error) {
	// Check that condition is set
	if condition == nil {
		rq.AddLog("Delete condition is not set")
		rq.Status = Err500
		return 0, fail.MissingParams
	}

	// Build DeleteQuery
	q := rdb.NewDeleteQuery(table)
	q.Where(condition)

	// Execute DeleteQuery
	var result *sql.Result
	var err error
	if isTx {
		rq.AddTxStep(q)
		result, err = rdb.ExecTx(q, rq.DBTx, rq.Checker)
	} else {
		result, err = rdb.Exec(q, rq.DB)
	}
	if err != nil {
		rq.AddFmtLog("Failed to delete %s", name)
		rq.Status = Err500
		return 0, err
	}

	rowsAffected := rdb.RowsAffected(result)
	if rowsAffected != 1 {
		rq.AddFmtLog("Deleted: %d %s", rowsAffected, name)
	}
	return rowsAffected, nil
}
