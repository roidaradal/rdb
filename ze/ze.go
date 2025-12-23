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

var (
	ErrInactiveItem  = errors.New("public: Inactive item")
	ErrInvalidField  = errors.New("public: Invalid field")
	ErrMissingField  = errors.New("public: Missing required field")
	ErrMissingParams = errors.New("public: Missing required parameters")
	ErrNotFoundItem  = errors.New("public: Item not found")
	ErrMissingSchema = errors.New("schema is not initialized")
)

var (
	errMismatchCount  = errors.New("count mismatch")
	errNoDBConnection = errors.New("no db connection")
	errNoDBTx         = errors.New("no db transaction")
	errNoLastInsertID = errors.New("no last insert id")
	errNoRowsInserted = errors.New("no rows inserted")
)

var (
	Items  *Schema[Item] = nil // Items Schema
	dbConn *sql.DB       = nil // db connection pool
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

	return err
}

// Get Item reference object
func ItemsRef() *Item {
	if Items == nil {
		return nil
	}
	return Items.Ref
}
