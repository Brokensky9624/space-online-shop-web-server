package tool

import (
	"reflect"
)

func CopyFields(targetPointer interface{}, input reflect.Value) {
	targetValue := reflect.Indirect(reflect.ValueOf(targetPointer))
	for i := 0; i < input.NumField(); i++ {
		inputFieldName := input.Type().Field(i).Name
		inputField := input.Field(i)
		targetValueField := targetValue.FieldByName(inputFieldName)
		if targetValueField.IsValid() {
			if targetValueField.CanSet() {
				if inputField.Type() == targetValueField.Type() {
					targetValueField.Set(inputField)
				}
			}
		} else if inputField.Kind() == reflect.Struct {
			CopyFields(targetValue.Addr().Interface(), inputField)
		}
	}
}
