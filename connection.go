package rdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// Parameters for SQL connection
type SQLConnParams struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// Create new MySQL DB connection pool
func NewSQLConnection(p *SQLConnParams) (*sql.DB, error) {
	if p == nil {
		return nil, errors.New("sql connection params are not set")
	}
	dbAddr := fmt.Sprintf("%s:%s", p.Host, p.Port)
	dbCfg := mysql.Config{
		User:                 p.Username,
		Passwd:               p.Password,
		Net:                  "tcp",
		Addr:                 dbAddr,
		DBName:               p.Database,
		AllowNativePasswords: true,
	}
	dbc, err := sql.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("cannot open db conn: %w", err)
	}
	err = dbc.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot ping db conn: %w", err)
	}
	return dbc, nil
}
