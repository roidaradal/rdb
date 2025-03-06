package columns

import (
	"fmt"
	"reflect"

	"github.com/roidaradal/rdb/internal/types"
)

const (
	columnTag    string = "col"
	skipTagValue string = "-"
)

type Result struct {
	Columns     []string
	AddressOf   map[string]string
	ColumnField map[string]string
}

func newResult() *Result {
	return &Result{
		Columns:     make([]string, 0),
		AddressOf:   make(map[string]string),
		ColumnField: make(map[string]string),
	}
}

func getColumnName(field reflect.StructField) string {
	column := field.Tag.Get(columnTag)
	if column == skipTagValue {
		return ""
	} else if column != "" {
		return column
	}
	return field.Name
}

func getAddress(schemaValue reflect.Value, index int) string {
	return fmt.Sprintf("%v", schemaValue.Field(index).Addr().Interface())
}

func All(schema any) *Result {
	result := newResult()
	if !types.IsStructPointer(schema) {
		return result
	}
	schemaValue := reflect.ValueOf(schema).Elem() // dereference pointer
	schemaType := schemaValue.Type()
	numFields := schemaType.NumField()
	for i := range numFields {
		field := schemaType.Field(i)
		if field.Anonymous {
			// Embedded struct
			embeddedStructPtr := schemaValue.FieldByName(field.Name).Addr().Interface()
			embedded := All(embeddedStructPtr)
			result.Columns = append(result.Columns, embedded.Columns...)
			for k, v := range embedded.AddressOf {
				result.AddressOf[k] = v
			}
			for k, v := range embedded.ColumnField {
				result.ColumnField[k] = v
			}
		} else {
			// Normal field
			column := getColumnName(field)
			if column == "" {
				continue
			}
			column = fmt.Sprintf("`%s`", column) // wrap in backticks
			address := getAddress(schemaValue, i)
			result.Columns = append(result.Columns, column)
			result.AddressOf[address] = column
			result.ColumnField[column] = field.Name
		}
	}
	return result
}
