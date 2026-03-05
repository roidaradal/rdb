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
