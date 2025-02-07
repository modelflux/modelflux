package util

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

// ValidateStructFields checks:
//   - The "yaml" tags (if present) match keys in the map
//   - Each value in 'fields' is convertible to the corresponding field's type
//   - No unexpected fields remain in 'fields'
func ValidateStructFields[T any](fields map[string]interface{}) error {
	var empty T
	structVal := reflect.ValueOf(&empty).Elem()
	structType := structVal.Type()

	usedKeys := make(map[string]bool)

	// Check each exported field in T.
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

		// Check if the map has this key; if so, ensure convertibility.
		if val, ok := fields[key]; ok {
			usedKeys[key] = true
			vVal := reflect.ValueOf(val)
			if !vVal.Type().ConvertibleTo(field.Type) {
				return fmt.Errorf("cannot convert field '%s' to type '%s'",
					field.Name, field.Type)
			}
		}
	}

	// Check for unexpected fields in the map.
	for key := range fields {
		if !usedKeys[key] {
			return fmt.Errorf("unexpected field: '%s'", key)
		}
	}

	return nil
}

// CreateStruct creates a struct of type T from 'fields' after validating them.
func CreateStruct[T any](fields map[string]interface{}) (T, error) {
	var s T
	structVal := reflect.ValueOf(&s).Elem()
	structType := structVal.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
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
			fieldVal := structVal.Field(i)
			vVal := reflect.ValueOf(val)
			fieldVal.Set(vVal.Convert(fieldVal.Type()))
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
