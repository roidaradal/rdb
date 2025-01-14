package rdb

import "database/sql"

func prepareAndExec(q buildableQuery, dbc *sql.DB) (*sql.Result, error) {
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
