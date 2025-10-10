package ze

import (
	"database/sql"

	"github.com/roidaradal/fn"
	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb"
)

type AddParams[T any] struct {
	Item  *T     // required for Insert*
	Items []*T   // required for InsertRows*
	Table string // required for Add*At
}

// Validate new item, check required fields, and apply transformations
func (s Schema[T]) ValidateNew(item *T) (*T, error) {
	// Validate struct
	if !check.IsValidStruct(item) {
		return nil, errMissingParams
	}

	for _, fieldName := range s.required {
		value := dyn.GetFieldValue(item, fieldName)
		// Apply transformer, if applicable
		if transform, ok := s.transformers[fieldName]; ok {
			value = transform(value)
		}
		// Check if zero value
		if dyn.IsZero(value) {
			return nil, errMissingField
		}
		// Update transformed field
		dyn.SetFieldValue(item, fieldName, value)
	}
	return item, nil
}

// InsertRowQuery: schema.Table
func (s Schema[T]) Insert(rq *Request, p *AddParams[T]) error {
	_, err := insertAt(rq, p, s.Name, s.Table, false, false)
	return err
}

// InsertRowQuery: params.Table
func (s Schema[T]) InsertAt(rq *Request, p *AddParams[T]) error {
	_, err := insertAt(rq, p, s.Name, p.Table, false, false)
	return err
}

// InsertRowQuery: schema.Table, with ID
func (s Schema[T]) InsertID(rq *Request, p *AddParams[T]) (ID, error) {
	return insertAt(rq, p, s.Name, s.Table, true, false)
}

// InsertRowQuery: params.Table, with ID
func (s Schema[T]) InsertIDAt(rq *Request, p *AddParams[T]) (ID, error) {
	return insertAt(rq, p, s.Name, p.Table, true, false)
}

// InsertRowQuery: schema.Table, as transaction
func (s Schema[T]) InsertTx(rq *Request, p *AddParams[T]) error {
	_, err := insertAt(rq, p, s.Name, s.Table, false, true)
	return err
}

// InsertRowQuery: params.Table, as transaction
func (s Schema[T]) InsertTxAt(rq *Request, p *AddParams[T]) error {
	_, err := insertAt(rq, p, s.Name, p.Table, false, true)
	return err
}

// InsertRowQuery: schema.Table, with ID, as transaction
func (s Schema[T]) InsertTxID(rq *Request, p *AddParams[T]) (ID, error) {
	return insertAt(rq, p, s.Name, s.Table, true, true)
}

// InsertRowQuery: params.Table, with ID, as transaction
func (s Schema[T]) InsertTxIDAt(rq *Request, p *AddParams[T]) (ID, error) {
	return insertAt(rq, p, s.Name, p.Table, true, true)
}

// InsertRowsQuery: schema.Table
func (s Schema[T]) InsertRows(rq *Request, p *AddParams[T]) error {
	return insertRowsAt(rq, p, s.Name, s.Table, false)
}

// InsertRowsQuery: params.Table
func (s Schema[T]) InsertRowsAt(rq *Request, p *AddParams[T]) error {
	return insertRowsAt(rq, p, s.Name, p.Table, false)
}

// InsertRowsQuery: schema.Table, as transaction
func (s Schema[T]) InsertTxRows(rq *Request, p *AddParams[T]) error {
	return insertRowsAt(rq, p, s.Name, s.Table, true)
}

// InsertRowsQuery: params.Table, as transaction
func (s Schema[T]) InsertTxRowsAt(rq *Request, p *AddParams[T]) error {
	return insertRowsAt(rq, p, s.Name, p.Table, true)
}

// Common: create and execute InsertRowQuery at given table
func insertAt[T any](rq *Request, p *AddParams[T], name, table string, getID bool, isTx bool) (ID, error) {
	var id ID = 0

	// Check that item is not null
	if p.Item == nil {
		rq.AddLog("Item to be added is null")
		return id, errMissingParams
	}

	// Build InsertRowQuery
	q := rdb.NewInsertRowQuery(table)
	q.Row(rdb.ToRow(p.Item))

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
		return id, err
	}
	rowsAffected := rdb.RowsAffected(result)

	// If not transaction, check if rowsAffected is 1
	if !isTx && rowsAffected != 1 {
		rq.AddFmtLog("No %s inserted", name)
		return id, errNoRowsInserted
	}

	// If getID flag is on, get the last insert ID
	if getID {
		var ok bool
		id, ok = rdb.LastInsertID(result)
		if !ok {
			rq.AddFmtLog("Failed to get %s insertID", name)
			return id, errNoLastInsertID
		}
	}

	rq.AddFmtLog("Added: %d", rowsAffected)
	return id, nil
}

// Common: create and execute InsertRowsQuery at given table
func insertRowsAt[T any](rq *Request, p *AddParams[T], name, table string, isTx bool) error {
	// Check that items are set
	if p.Items == nil {
		rq.AddLog("Items to be added are not set")
		return errMissingParams
	}
	numItems := len(p.Items)

	// Build InsertRowsQuery
	rows := fn.Map(p.Items, rdb.ToRow)
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
		return err
	}
	rowsAffected := rdb.RowsAffected(result)

	// If not transaction, check if rowsAffected == numItems
	if !isTx && rowsAffected != numItems {
		rq.AddFmtLog("Count mismatch: items = %d, rows = %d", numItems, rowsAffected)
		return errMismatchCount
	}

	rq.AddFmtLog("Added: %d", rowsAffected)
	return nil
}
