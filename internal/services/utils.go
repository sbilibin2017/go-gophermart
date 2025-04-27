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
					} else if field.Kind() == reflect.Float64 {
						if val.Kind() == reflect.Int64 {
							field.SetFloat(float64(val.Int()))
						} else if val.Kind() == reflect.String {
							parsed, err := strconv.ParseFloat(val.String(), 64)
							if err == nil {
								field.SetFloat(parsed)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func mapListToStruct(result any, m []map[string]any) error {
	// Проверяем, что result - это срез указателей на структуры
	val := reflect.ValueOf(result)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("expected a pointer to a slice, got %T", result)
	}

	// Получаем срез, на который указывает result
	slice := val.Elem()
	if slice.Kind() != reflect.Slice {
		return fmt.Errorf("expected a slice, got %T", slice)
	}

	// Проверка на соответствие длины
	if slice.Len() != len(m) {
		return fmt.Errorf("length of result slice and map slice must match")
	}

	// Для каждого элемента списка маппим данные
	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i)
		if elem.Kind() != reflect.Ptr || elem.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("expected a pointer to a struct at index %d, got %T", i, elem)
		}

		// Преобразуем map в структуру
		err := mapToStruct(elem.Interface(), m[i])
		if err != nil {
			return fmt.Errorf("error mapping item at index %d: %v", i, err)
		}
	}

	return nil
}
