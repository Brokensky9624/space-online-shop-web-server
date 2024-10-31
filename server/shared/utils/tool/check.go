package tool

import (
	"fmt"
	"reflect"
)

func CheckRequiredFields(param interface{}) error {
	t := reflect.TypeOf(param)
	v := reflect.ValueOf(param)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).String()
		jsonTag := field.Tag.Get("json")
		requiredTag := field.Tag.Get("required")
		if requiredTag == "true" && value == "" {
			return fmt.Errorf("%s is required but empty", jsonTag)
		}
	}
	return nil
}
