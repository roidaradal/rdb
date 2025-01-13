package rdb

import "errors"

var (
	errNotStruct        = errors.New("not a struct")
	errIncompleteFields = errors.New("incomplete fields")
)
