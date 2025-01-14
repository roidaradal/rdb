package rdb

import "database/sql"

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
