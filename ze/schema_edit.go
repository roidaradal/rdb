package ze

import (
	"database/sql"

	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/fn/fail"
	"github.com/roidaradal/rdb"
)

// Get field updates by comparing existing object and patch object
func (s Schema[T]) FieldUpdates(rq *Request, oldItem *T, patchObject dict.Object) (*T, rdb.FieldUpdates, error) {
	updates := make(rdb.FieldUpdates)
	for _, fieldName := range s.editable {
		if dict.NoKey(patchObject, fieldName) {
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
				return nil, nil, fail.InvalidField
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
		return nil, nil, fail.InvalidField
	}

	return oldItem, updates, nil
}

// UpdateQuery at schema.Table
func (s Schema[T]) Update(rq *Request, updates rdb.FieldUpdates, condition rdb.Condition) error {
	return updateAt[T](rq, updates, condition, s.Name, s.Table, false)
}

// UpdateQuery at table
func (s Schema[T]) UpdateAt(rq *Request, updates rdb.FieldUpdates, condition rdb.Condition, table string) error {
	return updateAt[T](rq, updates, condition, s.Name, table, false)
}

// UpdateQuery transaction at schema.Table
func (s Schema[T]) UpdateTx(rqtx *Request, updates rdb.FieldUpdates, condition rdb.Condition) error {
	return updateAt[T](rqtx, updates, condition, s.Name, s.Table, true)
}

// UpdateQuery transaction at table
func (s Schema[T]) UpdateTxAt(rqtx *Request, updates rdb.FieldUpdates, condition rdb.Condition, table string) error {
	return updateAt[T](rqtx, updates, condition, s.Name, table, true)
}

// Common: create and execute UpdateQuery at given table
func updateAt[T any](rq *Request, updates rdb.FieldUpdates, condition rdb.Condition, name, table string, isTx bool) error {
	// Check that condition and updates are set
	if condition == nil || updates == nil {
		rq.AddLog("Condition/updates not set")
		rq.Status = Err500
		return fail.MissingParams
	}

	// Build UpdateQuery
	q := rdb.NewUpdateQuery[T](table)
	q.Where(condition)
	q.Updates(updates)

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
