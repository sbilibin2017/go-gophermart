package services

import (
	"fmt"
	"reflect"
	"strconv"
)

func mapToStruct(result any, m map[string]any) error {
	val := reflect.ValueOf(result)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct, got %T", result)
	}
	elem := val.Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldName := elem.Type().Field(i).Name
		jsonTag := elem.Type().Field(i).Tag.Get("json")
		if jsonTag != "" {
			fieldName = jsonTag
		}
		if value, ok := m[fieldName]; ok {
			if field.CanSet() {
				val := reflect.ValueOf(value)
				if field.Kind() == val.Kind() {
					field.Set(val)
				} else {
					if field.Kind() == reflect.Int64 {
						if val.Kind() == reflect.Float64 {
							field.SetInt(int64(val.Float()))
						} else if val.Kind() == reflect.String {
							parsed, err := strconv.ParseInt(val.String(), 10, 64)
							if err == nil {
								field.SetInt(parsed)
							}
						} else if val.Kind() == reflect.Int {
							field.SetInt(val.Int())
						}
					} else if field.Kind() == reflect.String && val.Kind() == reflect.String {
						field.SetString(val.String())
					}
				}
			}
		}
	}
	return nil
}
