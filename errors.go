package rdb

import "errors"

var (
	errNotStruct        = errors.New("not a struct")
	errIncompleteFields = errors.New("incomplete fields")
)

var (
	errEmptyQuery     = errors.New("empty query")
	errNoDBConnection = errors.New("nil db connection")
	errNoRowReader    = errors.New("nil row reader")
)
