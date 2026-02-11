package ze

import (
	"database/sql"

	"github.com/roidaradal/fn/fail"
	"github.com/roidaradal/rdb"
)

// SetFlag: updateQuery of boolean=flag at schema.Table
func (s Schema[T]) SetFlag(rq *Request, condition rdb.Condition, field *bool, flag bool) error {
	return setFlagAt[T](rq, condition, field, flag, s.Table, s.Name, false)
}

// SetFlagAt: updateQuery of boolean=flag at table
func (s Schema[T]) SetFlagAt(rq *Request, condition rdb.Condition, field *bool, flag bool, table string) error {
	return setFlagAt[T](rq, condition, field, flag, table, s.Name, false)
}

// SetTxFlag: updateQuery transaction of boolean=flag at schema.Table
func (s Schema[T]) SetTxFlag(rqtx *Request, condition rdb.Condition, field *bool, flag bool) error {
	return setFlagAt[T](rqtx, condition, field, flag, s.Table, s.Name, true)
}

// SetTxFlagAt: updateQuery transaction of boolean=flag at table
func (s Schema[T]) SetTxFlagAt(rqtx *Request, condition rdb.Condition, field *bool, flag bool, table string) error {
	return setFlagAt[T](rqtx, condition, field, flag, table, s.Name, true)
}

// Common: create and execute UpdateQuery of boolean=flag at given table
func setFlagAt[T any](rq *Request, condition rdb.Condition, field *bool, flag bool, table, name string, isTx bool) error {
	// Check that condition is set
	if condition == nil {
		rq.AddLog("Condition not set")
		rq.Status = Err500
		return fail.MissingParams
	}

	// Build UpdateQuery
	q := rdb.NewUpdateQuery[T](table)
	q.Where(condition)
	rdb.Update(q, field, flag)

	// Execute UpdateQuery
	var result *sql.Result
	var err error
	if isTx {
		rq.AddTxStep(q)
		result, err = rdb.ExecTx(q, rq.DBTx, rq.Checker)
	} else {
		result, err = rdb.Exec(q, rq.DB)
	}
	if err != nil {
		rq.AddFmtLog("Failed to update %s", name)
		rq.Status = Err500
		return err
	}

	rowsAffected := rdb.RowsAffected(result)
	if rowsAffected != 1 {
		rq.AddFmtLog("Updated: %d %s", rowsAffected, name)
	}
	return nil
}
