package query

import (
	"database/sql"
	"fmt"
)

// Checks SQL result if condition is satisfied
type QueryResultChecker func(*sql.Result) bool

// QueryResultChecker that does nothing (used as default)
func AssertNothing(result *sql.Result) bool {
	return true // dont check results
}

// Create QueryResultChecker asserting number of rows affected
func AssertRowsAffected(target int) QueryResultChecker {
	return func(result *sql.Result) bool {
		return RowsAffected(result) == target
	}
}

// Gets the number of rows affected from SQL result,
// On error: returns 0
func RowsAffected(result *sql.Result) int {
	count := 0
	if result != nil {
		affected, err := (*result).RowsAffected()
		if err == nil {
			count = int(affected)
		}
	}
	return count
}

// Gets the last insert ID (uint) from SQL result,
// On error: returns 0
func LastInsertID(result *sql.Result) (uint, bool) {
	var insertID uint = 0
	ok := false
	if result != nil {
		insert, err := (*result).LastInsertId()
		if err == nil {
			insertID = uint(insert)
			ok = true
		}
	}
	return insertID, ok
}

// Executes an SQL query
func Exec(q Query, dbc *sql.DB) (*sql.Result, error) {
	query, values, err := preQueryCheck(q, dbc)
	if err != nil {
		return nil, err
	}
	stmt, err := dbc.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Executes an SQL query as part of a transaction,
// Applies Rollback on any errors
func ExecTx(q Query, dbtx *sql.Tx, checker QueryResultChecker) (*sql.Result, error) {
	var err error = nil
	query, values := q.Build()
	if dbtx == nil {
		err = errNoDBTx
	} else if query == "" {
		err = errEmptyQuery
	} else if checker == nil {
		err = errNoChecker
	}
	if err != nil {
		return nil, Rollback(dbtx, err)
	}

	stmt, err := dbtx.Prepare(query)
	if err != nil {
		return nil, Rollback(dbtx, err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return nil, Rollback(dbtx, err)
	}

	if ok := checker(&result); !ok {
		return nil, Rollback(dbtx, errFailedResultCheck)
	}

	return &result, nil
}

// Rolls back the SQL transaction
func Rollback(dbtx *sql.Tx, err error) error {
	err2 := dbtx.Rollback()
	if err2 != nil {
		// Combine original error and rollback error
		return fmt.Errorf("error: %w, rollback error: %w", err, err2)
	}
	// return original error if rollback successful
	return err
}
