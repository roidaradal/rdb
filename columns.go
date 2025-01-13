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
Input: t = struct

Output: List of column names from struct, extracted from struct tag `col:"Column"`

Example: AllColumns(Account{})
*/
func AllColumns(t any) []string {
	// Check if t is a struct
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

/*
Input: t = &struct, fields = &struct.Field, ...

Constraint: &struct has to be the same one used by &struct.Field

Constraint: cannot pass a field with a blank column tag

Output: List of column names corresponding to passed in fields

Output: Returns empty list if any of the fields have blank column tag

Example: a = Account{}; Columns(&a, &a.Name, &a.Type)
*/
func Columns(t any, fields ...any) []string {
	// Check if t is a struct pointer
	if !isStructPointer(t) {
		return []string{}
	}
	// Get reflect.Type and reflect.Value for struct
	structValue := reflect.ValueOf(t).Elem() // dereference pointer
	structType := structValue.Type()
	numFields := structValue.NumField()
	numColumns := len(fields)
	columns := make([]string, 0, numColumns)
	// Get column names for each field argument
	for _, field := range fields {
		// Check if field is pointer
		if !isPointer(field) {
			return []string{}
		}
		// Get the goal address
		// Note: Important to have &t and &t.Field have the same struct
		// otherwise this part doesnt work (e.g. &t, &p.Field)
		goalFieldValue := reflect.ValueOf(field).Elem()  // dereference pointer
		goalAddress := goalFieldValue.Addr().Interface() // get goal field value's address
		// Look through fields to find goal address
		for i := range numFields {
			structField := structType.Field(i)
			if structField.Anonymous {
				// Embedded struct
				embeddedStructValue := structValue.FieldByName(structField.Name) // Get embedded struct as reflect.Value
				embeddedStructPointer := embeddedStructValue.Addr().Interface()  // Convert back to struct pointer
				column := Column(embeddedStructPointer, field)
				if column != "" {
					columns = append(columns, column)
					break
				}
			} else {
				// Normal field
				fieldAddress := structValue.Field(i).Addr().Interface() // get current field value's address
				if fieldAddress == goalAddress {
					column := getColumnTag(structField)
					if column != "" {
						columns = append(columns, column)
					}
					break
				}
			}
		}
	}
	// Return empty list if columns and fields count mismatch
	if len(columns) != numColumns {
		return []string{}
	}
	return columns
}

/*
Input: t = &struct, field = &struct.Field

Constraint: &struct has to be the same one used by &struct.Field

Output: Column name corresponding to passed in field, or blank string

Example: a = Account{}; Column(&a, &a.Name)
*/
func Column(t any, field any) string {
	columns := Columns(t, field)
	if len(columns) == 0 {
		return ""
	}
	return columns[0]
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
