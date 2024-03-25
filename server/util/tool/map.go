package tool

import "reflect"

func StructToMap(a interface{}) map[string]interface{} {
	tps := reflect.TypeOf(a)
	vs := reflect.ValueOf(a)
	ret := make(map[string]interface{})
	for i := 0; i < tps.NumField(); i++ {
		tp := tps.Field(i)
		key := tp.Tag.Get("json")
		if key == "" {
			key = tp.Name
		}
		val := vs.Field(i).Interface()
		ret[key] = val
	}
	return ret
}
