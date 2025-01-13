package rdb

import "reflect"

// Unifies *sql.Row and *sql.Rows
type RowScannable interface {
	Scan(...any) error
}

type RowReader[T any] func(RowScannable) (*T, error)

/*
Input: columns []string

Output: RowReader function for given columns
*/
func Reader[T any](columns []string) RowReader[T] {
	return func(row RowScannable) (*T, error) {
		var x T
		// Make sure T is struct
		if !isStruct(x) {
			return nil, ErrNotStruct
		}
		numColumns := len(columns)
		fields := make([]any, 0, numColumns)
		for _, column := range columns {
			field, ok := findColumnField(&x, column)
			if !ok {
				continue
			}
			fields = append(fields, field)
		}
		// dont scan if number of fields and columns mismatch
		if len(fields) != numColumns {
			return nil, ErrIncompleteFields
		}
		err := row.Scan(fields...)
		return &x, err
	}
}

/*************************** PRIVATE FUNCTIONS ***************************/

/*
Input: &struct, column string

Note: t is assumed to be isStructPtr-validated already, or isStruct validated and passed as &t

Output: pointer to field matched by column tag, boolean to indicate whether column was found or not
*/
func findColumnField(t any, column string) (any, bool) {
	structValue := reflect.ValueOf(t).Elem()
	structType := structValue.Type()
	for i := range structValue.NumField() {
		structField := structType.Field(i)
		if structField.Anonymous {
			// Embedded struct
			embeddedStructValue := structValue.FieldByName(structField.Name) // Get embedded struct as reflect.Value
			embeddedStructPointer := embeddedStructValue.Addr().Interface()  // Convert back to struct pointer
			fieldAddress, ok := findColumnField(embeddedStructPointer, column)
			if ok {
				return fieldAddress, true
			}
		} else {
			// Normal field
			if getColumnTag(structField) == column {
				// Get struct value at this index
				fieldValue := structValue.Field(i)
				// Get address (as reflect.Value), convert back to pointer
				fieldAddress := fieldValue.Addr().Interface()
				return fieldAddress, true
			}
		}
	}
	return nil, false
}
