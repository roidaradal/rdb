package rdb

import (
	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/row"
	"github.com/roidaradal/rdb/internal/types"
)

func ToRow[T any](x *T) map[string]any {
	name := types.NameOf(x)
	rowFn, ok := memo.RowCreator[name]
	if !ok {
		return map[string]any{}
	}
	return rowFn(x)
}

func RowReader[T any](columns []string) row.RowReader[T] {
	return row.Reader[T](columns)
}
