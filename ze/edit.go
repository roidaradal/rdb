package ze

import (
	"database/sql"

	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb"
)

type UpdateParams struct {
	Updates   rdb.FieldUpdates // required
	Condition rdb.Condition    // required
	Table     string           // required for Update*At
}

// Get field updates by comparing existing object and patch object
func (s Schema[T]) FieldUpdates(rq *Request, oldItem *T, patchObject dict.Object) (*T, rdb.FieldUpdates, error) {
	updates := make(rdb.FieldUpdates)
	for _, fieldName := range s.editable {
		if !dict.HasKey(patchObject, fieldName) {
			continue // skip if fieldName is not in patch
		}
		newValue := patchObject[fieldName]
		oldValue := dyn.GetFieldValue(oldItem, fieldName)
		// Apply transformer, if applicable
		if transform, ok := s.transformers[fieldName]; ok {
			newValue = transform(newValue)
		}
		// Check custom validator, if any
		if validator, ok := s.validators[fieldName]; ok {
			if !validator(newValue) {
				rq.Status = Err400
				return nil, nil, errInvalidField
			}
		}
		// Add field update if old and new values are not equal
		if dyn.NotEqual(oldValue, newValue) {
			updates[fieldName] = rdb.FieldUpdate{oldValue, newValue}
			dyn.SetFieldValue(oldItem, fieldName, newValue)
		}
	}

	// Validate old item with updated values
	if !check.IsValidStruct(oldItem) {
		rq.Status = Err400
		return nil, nil, errInvalidField
	}

	return oldItem, updates, nil
}

// UpdateQuery: schema.Table
func (s Schema[T]) Update(rq *Request, p *UpdateParams) error {
	return updateAt[T](rq, p, s.Name, s.Table, false)
}

// UpdateQuery: params.Table
func (s Schema[T]) UpdateAt(rq *Request, p *UpdateParams) error {
	return updateAt[T](rq, p, s.Name, p.Table, false)
}

// UpdateQuery Transaction: schema.Table
func (s Schema[T]) UpdateTx(rq *Request, p *UpdateParams) error {
	return updateAt[T](rq, p, s.Name, s.Table, true)
}

// UpdateQuery Transaction: params.Table
func (s Schema[T]) UpdateTxAt(rq *Request, p *UpdateParams) error {
	return updateAt[T](rq, p, s.Name, p.Table, true)
}

// Common: create and execute UpdateQuery at given table
func updateAt[T any](rq *Request, p *UpdateParams, name, table string, isTx bool) error {
	// Check that condition and updates are set
	if p.Condition == nil || p.Updates == nil {
		rq.AddLog("Condition/updates not set")
		rq.Status = Err400
		return errMissingParams
	}

	// Build UpdateQuery
	q := rdb.NewUpdateQuery[T](table)
	q.Where(p.Condition)
	q.Updates(p.Updates)

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

	rq.AddFmtLog("Updated: %d", rdb.RowsAffected(result))
	rq.Status = OK200
	return nil
}
