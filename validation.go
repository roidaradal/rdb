package rdb

import "reflect"

func isStruct(t any) bool {
	return reflect.TypeOf(t).Kind() == reflect.Struct
}

func isPointer(t any) bool {
	return reflect.TypeOf(t).Kind() == reflect.Pointer
}

func isStructPointer(t any) bool {
	if !isPointer(t) {
		return false
	}
	// Get underlying struct
	_struct := reflect.ValueOf(t).Elem().Interface()
	return isStruct(_struct)
}
