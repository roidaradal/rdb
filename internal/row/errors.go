package row

import "errors"

var (
	errNotStruct        = errors.New("type is not a struct")
	errIncompleteFields = errors.New("incomplete fields")
)
