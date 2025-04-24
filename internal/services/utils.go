package services

import (
	"fmt"
	"reflect"
)

func convertToMap(v interface{}) map[string]any {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return nil
	}
	result := make(map[string]any)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		result[fieldName] = field.Interface()
	}
	return result
}

func convertToStruct(structPtr interface{}, data map[string]any) error {
	structValue := reflect.ValueOf(structPtr)
	if structValue.Kind() != reflect.Ptr || structValue.IsNil() {
		return fmt.Errorf("expected a non-nil pointer to a struct")
	}
	structValue = structValue.Elem()
	if structValue.Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct, but got %s", structValue.Kind())
	}
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldName := structValue.Type().Field(i).Name
		if value, found := data[fieldName]; found {
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.New(field.Type().Elem()))
				fieldValue := field.Elem()
				if value != nil {
					fieldValue.Set(reflect.ValueOf(value))
				}
			} else if field.CanSet() {
				field.Set(reflect.ValueOf(value))
			}
		}
	}
	return nil
}
