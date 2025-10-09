package ze

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/rdb"
)

// ToggleParams
type ToggleParams struct {
	*DBTrio
	IsActive bool   // required
	ID       ID     // required for Toggle*ID
	Code     string // required for Toggle*Code
	Table    string // required for Toggle*At
}

// ToggleID: schema.Table, by ID
func (s Schema[T]) ToggleID(p *ToggleParams) ([]string, error) {
	return toggleAt(&s, p, s.Table, true, false)
}

// ToggleID: params.Table, by ID
func (s Schema[T]) ToggleIDAt(p *ToggleParams) ([]string, error) {
	return toggleAt(&s, p, p.Table, true, false)
}

// ToggleCode: schema.Table, by Code
func (s Schema[T]) ToggleCode(p *ToggleParams) ([]string, error) {
	return toggleAt(&s, p, s.Table, false, false)
}

// ToggleCode: params.Table, by Code
func (s Schema[T]) ToggleCodeAt(p *ToggleParams) ([]string, error) {
	return toggleAt(&s, p, p.Table, false, false)
}

// ToggleTxID: schema.Table, transaction by ID
func (s Schema[T]) ToggleTxID(p *ToggleParams) ([]string, error) {
	return toggleAt(&s, p, s.Table, true, true)
}

// ToggleTxID: params.Table, transaction by ID
func (s Schema[T]) ToggleTxIDAt(p *ToggleParams) ([]string, error) {
	return toggleAt(&s, p, p.Table, true, true)
}

// ToggleTxCode: schema.Table, transaction by Code
func (s Schema[T]) ToggleTxCode(p *ToggleParams) ([]string, error) {
	return toggleAt(&s, p, s.Table, false, true)
}

// ToggleTxCode: params.Table, transaction by Code
func (s Schema[T]) ToggleTxCodeAt(p *ToggleParams) ([]string, error) {
	return toggleAt(&s, p, p.Table, false, true)
}

// Common: create and execute UpdateQuery, which toggles IsActive true/false,
// at given table, using ID/Code
func toggleAt[T any](s *Schema[T], p *ToggleParams, table string, byID bool, isTx bool) ([]string, error) {
	logs := make([]string, 0)

	// Check that params has ID or Code
	hasIdentity := false
	if byID && p.ID != 0 {
		hasIdentity = true
	} else if !byID && p.Code != "" {
		hasIdentity = true
	}
	if !hasIdentity {
		logs = append(logs, "Incomplete ToggleParams")
		return logs, errMissingParams
	}
	// Make sure Items schema is initialized
	if Items == nil {
		return logs, errMissingItems
	}

	// Build UpdateQuery using the Items schema
	item := Items.Ref
	q := rdb.NewUpdateQuery(table)
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
		result, err = rdb.ExecTx(q, p.DBTx, p.Checker)
	} else {
		result, err = rdb.Exec(q, p.DB)
	}
	if err != nil {
		logs = append(logs, fmt.Sprintf("Failed to toggle %s", s.Name))
		return logs, err
	}

	logs = append(logs, fmt.Sprintf("Toggled: %d", rdb.RowsAffected(result)))
	return logs, nil
}
