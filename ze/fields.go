package ze

import (
	"reflect"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/fn/str"
)

const (
	rdbTag   string = "rdb"  // Struct tag for rdb required and editable fields
	fxTag    string = "fx"   // Struct tag for transformation function
	required string = "must" // Struct tag value for required field
	editable string = "edit" // Struct tag value for editable field
)

type fieldsInfo struct {
	required     []string
	editable     []string
	transformers map[string]TransformFn
}

// Get all required, editable fields and field transformers from given struct pointer
func getFieldsInfo(structRef any) fieldsInfo {
	fields := fieldsInfo{
		required:     make([]string, 0),
		editable:     make([]string, 0),
		transformers: make(map[string]TransformFn),
	}

	if !dyn.IsStructPointer(structRef) {
		return fields
	}

	structValue := reflect.ValueOf(structRef).Elem() // dereference pointer
	structType := structValue.Type()
	numFields := structType.NumField()

	for i := range numFields {
		structField := structType.Field(i)
		fieldName := structField.Name
		if structField.Anonymous {
			// Embedded struct, use recursion
			embeddedStructRef := structValue.FieldByName(fieldName).Addr().Interface()
			embedded := getFieldsInfo(embeddedStructRef)
			fields.required = append(fields.required, embedded.required...)
			fields.editable = append(fields.editable, embedded.editable...)
			fields.transformers = dict.Update(fields.transformers, embedded.transformers)
		} else {
			// Normal field
			values := getRdbTagValues(structField)
			for _, value := range values {
				switch value {
				case required:
					fields.required = append(fields.required, fieldName)
				case editable:
					fields.editable = append(fields.editable, fieldName)
				}
			}
			fxKey := structField.Tag.Get(fxTag)
			if fn, ok := transformers[fxKey]; ok {
				fields.transformers[fieldName] = fn
			}
		}
	}
	return fields
}

// Extract rdb tag values from struct field
func getRdbTagValues(structField reflect.StructField) []string {
	values := make([]string, 0)
	tagValue := structField.Tag.Get(rdbTag)
	if tagValue != "" {
		values = str.CleanSplit(tagValue, ",")
	}
	return values
}
