package rdb

import "github.com/roidaradal/rdb/internal/query"

var (
	Rollback = query.Rollback
	Exec     = query.Exec
	ExecTx   = query.ExecTx
)

var (
	RowsAffected       = query.RowsAffected
	LastInsertID       = query.LastInsertID
	AssertRowsAffected = query.AssertRowsAffected
	AssertNothing      = query.AssertNothing
)

type QueryResultChecker = query.QueryResultChecker
