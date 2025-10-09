package ze

import "errors"

var (
	errMissingItems   = errors.New("items schema is not initialized")
	errMissingParams  = errors.New("missing required parameters")
	errMissingField   = errors.New("missing required field")
	errNoLastInsertID = errors.New("no last insert id")
	errNoRowsInserted = errors.New("no rows inserted")
)

var Items *Schema[Item] = nil

// Initialize the ze package
func Initialize() error {
	var err error
	Items, err = NewSharedSchema(&Item{})
	return err
}
