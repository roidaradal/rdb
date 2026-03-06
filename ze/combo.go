package ze

import "github.com/roidaradal/rdb"

// Insert item to insertSchema and delete correpsonding item from deleteSchema, as part of a transaction
func MoveItem[T1, T2 any](rqtx *Request, insertSchema *Schema[T1], item *T1, deleteSchema *Schema[T2], deleteCondition rdb.Condition) error {
	// 1) Insert item to insertSchema
	err := insertSchema.InsertTx(rqtx, item)
	if err != nil {
		rqtx.AddFmtLog("Failed to insert %s", insertSchema.Name)
		return err
	}
	// 2) Delete from deleteSchema
	err = deleteSchema.DeleteTx(rqtx, deleteCondition)
	if err != nil {
		rqtx.AddFmtLog("Failed to delete %s", deleteSchema.Name)
		return err
	}
	return nil
}

// Update one item and fetch it
func UpdateAndGet[T any](rqtx *Request, schema *Schema[T], setUpdatesFn func(*rdb.UpdateQuery[T]), updateCondition, selectCondition rdb.Condition) (*T, error) {
	// Update one item
	q := rdb.NewUpdateQuery[T](schema.Table)
	q.Where(updateCondition)
	q.Limit(1)
	setUpdatesFn(q)
	rqtx.AddTxStep(q)
	_, err := rdb.ExecTx(q, rqtx.DBTx, rqtx.Checker)
	if err != nil {
		return nil, err
	}
	// Select one item
	item, err := schema.Get(rqtx, selectCondition)
	if err != nil {
		err = rdb.Rollback(rqtx.DBTx, err) // Manual rollback on error of Get
		return nil, err
	}
	return item, nil
}

// Gets the item and locks it
// Note: no need to include IsLocked = true/false in conditions, as this function adds it
func GetAndLock[T any](rqtx *Request, schema *Schema[T], lockField *bool, selectCondition rdb.Condition, lockConditionFn func(*T) rdb.Condition) (*T, error) {
	// Get unlocked item
	condition := rdb.And(
		selectCondition,
		rdb.Equal(lockField, false),
	)
	item, err := schema.Get(rqtx, condition)
	if err != nil {
		// Manual rollback on error of Get
		err = rdb.Rollback(rqtx.DBTx, err)
		return nil, err
	}
	// Lock item
	lockCondition := lockConditionFn(item)
	condition2 := rdb.And(
		lockCondition,
		rdb.Equal(lockField, false),
	)
	err = schema.SetTxFlag(rqtx, condition2, lockField, true)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Gets a list of items and locks all of them
// Note: no need to include IsLocked = true/false in conditions, as this function adds it
func GetAndLockItems[T any](rqtx *Request, schema *Schema[T], lockField *bool, selectCondition rdb.Condition, lockConditionFn func([]*T) rdb.Condition, numItems int) ([]*T, error) {
	// Get unlocked items
	condition := rdb.And(
		selectCondition,
		rdb.Equal(lockField, false),
	)
	items, err := schema.GetRows(rqtx, condition)
	if err != nil {
		// Manual rollback on error of GetRows
		err = rdb.Rollback(rqtx.DBTx, err)
		return nil, err
	}
	if len(items) != numItems {
		rqtx.Status = Err500
		// Manual rollback if mismatch count
		err = rdb.Rollback(rqtx.DBTx, errMismatchCount)
		return nil, err
	}
	// Lock items
	lockCondition := lockConditionFn(items)
	condition2 := rdb.And(
		lockCondition,
		rdb.Equal(lockField, false),
	)
	err = schema.SetTxFlag(rqtx, condition2, lockField, true)
	if err != nil {
		return nil, err
	}
	return items, nil
}
