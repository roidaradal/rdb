package query

import "errors"

var (
	errEmptyQuery     = errors.New("empty query")
	errNoDBConnection = errors.New("no db connection")
	errNoDBTx         = errors.New("no dbtx connection")
	errNoReader       = errors.New("no row reader")
	errResultCheck    = errors.New("result check failed")
)

var (
	errFieldNotFound = errors.New("field not found")
	errTypeAssertion = errors.New("type assertion failed")
)
