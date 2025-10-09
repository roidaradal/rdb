package query

import "errors"

var (
	errNoDBConnection      = errors.New("no db connection")
	errNoDBTx              = errors.New("no db transaction")
	errNoReader            = errors.New("no row reader")
	errEmptyQuery          = errors.New("empty query")
	errEmptyTable          = errors.New("empty table")
	errNotFoundField       = errors.New("field not found")
	errFailedTypeAssertion = errors.New("type assertion failed")
	errFailedResultCheck   = errors.New("result check failed")
)
