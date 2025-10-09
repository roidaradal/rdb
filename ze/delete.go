package ze

import (
	"database/sql"
	"fmt"

	"github.com/roidaradal/rdb"
)

// DeleteParams
type DeleteParams struct {
	*DBTrio
	Condition rdb.Condition // required
	Table     string        // required for DeleteAt, DeleteTxAt
}

// DeleteQuery: schema.Table
func (s Schema[T]) Delete(p *DeleteParams) ([]string, error) {
	return deleteAt(&s, p, s.Table, false)
}

// DeleteQuery: params.Table
func (s Schema[T]) DeleteAt(p *DeleteParams) ([]string, error) {
	return deleteAt(&s, p, p.Table, false)
}

// DeleteQuery Transaction: schema.Table
func (s Schema[T]) DeleteTx(p *DeleteParams) ([]string, error) {
	return deleteAt(&s, p, s.Table, true)
}

// DeleteQuery Transaction: params.Table
func (s Schema[T]) DeleteTxAt(p *DeleteParams) ([]string, error) {
	return deleteAt(&s, p, p.Table, true)
}

// Common: create and execute DeleteQuery at given table using condition
func deleteAt[T any](s *Schema[T], p *DeleteParams, table string, isTx bool) ([]string, error) {
	logs := make([]string, 0)

	// Check that condition is set
	if p.Condition == nil {
		logs = append(logs, "Incomplete DeleteParams")
		return logs, errMissingParams
	}

	// Build DeleteQuery
	q := rdb.NewDeleteQuery(table)
	q.Where(p.Condition)

	// Execute DeleteQuery
	var result *sql.Result
	var err error
	if isTx {
		result, err = rdb.ExecTx(q, p.DBTx, p.Checker)
	} else {
		result, err = rdb.Exec(q, p.DB)
	}
	if err != nil {
		logs = append(logs, fmt.Sprintf("Failed to delete %s", s.Name))
		return logs, err
	}

	logs = append(logs, fmt.Sprintf("Deleted: %d", rdb.RowsAffected(result)))
	return logs, nil
}
