package ze

import (
	"fmt"

	"github.com/roidaradal/rdb"
)

// Checks if item exists using preCondition; if not, the item is inserted into the table.
// Finally, get the item using preCondition and postCondition
func (s Schema[T]) GetOrCreate(rq *Request, name, owner string, preCondition, postCondition rdb.Condition, newFn func() *T) (*T, error) {
	// Check if item exists
	numRows, err := s.Count(rq, preCondition)
	if err != nil {
		rq.AddFmtLog("Failed to check if %s exists", name)
		return nil, err
	}
	if numRows > 1 {
		rq.AddFmtLog("Failed to get one %s", name)
		return nil, fmt.Errorf("public: Multiple %s found", name)
	} else if numRows == 0 {
		// Not found = create
		newItem := newFn()
		err = s.Insert(rq, newItem)
		if err != nil {
			rq.AddFmtLog("Failed to create %s", name)
			return nil, err
		}
		rq.AddFmtLog("Created %s for %s", name, owner)
	}
	// Fetch item
	return s.Get(rq, rdb.And(preCondition, postCondition))
}
