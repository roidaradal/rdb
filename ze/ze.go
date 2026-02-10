// Package ze contains the Schema and Request types
package ze

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/roidaradal/rdb"
)

const (
	OK200  = http.StatusOK                  // OK
	OK201  = http.StatusCreated             // Created
	Err400 = http.StatusBadRequest          // client-side error
	Err401 = http.StatusUnauthorized        // unauthenticated
	Err403 = http.StatusForbidden           // unauthorized
	Err404 = http.StatusNotFound            // not found
	Err429 = http.StatusTooManyRequests     // rate limiting
	Err500 = http.StatusInternalServerError // server-side error
)

var ErrMissingSchema = errors.New("schema is not initialized")

var (
	errMismatchCount  = errors.New("count mismatch")
	errNoDBConnection = errors.New("no db connection")
	errNoDBTx         = errors.New("no db transaction")
	errNoLastInsertID = errors.New("no last insert id")
	errNoRowsInserted = errors.New("no rows inserted")
)

var (
	Items     *Schema[Item]      = nil // Items Schema
	dbConn    *sql.DB            = nil // db connection pool
	dbConnMap map[string]*sql.DB = nil // map of custom db connection pools
)

// Initialize ze package; create Items schema and initialize db connection pool
func Initialize(dbConnParams *rdb.SQLConnParams) error {
	var err error

	// Create Items schema
	Items, err = NewSharedSchema(&Item{})
	if err != nil {
		return err
	}

	// Create db connection pool
	dbConn, err = rdb.NewSQLConnection(dbConnParams)
	if err != nil {
		return err
	}

	// Initialize custom db connection pools
	dbConnMap = make(map[string]*sql.DB)

	return nil
}

// Add custom DB connection
func AddDBConnection(name string, dbConnParams *rdb.SQLConnParams) error {
	customDBConn, err := rdb.NewSQLConnection(dbConnParams)
	if err != nil {
		return err
	}
	dbConnMap[name] = customDBConn
	return nil
}

// Add a new schema with the given table, add to errors list if applicable
func AddSchema[T any](item *T, table string, errs []error) *Schema[T] {
	schema, err := NewSchema(item, table)
	if err != nil {
		errs = append(errs, err)
	}
	return schema
}

// Add shared schema (no table), add to errors list if applicable
func AddSharedSchema[T any](item *T, errs []error) *Schema[T] {
	schema, err := NewSharedSchema(item)
	if err != nil {
		errs = append(errs, err)
	}
	return schema
}

// Get Item reference object
func ItemsRef() *Item {
	if Items == nil {
		return nil
	}
	return Items.Ref
}
