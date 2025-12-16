// Package rdb contains types and functions to perform type-safe SQL queries
package rdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/roidaradal/fn/check"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/rdb/internal/rdb"
)

// Initialize rdb package
func Initialize() error {
	rdb.Initialize()
	return nil
}

// Parameters for SQL connection
type SQLConnParams struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required"`
	Database string `validate:"required"`
}

// Create new MySQL DB connection pool
func NewSQLConnection(p *SQLConnParams) (*sql.DB, error) {
	if p == nil {
		return nil, errors.New("sql connection params are not set")
	}
	if check.NotValidStruct(p) {
		return nil, errors.New("missing required sql conn params")
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

// Add new type, extract its columns and fields, create RowCreator.
// StructRef is expected to be a struct pointer
var AddType = rdb.AddType

var (
	AllColumns = rdb.ColumnsOf      // Get all column names of given item's type
	Column     = rdb.GetColumnName  // Get column name of given field pointer
	Columns    = rdb.GetColumnNames // Get column names of given field pointers
	Field      = rdb.GetFieldName   // Get field name of given field pointer
	Fields     = rdb.GetFieldNames  // Get field names of given field pointers
)

// Function that reads row values into struct
type RowReader[T any] = rdb.RowReader[T]

// Create new RowReader[T] with given columns
func NewReader[T any](columns ...string) RowReader[T] {
	return rdb.NewReader[T](columns...)
}

// Create new RowReader[T] wusing all columns
func FullReader[T any](structRef *T) RowReader[T] {
	return rdb.FullReader(structRef)
}

// Convert given struct to map[string]any for row insertion
func ToRow[T any](structRef *T) dict.Object {
	return rdb.ToRow(structRef)
}
