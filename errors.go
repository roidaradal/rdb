package rdb

import "errors"

var (
	errNotStruct        = errors.New("not a struct")
	errIncompleteFields = errors.New("incomplete fields")
)

var (
	errEmptyQuery     = errors.New("empty query")
	errNoDBConnection = errors.New("nil db connection")
	errNoDBTx         = errors.New("nil dbtx")
	errNoRowReader    = errors.New("nil row reader")
	errResultCheck    = errors.New("result check failed")
)
