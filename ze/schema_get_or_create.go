package ze

import (
	"fmt"

	"github.com/roidaradal/rdb"
)

type GetOrCreateParams[T any] struct {
	Name          string
	Owner         string
	PreCondition  rdb.Condition
	PostCondition rdb.Condition
	NewFn         func() *T
}

// Get or create item
func (s Schema[T]) GetOrCreate(rq *Request, cfg *GetOrCreateParams[T]) (*T, error) {
	return getOrCreate(rq, cfg, &s, false)
}

func (s Schema[T]) GetOrCreateTx(rqtx *Request, cfg *GetOrCreateParams[T]) (*T, error) {
	return getOrCreate(rqtx, cfg, &s, true)
}

// Common: get or create item
// Checks if item exists using preCondition; if not, the item is inserted into the table.
// Finally, get the item using preCondition and postCondition
func getOrCreate[T any](rq *Request, cfg *GetOrCreateParams[T], schema *Schema[T], isTx bool) (*T, error) {
	rqtx := rq
	// Check if item exists
	numRows, err := schema.Count(rq, cfg.PreCondition)
	if err != nil {
		rq.AddFmtLog("Failed to check if %s exists", cfg.Name)
		if isTx {
			err = rdb.Rollback(rqtx.DBTx, err) // manual rollback
		}
		return nil, err
	}

	if numRows > 1 {
		rq.AddFmtLog("Failed to get one %s", cfg.Name)
		err = fmt.Errorf("public: Multiple %s found", cfg.Name)
		if isTx {
			err = rdb.Rollback(rqtx.DBTx, err) // manual rollback
		}
		return nil, err
	} else if numRows == 0 {
		// Not found = create
		newItem := cfg.NewFn()
		if isTx {
			err = schema.InsertTx(rqtx, newItem)
		} else {
			err = schema.Insert(rq, newItem)
		}
		if err != nil {
			rq.AddFmtLog("Failed to create %s", cfg.Name)
			return nil, err
		}
		rq.AddFmtLog("Created %s for %s", cfg.Name, cfg.Owner)
	}

	// Fetch item
	item, err := schema.Get(rq, rdb.And(cfg.PreCondition, cfg.PostCondition))
	if err != nil {
		rq.AddFmtLog("Failed to get %s", cfg.Name)
		if isTx {
			err = rdb.Rollback(rqtx.DBTx, err) // manual rollback
		}
		return nil, err
	}
	return item, nil
}
