package ze

import (
	"database/sql"

	"github.com/roidaradal/rdb"
)

// DeleteQuery at schema.Table
func (s Schema[T]) Delete(rq *Request, condition rdb.Condition) error {
	return deleteAt(rq, condition, s.Name, s.Table, false)
}

// DeleteQuery at table
func (s Schema[T]) DeleteAt(rq *Request, condition rdb.Condition, table string) error {
	return deleteAt(rq, condition, s.Name, table, false)
}

// DeleteQuery transaction at schema.Table
func (s Schema[T]) DeleteTx(rqtx *Request, condition rdb.Condition) error {
	return deleteAt(rqtx, condition, s.Name, s.Table, true)
}

// DeleteQuery transaction at table
func (s Schema[T]) DeleteTxAt(rqtx *Request, condition rdb.Condition, table string) error {
	return deleteAt(rqtx, condition, s.Name, table, true)
}

// Common: create and execute DeleteQuery at given table using condition
func deleteAt(rq *Request, condition rdb.Condition, name, table string, isTx bool) error {
	// Check that condition is set
	if condition == nil {
		rq.AddLog("Delete condition is not set")
		rq.Status = Err500
		return ErrMissingParams
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
		return err
	}

	rowsAffected := rdb.RowsAffected(result)
	if rowsAffected != 1 {
		rq.AddFmtLog("Deleted: %d %s", rowsAffected, name)
	}
	return nil
}
