package util

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

// BuildStruct validates the given fields against the struct T and,
// if successful, returns a new instance of T populated with those fields.
func BuildStruct[T any](input map[string]interface{}) (T, error) {
	var empty T
	var s T
	structVal := reflect.ValueOf(&s).Elem()
	structType := structVal.Type()
	usedKeys := make(map[string]bool)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.PkgPath != "" {
			// Unexported field; skip.
			continue
		}
		// Use the "yaml" tag if available.
		key := field.Name
		if tag, ok := field.Tag.Lookup("yaml"); ok {
			parts := strings.Split(tag, ",")
			if parts[0] != "" && parts[0] != "-" {
				key = parts[0]
			}
		}

		// If the key exists in the input, check convertibility and set the field.
		if val, ok := input[key]; ok {
			usedKeys[key] = true
			fieldVal := structVal.Field(i)
			vVal := reflect.ValueOf(val)
			// If the field is a map and the value is a map, use the helper.
			if field.Type.Kind() == reflect.Map && vVal.Kind() == reflect.Map {
				newMap, err := convertMapValue(field.Type, vVal)
				if err != nil {
					return empty, err
				}
				fieldVal.Set(newMap)
			} else {
				if !vVal.Type().ConvertibleTo(field.Type) {
					return empty, fmt.Errorf("cannot convert field '%s': of type %s to type '%s'", field.Name, vVal.Type(), field.Type)
				}
				fieldVal.Set(vVal.Convert(field.Type))
			}
		}
	}

	// Check for unexpected fields.
	for key := range input {
		if !usedKeys[key] {
			return s, fmt.Errorf("unexpected field: '%s'", key)
		}
	}

	return s, nil
}

// convertMapValue converts an input map (as a reflect.Value) into a new map of type targetType.
// It checks that each element is convertible and returns an error if any element fails.
func convertMapValue(targetType reflect.Type, inputMap reflect.Value) (reflect.Value, error) {
	newMap := reflect.MakeMap(targetType)
	for _, mapKey := range inputMap.MapKeys() {
		mapVal := inputMap.MapIndex(mapKey)
		// If the value is an interface, unwrap it.
		if mapVal.Kind() == reflect.Interface && !mapVal.IsNil() {
			mapVal = mapVal.Elem()
		}
		if !mapVal.Type().ConvertibleTo(targetType.Elem()) {
			return reflect.Value{}, fmt.Errorf("cannot convert map element for key '%v': %s to %s", mapKey, mapVal.Type(), targetType.Elem())
		}
		convertedVal := mapVal.Convert(targetType.Elem())
		newMap.SetMapIndex(mapKey, convertedVal)
	}
	return newMap, nil
}

// GenerateRandomID generates a random string of the specified size.
// The string consists of lowercase and uppercase letters and digits.
//
// Parameters:
//   - size: The length of the random string to generate.
//
// Returns:
//
//	A random string of the specified length.
func GenerateRandomID(size int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, int(size))
	rand.Seed(uint64(time.Now().UnixNano()))

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
