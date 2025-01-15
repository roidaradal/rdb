package rdb

import "database/sql"

type QueryResultChecker func(*sql.Result) bool

func AssertRowsAffected(target int) QueryResultChecker {
	return func(result *sql.Result) bool {
		return RowsAffected(result) == target
	}
}

func AssertNothing() QueryResultChecker {
	return func(result *sql.Result) bool {
		return true // dont check results
	}
}

func Rollback(dbtx *sql.Tx, err error) (*sql.Result, error) {
	err2 := dbtx.Rollback()
	if err2 != nil {
		return nil, err2
	}
	return nil, err
}

func prepareAndExecTx(q BuildableQuery, dbtx *sql.Tx, resultChecker QueryResultChecker) (*sql.Result, error) {
	if dbtx == nil {
		return nil, errNoDBTx
	}

	query, values := q.Build()
	if query == "" {
		return Rollback(dbtx, errEmptyQuery)
	}

	stmt, err := dbtx.Prepare(query)
	if err != nil {
		return Rollback(dbtx, err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return Rollback(dbtx, err)
	}

	if ok := resultChecker(&result); !ok {
		return Rollback(dbtx, errResultCheck)
	}

	return &result, nil
}

func prepareAndExec(q BuildableQuery, dbc *sql.DB) (*sql.Result, error) {
	if dbc == nil {
		return nil, errNoDBConnection
	}

	query, values := q.Build()
	if query == "" {
		return nil, errEmptyQuery
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
