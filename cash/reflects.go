package cash

import (
	"reflect"
)

func outValueOf(v interface{}) reflect.Value {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		panic("cash: non-pointer " + rv.Type().String())
	} else if rv.IsNil() {
		panic("cash: nil value")
	}
	return reflect.Indirect(rv)
}
