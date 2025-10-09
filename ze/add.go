package ze

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/rdb"
)

// AddParams
type AddParams[T any] struct {
	*DBTrio
	Item  *T     // required
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
func (s Schema[T]) Insert(p *AddParams[T]) (ID, []string, error) {
	return insertAt(&s, p, s.Table, false, false)
}

// InsertRowQuery: params.Table
func (s Schema[T]) InsertAt(p *AddParams[T]) (ID, []string, error) {
	return insertAt(&s, p, p.Table, false, false)
}

// InsertRowQuery: schema.Table, with ID
func (s Schema[T]) InsertID(p *AddParams[T]) (ID, []string, error) {
	return insertAt(&s, p, s.Table, true, false)
}

// InsertRowQuery: params.Table, with ID
func (s Schema[T]) InsertIDAt(p *AddParams[T]) (ID, []string, error) {
	return insertAt(&s, p, p.Table, true, false)
}

// InsertRowQuery: schema.Table
func (s Schema[T]) InsertTx(p *AddParams[T]) (ID, []string, error) {
	return insertAt(&s, p, s.Table, false, true)
}

// InsertRowQuery: params.Table
func (s Schema[T]) InsertTxAt(p *AddParams[T]) (ID, []string, error) {
	return insertAt(&s, p, p.Table, false, true)
}

// InsertRowQuery: schema.Table, with ID
func (s Schema[T]) InsertTxID(p *AddParams[T]) (ID, []string, error) {
	return insertAt(&s, p, s.Table, true, true)
}

// InsertRowQuery: params.Table, with ID
func (s Schema[T]) InsertTxIDAt(p *AddParams[T]) (ID, []string, error) {
	return insertAt(&s, p, p.Table, true, true)
}

// Common: create and execute InsertRowQuery at given table
func insertAt[T any](s *Schema[T], p *AddParams[T], table string, getID bool, isTx bool) (ID, []string, error) {
	var id ID = 0
	logs := make([]string, 0)

	// Check that item is not null
	if p.Item == nil {
		logs = append(logs, "Incomplete AddParams")
		return id, logs, errMissingParams
	}

	// Build InsertRowQuery
	q := rdb.NewInsertRowQuery(table)
	q.Row(rdb.ToRow(p.Item))

	// Execute InsertRowQuery
	var result *sql.Result
	var err error
	if isTx {
		result, err = rdb.ExecTx(q, p.DBTx, p.Checker)
	} else {
		result, err = rdb.Exec(q, p.DB)
	}
	if err != nil {
		logs = append(logs, fmt.Sprintf("Failed to insert %s", s.Name))
		return id, logs, err
	}

	// If not transaction, check if rowsAffected is 1
	if !isTx && rdb.RowsAffected(result) != 1 {
		logs = append(logs, fmt.Sprintf("No %s inserted", s.Name))
		return id, logs, errNoRowsInserted
	}

	// If getID flag is on, get the last insert ID
	if getID {
		var ok bool
		id, ok = rdb.LastInsertID(result)
		if !ok {
			logs = append(logs, fmt.Sprintf("Failed to get %s insertID", s.Name))
			return id, logs, errNoLastInsertID
		}
	}

	return id, logs, nil
}
