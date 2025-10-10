package ze

import (
	"database/sql"
	"net/http"

	"github.com/roidaradal/rdb"
)

type DeleteParams struct {
	Condition rdb.Condition // required
	Table     string        // required for Delete*At
}

// DeleteQuery: schema.Table
func (s Schema[T]) Delete(rq *Request, p *DeleteParams) error {
	return deleteAt(rq, p, s.Name, s.Table, false)
}

// DeleteQuery: params.Table
func (s Schema[T]) DeleteAt(rq *Request, p *DeleteParams) error {
	return deleteAt(rq, p, s.Name, p.Table, false)
}

// DeleteQuery Transaction: schema.Table
func (s Schema[T]) DeleteTx(rq *Request, p *DeleteParams) error {
	return deleteAt(rq, p, s.Name, s.Table, true)
}

// DeleteQuery Transaction: params.Table
func (s Schema[T]) DeleteTxAt(rq *Request, p *DeleteParams) error {
	return deleteAt(rq, p, s.Name, p.Table, true)
}

// Common: create and execute DeleteQuery at given table using condition
func deleteAt(rq *Request, p *DeleteParams, name, table string, isTx bool) error {
	// Check that condition is set
	if p.Condition == nil {
		rq.AddLog("Delete condition is not set")
		rq.Status = http.StatusBadRequest
		return errMissingParams
	}

	// Build DeleteQuery
	q := rdb.NewDeleteQuery(table)
	q.Where(p.Condition)

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
		rq.Status = http.StatusInternalServerError
		return err
	}

	rq.AddFmtLog("Deleted: %d", rdb.RowsAffected(result))
	rq.Status = http.StatusOK
	return nil
}
