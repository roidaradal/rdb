package rdb

import (
	"github.com/roidaradal/rdb/internal/memo"
)

func AllColumns(x any) []string {
	return memo.ColumnsOf(x)
}

func Column(field any) string {
	return memo.GetColumn(field)
}

func Columns(fields ...any) []string {
	return memo.GetColumns(fields...)
}
