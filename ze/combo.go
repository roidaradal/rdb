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
