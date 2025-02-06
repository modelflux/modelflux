package util

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

// CreateStruct creates an instance of a struct of type T and populates its fields
// based on the provided map of field names to values. The function uses reflection
// to set the struct fields and supports struct tags for field names.
//
// Type Parameters:
//   - T: The type of the struct to be created.
//
// Parameters:
//   - fields: A map where the keys are field names (or their "yaml" tag values) and
//     the values are the values to be set for those fields.
//
// Returns:
//   - T: An instance of the struct with the populated fields.
//   - error: An error if there is an issue with setting the fields, such as a type
//     mismatch or an unexpected field in the map.
//
// Notes:
//   - Unexported fields in the struct are skipped.
//   - If a field has a "yaml" tag, the tag value is used as the key in the map.
//   - If a field value in the map cannot be converted to the field's type, an error
//     is returned.
//   - If the map contains keys that do not correspond to any struct fields, an error
//     is returned.
func CreateStruct[T any](fields map[string]interface{}) (T, error) {
	var s T

	structVal := reflect.ValueOf(&s).Elem()
	structType := structVal.Type()
	usedKeys := make(map[string]bool)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		// Skip unexported fields.
		if field.PkgPath != "" {
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

		if val, ok := fields[key]; ok {
			usedKeys[key] = true
			fieldVal := structVal.Field(i)
			if !fieldVal.CanSet() {
				continue
			}
			vVal := reflect.ValueOf(val)
			if !vVal.Type().ConvertibleTo(fieldVal.Type()) {
				return s, fmt.Errorf("cannot convert field %s to type %s", field.Name, fieldVal.Type())
			}
			fieldVal.Set(vVal.Convert(fieldVal.Type()))
		}
	}

	// Optionally, check for unexpected fields.
	for key := range fields {
		if !usedKeys[key] {
			return s, fmt.Errorf("unexpected field: %s", key)
		}
	}

	return s, nil
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
