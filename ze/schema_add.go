package ze

import (
	"database/sql"

	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/fn/fail"
	"github.com/roidaradal/fn/list"
	"github.com/roidaradal/rdb"
)

// Validate new item, check required fields, and apply transformations
func (s Schema[T]) ValidateNew(rq *Request, item *T) (*T, error) {
	// Validate struct
	if !check.IsValidStruct(item) {
		rq.Status = Err400
		return nil, fail.MissingParams
	}

	for _, fieldName := range s.required {
		value := dyn.GetFieldValue(item, fieldName)
		// Apply transformer, if applicable
		if transform, ok := s.transformers[fieldName]; ok {
			value = transform(value)
		}
		// Check custom validator, if any
		if validator, ok := s.validators[fieldName]; ok {
			if !validator(value) {
				rq.Status = Err400
				return nil, fail.InvalidField
			}
		}
		// Check if zero value
		if dyn.IsZero(value) {
			rq.Status = Err400
			return nil, fail.MissingField
		}
		// Update transformed field
		dyn.SetFieldValue(item, fieldName, value)
	}
	return item, nil
}

// InsertRowQuery at schema.Table
func (s Schema[T]) Insert(rq *Request, item *T) error {
	_, err := insertAt(rq, item, s.Name, s.Table, false, false)
	return err
}

// InsertRowQuery at table
func (s Schema[T]) InsertAt(rq *Request, item *T, table string) error {
	_, err := insertAt(rq, item, s.Name, table, false, false)
	return err
}

// InsertRowQuery at schema.Table with ID
func (s Schema[T]) InsertID(rq *Request, item *T) (ID, error) {
	return insertAt(rq, item, s.Name, s.Table, true, false)
}

// InsertRowQuery at table with ID
func (s Schema[T]) InsertIDAt(rq *Request, item *T, table string) (ID, error) {
	return insertAt(rq, item, s.Name, table, true, false)
}

// InsertRowQuery transaction at schema.Table
func (s Schema[T]) InsertTx(rqtx *Request, item *T) error {
	_, err := insertAt(rqtx, item, s.Name, s.Table, false, true)
	return err
}

// InsertRowQuery transaction at table
func (s Schema[T]) InsertTxAt(rqtx *Request, item *T, table string) error {
	_, err := insertAt(rqtx, item, s.Name, table, false, true)
	return err
}

// InsertRowQuery transaction at schema.Table with ID
func (s Schema[T]) InsertTxID(rqtx *Request, item *T) (ID, error) {
	return insertAt(rqtx, item, s.Name, s.Table, true, true)
}

// InsertRowQuery transaction at table, with ID
func (s Schema[T]) InsertTxIDAt(rqtx *Request, item *T, table string) (ID, error) {
	return insertAt(rqtx, item, s.Name, table, true, true)
}

// InsertRowsQuery at schema.Table
func (s Schema[T]) InsertRows(rq *Request, items []*T) error {
	return insertRowsAt(rq, items, s.Name, s.Table, false)
}

// InsertRowsQuery at table
func (s Schema[T]) InsertRowsAt(rq *Request, items []*T, table string) error {
	return insertRowsAt(rq, items, s.Name, table, false)
}

// InsertRowsQuery transaction at schema.Table
func (s Schema[T]) InsertTxRows(rqtx *Request, items []*T) error {
	return insertRowsAt(rqtx, items, s.Name, s.Table, true)
}

// InsertRowsQuery transaction at table
func (s Schema[T]) InsertTxRowsAt(rqtx *Request, items []*T, table string) error {
	return insertRowsAt(rqtx, items, s.Name, table, true)
}

// Common: create and execute InsertRowQuery at given table
func insertAt[T any](rq *Request, item *T, name, table string, getID bool, isTx bool) (ID, error) {
	var id ID = 0

	// Check that item is not null
	if item == nil {
		rq.AddLog("Item to be added is null")
		return id, fail.MissingParams
	}

	// Build InsertRowQuery
	q := rdb.NewInsertRowQuery(table)
	q.Row(rdb.ToRow(item))

	// Execute InsertRowQuery
	var result *sql.Result
	var err error
	if isTx {
		rq.AddTxStep(q)
		result, err = rdb.ExecTx(q, rq.DBTx, rq.Checker)
	} else {
		result, err = rdb.Exec(q, rq.DB)
	}
	if err != nil {
		rq.AddFmtLog("Failed to insert %s", name)
		rq.Status = Err500
		return id, err
	}
	rowsAffected := rdb.RowsAffected(result)

	// If not transaction, check if rowsAffected is 1
	if !isTx && rowsAffected != 1 {
		rq.AddFmtLog("No %s inserted", name)
		rq.Status = Err500
		return id, errNoRowsInserted
	}

	// If getID flag is on, get the last insert ID
	if getID {
		var ok bool
		id, ok = rdb.LastInsertID(result)
		if !ok {
			rq.AddFmtLog("Failed to get %s insertID", name)
			rq.Status = Err500
			return id, errNoLastInsertID
		}
	}

	// rq.AddFmtLog("Added: %d %s", rowsAffected, name)
	rq.Status = OK201
	return id, nil
}

// Common: create and execute InsertRowsQuery at given table
func insertRowsAt[T any](rq *Request, items []*T, name, table string, isTx bool) error {
	// Check that items are set
	if items == nil {
		rq.AddLog("Items to be added are not set")
		return fail.MissingParams
	}
	numItems := len(items)

	// Build InsertRowsQuery
	rows := list.Map(items, rdb.ToRow)
	q := rdb.NewInsertRowsQuery(table)
	q.Rows(rows)

	// Execute InsertRowsQuery
	var result *sql.Result
	var err error
	if isTx {
		rq.AddTxStep(q)
		checker := rdb.AssertRowsAffected(numItems)
		result, err = rdb.ExecTx(q, rq.DBTx, checker)
	} else {
		result, err = rdb.Exec(q, rq.DB)
	}
	if err != nil {
		rq.AddFmtLog("Failed to insert %d %s rows", numItems, name)
		rq.Status = Err500
		return err
	}
	rowsAffected := rdb.RowsAffected(result)

	// If not transaction, check if rowsAffected == numItems
	if !isTx && rowsAffected != numItems {
		rq.AddFmtLog("Count mismatch: items = %d, rows = %d", numItems, rowsAffected)
		rq.Status = Err500
		return errMismatchCount
	}

	rq.AddFmtLog("Added: %d %s", rowsAffected, name)
	rq.Status = OK201
	return nil
}
