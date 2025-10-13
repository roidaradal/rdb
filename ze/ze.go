package ze

import (
	"database/sql"
	"errors"

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
	errInvalidField  = errors.New("public: Invalid field")
	errMissingParams = errors.New("public: Missing required parameters")
	errMissingField  = errors.New("public: Missing required field")
)

var (
	Items  *Schema[Item] = nil
	dbConn *sql.DB       = nil
)

const Dot string = memo.Dot

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
