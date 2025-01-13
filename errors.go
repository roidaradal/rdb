package rdb

import "errors"

var (
	ErrNotStruct        = errors.New("not a struct")
	ErrIncompleteFields = errors.New("incomplete fields")
)
