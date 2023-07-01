package utils

import "reflect"

func GetValue(v interface{}, field string) reflect.Value {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f
}
