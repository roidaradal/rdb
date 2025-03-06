package query

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
