package ze

import (
	"database/sql"

	"github.com/roidaradal/rdb"
)

type ToggleParams struct {
	IsActive bool   // required
	ID       ID     // required for Toggle*ID
	Code     string // required for Toggle*Code
	Table    string // required for Toggle*At
}

// ToggleID: schema.Table, by ID
func (s Schema[T]) ToggleID(rq *Request, p *ToggleParams) error {
	return toggleAt[T](rq, p, s.Name, s.Table, true, false)
}

// ToggleID: params.Table, by ID
func (s Schema[T]) ToggleIDAt(rq *Request, p *ToggleParams) error {
	return toggleAt[T](rq, p, s.Name, p.Table, true, false)
}

// ToggleCode: schema.Table, by Code
func (s Schema[T]) ToggleCode(rq *Request, p *ToggleParams) error {
	return toggleAt[T](rq, p, s.Name, s.Table, false, false)
}

// ToggleCode: params.Table, by Code
func (s Schema[T]) ToggleCodeAt(rq *Request, p *ToggleParams) error {
	return toggleAt[T](rq, p, s.Name, p.Table, false, false)
}

// ToggleTxID: schema.Table, transaction by ID
func (s Schema[T]) ToggleTxID(rq *Request, p *ToggleParams) error {
	return toggleAt[T](rq, p, s.Name, s.Table, true, true)
}

// ToggleTxID: params.Table, transaction by ID
func (s Schema[T]) ToggleTxIDAt(rq *Request, p *ToggleParams) error {
	return toggleAt[T](rq, p, s.Name, p.Table, true, true)
}

// ToggleTxCode: schema.Table, transaction by Code
func (s Schema[T]) ToggleTxCode(rq *Request, p *ToggleParams) error {
	return toggleAt[T](rq, p, s.Name, s.Table, false, true)
}

// ToggleTxCode: params.Table, transaction by Code
func (s Schema[T]) ToggleTxCodeAt(rq *Request, p *ToggleParams) error {
	return toggleAt[T](rq, p, s.Name, p.Table, false, true)
}

// Common: create and execute UpdateQuery, which toggles IsActive true/false,
// at given table, using ID/Code
func toggleAt[T any](rq *Request, p *ToggleParams, name, table string, byID bool, isTx bool) error {
	// Check that params has ID or Code
	hasIdentity := false
	if byID && p.ID != 0 {
		hasIdentity = true
	} else if !byID && p.Code != "" {
		hasIdentity = true
	}
	if !hasIdentity {
		rq.AddLog("ID/Code to toggle is not set")
		return errMissingParams
	}
	// Make sure Items schema is initialized
	if Items == nil {
		rq.AddLog("Items schema is null")
		return errMissingItems
	}

	// Build UpdateQuery using the Items schema
	item := Items.Ref
	q := rdb.NewUpdateQuery[T](table)
	var condition1 rdb.Condition
	if byID {
		condition1 = rdb.Equal(&item.ID, p.ID)
	} else {
		condition1 = rdb.Equal(&item.Code, p.Code)
	}
	condition2 := rdb.Equal(&item.IsActive, !p.IsActive)
	q.Where(rdb.And(condition1, condition2))
	rdb.Update(q, &item.IsActive, p.IsActive)

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
		rq.AddFmtLog("Failed to toggle %s", name)
		return err
	}

	rq.AddFmtLog("Toggled: %d", rdb.RowsAffected(result))
	return nil
}
