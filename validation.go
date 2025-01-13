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

func isNil(t any) bool {
	if t == nil {
		return true
	}
	switch reflect.TypeOf(t).Kind() {
	case reflect.Pointer, reflect.Map, reflect.Array, reflect.Slice, reflect.Chan:
		return reflect.ValueOf(t).IsNil()
	}
	return false
}
