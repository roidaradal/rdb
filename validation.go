package rdb

import "reflect"

func isStruct(t any) bool {
	return reflect.TypeOf(t).Kind() == reflect.Struct
}
