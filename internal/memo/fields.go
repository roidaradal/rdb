package memo

import (
	"reflect"

	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/dyn"
	"github.com/roidaradal/fn/str"
)

const (
	rdbTag   string = "rdb"
	fxTag    string = "fx"
	required string = "must"
	editable string = "edit"
)

type fieldsResult struct {
	Required     []string
	Editable     []string
	Transformers map[string]TransformFn
}

// Get all required, editable fields and field transformers from the given struct pointer
func GetFieldsInfo(structRef any) *fieldsResult {
	result := &fieldsResult{
		Required:     make([]string, 0),
		Editable:     make([]string, 0),
		Transformers: make(map[string]TransformFn),
	}

	if !dyn.IsStructPointer(structRef) {
		return result
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
			embedded := GetFieldsInfo(embeddedStructRef)
			result.Required = append(result.Required, embedded.Required...)
			result.Editable = append(result.Editable, embedded.Editable...)
			result.Transformers = dict.Update(result.Transformers, embedded.Transformers)
		} else {
			// Normal field
			values := getRdbTagValues(structField)
			for _, value := range values {
				switch value {
				case required:
					result.Required = append(result.Required, fieldName)
				case editable:
					result.Editable = append(result.Editable, fieldName)
				}
			}
			fxKey := getFxKey(structField)
			if fn, ok := transformers[fxKey]; ok {
				result.Transformers[fieldName] = fn
			}
		}
	}
	return result
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

// Extract fx tag value from struct field
func getFxKey(structField reflect.StructField) string {
	return structField.Tag.Get(fxTag)
}
