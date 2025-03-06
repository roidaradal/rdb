package types

import (
	"fmt"
	"reflect"
)

type RowFn func(any) map[string]any

func NameOf(x any) string {
	if IsPointer(x) {
		return NameOf(Deref(x))
	}
	return reflect.TypeOf(x).Name()
}

func AddressOf(x any) string {
	return fmt.Sprintf("%p", x)
}

func Deref(x any) any {
	return reflect.ValueOf(x).Elem().Interface()
}
