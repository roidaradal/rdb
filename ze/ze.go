package ze

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/roidaradal/rdb"
	"github.com/roidaradal/rdb/internal/memo"
)

var (
	errMismatchCount  = errors.New("count mismatch")
	errMissingItems   = errors.New("items schema is not initialized")
	errNoDBConnection = errors.New("no db connection")
	errNoDBTx         = errors.New("no db transaction")
	errNoLastInsertID = errors.New("no last insert id")
	errNoRowsInserted = errors.New("no rows inserted")
)

var (
	ErrInvalidField  = errors.New("public: Invalid field")
	ErrMissingField  = errors.New("public: Missing required field")
	ErrMissingParams = errors.New("public: Missing required parameters")
	ErrMissingSchema = errors.New("schema is not initialized")
)

var (
	Items  *Schema[Item] = nil
	dbConn *sql.DB       = nil
)

const Dot string = memo.Dot

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

// Initialize the ze package:
// Creates the Items schema,
// Initializes the db connection pool
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
