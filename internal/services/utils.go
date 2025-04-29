package services

import (
	"fmt"
	"reflect"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func mapToStruct(target interface{}, data map[string]any) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}
	structValue := targetValue.Elem()
	for key, value := range data {
		field := structValue.FieldByName(key)
		if field.IsValid() && field.CanSet() {
			fieldValue := reflect.ValueOf(value)

			// Handle time fields
			if field.Type() == reflect.TypeOf((*time.Time)(nil)) && fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				field.Set(reflect.ValueOf(value))
			} else if field.Type() == reflect.TypeOf(time.Time{}) && fieldValue.Type() == reflect.TypeOf("") {
				strValue, ok := value.(string)
				if !ok {
					return fmt.Errorf("expected string for field %s, got %T", key, value)
				}
				parsedTime, err := time.Parse(time.RFC3339, strValue)
				if err != nil {
					return fmt.Errorf("cannot parse string to time for field %s: %v", key, err)
				}
				location := time.FixedZone("UTC+3", 3*60*60)
				parsedTimeInZone := parsedTime.In(location)
				field.Set(reflect.ValueOf(parsedTimeInZone))
			} else if field.Type().Kind() == reflect.Ptr {
				// Handle pointer fields
				if field.Type().Elem().Kind() == reflect.Int64 {
					// Handle *int64 case
					if fieldValue.Kind() == reflect.Int {
						ptr := reflect.New(reflect.TypeOf(int64(0)))
						ptr.Elem().Set(reflect.ValueOf(int64(fieldValue.Int())))
						field.Set(ptr)
					} else if fieldValue.Kind() == reflect.Int64 {
						ptr := reflect.New(reflect.TypeOf(int64(0)))
						ptr.Elem().Set(fieldValue)
						field.Set(ptr)
					}
				} else if field.Type().Elem().Kind() == reflect.String {
					// Handle *string case
					if fieldValue.Kind() == reflect.String {
						ptr := reflect.New(reflect.TypeOf(""))
						ptr.Elem().Set(fieldValue)
						field.Set(ptr)
					}
				} else if field.Type().Elem().Kind() == reflect.Struct && field.Type().Elem() == reflect.TypeOf(time.Time{}) {
					// Handle *time.Time case
					if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
						ptr := reflect.New(reflect.TypeOf(time.Time{}))
						ptr.Elem().Set(fieldValue)
						field.Set(ptr)
					}
				} else {
					// Handle other cases
					field.Set(reflect.New(field.Type().Elem()))
					field.Elem().Set(fieldValue)
				}
			} else if fieldValue.Type().ConvertibleTo(field.Type()) {
				// Handle direct conversion
				field.Set(fieldValue.Convert(field.Type()))
			} else if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
				// Handle nil pointers
				field.Set(reflect.Zero(field.Type()))
			} else {
				return fmt.Errorf("cannot set field %s: incompatible type %s -> %s", key, fieldValue.Type(), field.Type())
			}
		}
	}
	return nil
}

func mapListToStruct(target interface{}, data []map[string]any) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Slice {
		// Added detailed log to check the type of target
		logger.Logger.Errorf("mapListToStruct: target is not a pointer to a slice, but %s", targetValue.Kind())
		return fmt.Errorf("target must be a pointer to a slice, but got %s", targetValue.Kind())
	}

	sliceElemType := targetValue.Elem().Type().Elem()

	// Log the data to be processed
	logger.Logger.Infof("Mapping data to slice of %s", sliceElemType)

	for _, item := range data {
		structInstance := reflect.New(sliceElemType).Elem()
		err := mapToStruct(structInstance.Addr().Interface(), item)
		if err != nil {
			// Log more details about the mapping failure
			logger.Logger.Errorf("Failed to map data for item: %v, error: %v", item, err)
			return fmt.Errorf("failed to map data to struct: %v", err)
		}
		targetValue.Elem().Set(reflect.Append(targetValue.Elem(), structInstance))
	}

	// Log successful mapping
	logger.Logger.Infof("Successfully mapped %d items to slice", len(data))
	return nil
}
