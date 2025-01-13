package rdb

import (
	"reflect"
)

const (
	columnTag    string = "col"
	skipTagValue string = "-"
)

/******************************** PUBLIC FUNCTIONS ********************************/

/*
Constraint: T has to be a struct

Output: List of column names from struct, extracted from struct tag `col:"Column"`
*/
func AllColumns(t any) []string {
	// Check if T is a struct
	if !isStruct(t) {
		return []string{}
	}
	// Get reflect.Type and reflect.Value for struct
	structType := reflect.TypeOf(t)
	structValue := reflect.ValueOf(t)
	numFields := structType.NumField()
	// Collect column names from column tags
	columns := make([]string, 0, numFields)
	for i := range numFields {
		field := structType.Field(i)
		if field.Anonymous {
			// Embedded struct
			fieldValue := structValue.FieldByName(field.Name) // Get embedded struct as reflect.Value
			embeddedStruct := fieldValue.Interface()          // Convert back to struct
			embeddedColumns := AllColumns(embeddedStruct)     // Use recursion
			columns = append(columns, embeddedColumns...)
		} else {
			// Normal field
			column := getColumnTag(field)
			if column == "" {
				continue
			}
			columns = append(columns, column)
		}
	}
	return columns
}

/******************************** PRIVATE FUNCTIONS ********************************/

/*
Input: reflect.StructField

Output: column name found in struct tag `col:"Column"`

Returns empty string if column tag not found or `col:"-"`
*/
func getColumnTag(field reflect.StructField) string {
	column := field.Tag.Get(columnTag)
	if column == skipTagValue {
		column = ""
	}
	return column
}
