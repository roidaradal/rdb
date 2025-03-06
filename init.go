package rdb

import (
	"fmt"

	"github.com/roidaradal/rdb/internal/columns"
	"github.com/roidaradal/rdb/internal/memo"
	"github.com/roidaradal/rdb/internal/row"
	"github.com/roidaradal/rdb/internal/types"
)

func Initialize(schemas ...any) error {
	memo.Initialize()
	for _, schema := range schemas {
		name := types.NameOf(schema)
		if !types.IsStructPointer(schema) {
			return fmt.Errorf("invalid schema: %s", name)
		}
		result := columns.All(schema)
		memo.Instance[name] = schema
		memo.AllColumns[name] = result.Columns
		memo.ColumnField[name] = result.ColumnField
		memo.UpdateColumnAddress(result.AddressOf)
		memo.RowCreator[name] = row.Creator(result.Columns)
	}
	return nil
}

func Schema[T any](t T) *T {
	return memo.InstanceOf(t)
}
