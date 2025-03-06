package types

import "reflect"

func IsStruct(x any) bool {
	return reflect.TypeOf(x).Kind() == reflect.Struct
}

func IsPointer(x any) bool {
	return reflect.TypeOf(x).Kind() == reflect.Pointer
}

func IsStructPointer(x any) bool {
	if !IsPointer(x) {
		return false
	}
	return IsStruct(Deref(x))
}

func IsNil(x any) bool {
	if x == nil {
		return true
	}
	switch reflect.TypeOf(x).Kind() {
	case reflect.Pointer, reflect.Map, reflect.Array, reflect.Slice, reflect.Chan:
		return reflect.ValueOf(x).IsNil()
	}
	return false
}
