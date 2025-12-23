package ze

import (
	"github.com/roidaradal/rdb"
)

// ToggleID at schema.Table
func (s Schema[T]) ToggleID(rq *Request, id ID, isActive bool) error {
	p := &toggleParams{id: id, isActive: isActive}
	return toggleAt[T](rq, p, s.Name, s.Table, true, false)
}

// ToggleID at table
func (s Schema[T]) ToggleIDAt(rq *Request, id ID, isActive bool, table string) error {
	p := &toggleParams{id: id, isActive: isActive}
	return toggleAt[T](rq, p, s.Name, table, true, false)
}

// ToggleCode at schema.Table
func (s Schema[T]) ToggleCode(rq *Request, code string, isActive bool) error {
	p := &toggleParams{code: code, isActive: isActive}
	return toggleAt[T](rq, p, s.Name, s.Table, false, false)
}

// ToggleCode at table
func (s Schema[T]) ToggleCodeAt(rq *Request, code string, isActive bool, table string) error {
	p := &toggleParams{code: code, isActive: isActive}
	return toggleAt[T](rq, p, s.Name, table, false, false)
}

// ToggleID transaction at schema.Table
func (s Schema[T]) ToggleTxID(rqtx *Request, id ID, isActive bool) error {
	p := &toggleParams{id: id, isActive: isActive}
	return toggleAt[T](rqtx, p, s.Name, s.Table, true, true)
}

// ToggleID transaction at table
func (s Schema[T]) ToggleTxIDAt(rqtx *Request, id ID, isActive bool, table string) error {
	p := &toggleParams{id: id, isActive: isActive}
	return toggleAt[T](rqtx, p, s.Name, table, true, true)
}

// ToggleCode transaction at schema.Table
func (s Schema[T]) ToggleTxCode(rqtx *Request, code string, isActive bool) error {
	p := &toggleParams{code: code, isActive: isActive}
	return toggleAt[T](rqtx, p, s.Name, s.Table, false, true)
}

// ToggleCode transaction at table
func (s Schema[T]) ToggleTxCodeAt(rqtx *Request, code string, isActive bool, table string) error {
	p := &toggleParams{code: code, isActive: isActive}
	return toggleAt[T](rqtx, p, s.Name, table, false, true)
}

type toggleParams struct {
	isActive bool   // required
	id       ID     // required for Toggle*ID
	code     string // required for Toggle*Code
}

// Common: create and execute UpdateQuery, which toggles IsActive true/false,
// at given table, using ID/Code
func toggleAt[T any](rq *Request, p *toggleParams, name, table string, byID bool, isTx bool) error {
	// Check that params has ID or Code
	hasIdentity := false
	if byID && p.id != 0 {
		hasIdentity = true
	} else if !byID && p.code != "" {
		hasIdentity = true
	}
	if !hasIdentity {
		rq.AddLog("ID/Code to toggle is not set")
		return ErrMissingParams
	}
	// Make sure Items schema is initialized
	if Items == nil {
		rq.AddLog("Items schema is null")
		return ErrMissingSchema
	}

	// Build UpdateQuery using the Items schema
	item := Items.Ref
	q := rdb.NewUpdateQuery[T](table)
	var condition1 rdb.Condition
	if byID {
		condition1 = rdb.Equal(&item.ID, p.id)
	} else {
		condition1 = rdb.Equal(&item.Code, p.code)
	}
	condition2 := rdb.Equal(&item.IsActive, !p.isActive)
	q.Where(rdb.And(condition1, condition2))
	rdb.Update(q, &item.IsActive, p.isActive)

	// Execute UpdateQuery
	var err error
	if isTx {
		rq.AddTxStep(q)
		_, err = rdb.ExecTx(q, rq.DBTx, rq.Checker)
	} else {
		_, err = rdb.Exec(q, rq.DB)
	}
	if err != nil {
		rq.AddFmtLog("Failed to toggle %s", name)
		rq.Status = Err500
		return err
	}

	// rq.AddFmtLog("Toggled: %d %s", rdb.RowsAffected(result), name)
	return nil
}
