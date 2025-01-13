package rdb

import "reflect"

type rowCreator[T any] func(*T) map[string]any

/*
Input: columns []string

Output: rowCreator function for given columns
*/
func RowCreator[T any](columns []string) rowCreator[T] {
	return func(x *T) map[string]any {
		// Make sure T is struct
		if !isStructPointer(x) {
			return map[string]any{}
		}
		numColumns := len(columns)
		row := make(map[string]any, numColumns)
		for _, column := range columns {
			value, ok := findColumnValue(x, column)
			if !ok {
				continue
			}
			row[column] = value
		}
		// return empty if number of fields and columns mismatch
		if len(row) != numColumns {
			return map[string]any{}
		}
		return row
	}
}

/*************************** PRIVATE FUNCTIONS ***************************/

/*
Input: &struct, column string

Note: t is assumed to be isStructPtr-validated already, or isStruct validated and passed as &t

Output: field value matched by column tag, boolean to indicate whether column was found or not
*/
func findColumnValue(t any, column string) (any, bool) {
	structValue := reflect.ValueOf(t).Elem()
	structType := structValue.Type()
	for i := range structValue.NumField() {
		structField := structType.Field(i)
		if structField.Anonymous {
			// Embedded struct
			embeddedStructValue := structValue.FieldByName(structField.Name) // Get embedded struct as reflect.Value
			embeddedStructPointer := embeddedStructValue.Addr().Interface()  // Convert back to struct pointer
			fieldValue, ok := findColumnValue(embeddedStructPointer, column)
			if ok {
				return fieldValue, true
			}
		} else {
			// Normal field
			if getColumnTag(structField) == column {
				// Get struct value at this index
				fieldValue := structValue.Field(i).Interface()
				return fieldValue, true
			}
		}
	}
	return nil, false
}
