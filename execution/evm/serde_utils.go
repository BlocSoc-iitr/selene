package evm

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

func unmarshalJSON(data []byte, v interface{}) error {
	var serialized map[string]interface{}
	if err := json.Unmarshal(data, &serialized); err != nil {
		return fmt.Errorf("error unmarshalling into map: %w", err)
	}

	// Get the value of the struct using reflection
	value := reflect.ValueOf(v).Elem()
	typeOfValue := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := typeOfValue.Field(i)
		jsonTag := fieldType.Tag.Get("json")

		// Remove ",omitempty" if present in json tag
		if idx := strings.Index(jsonTag, ",omitempty"); idx != -1 {
			jsonTag = jsonTag[:idx]
		}

		// Convert the struct field name or its json tag to snake_case for matching
		fieldName := jsonTag
		if fieldName == "" {
			fieldName = snakeCase(fieldType.Name)
		}

		

		if rawValue, ok := serialized[fieldName]; ok {
			
			if err := setFieldValue(field, rawValue); err != nil {
				return fmt.Errorf("error setting field %s: %w", fieldName, err)
			}
		} 
	}
	return nil
}

// snakeCase converts CamelCase strings to snake_case (e.g., "ParentRoot" to "parent_root").
func snakeCase(input string) string {
	var result []rune
	for i, r := range input {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func setFieldValue(field reflect.Value, value interface{}) error {

	if field.Kind() == reflect.Ptr {
		// If fie
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}

		// Handle *big.Int specifically
		if field.Type() == reflect.TypeOf(&big.Int{}) {
			bi := new(big.Int)
			strVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("expected string value for big.Int")
			}

			var success bool
			if strings.HasPrefix(strVal, "0x") {
				strVal = strVal[2:]
				bi, success = bi.SetString(strVal, 16)
			} else {
				bi, success = bi.SetString(strVal, 10)
			}

			if !success {
				return fmt.Errorf("invalid string format for big.Int: %s", strVal)
			}

			field.Set(reflect.ValueOf(bi))
			return nil
		}

		// Handle *uint64 specifically
		if field.Type().Elem().Kind() == reflect.Uint64 {
			switch v := value.(type) {
			case string:
				var val uint64
				var err error
				if strings.HasPrefix(v, "0x") {
					val, err = strconv.ParseUint(v[2:], 16, 64)
				} else {
					val, err = strconv.ParseUint(v, 10, 64)
				}
				if err != nil {
					return fmt.Errorf("failed to convert to uint64: %w", err)
				}
				ptr := new(uint64)
				*ptr = val
				field.Set(reflect.ValueOf(ptr))
				return nil
			case float64:
				ptr := new(uint64)
				*ptr = uint64(v)
				field.Set(reflect.ValueOf(ptr))
				return nil
			}
			return fmt.Errorf("unsupported type for *uint64")
		}

		// For other pointer types, recurse on the element
		return setFieldValue(field.Elem(), value)
	}

	switch field.Kind() {
	case reflect.Uint64, reflect.Uint16, reflect.Uint8:
		var val uint64
		var err error
		strVal := value.(string)

		// Handle hex string
		if strings.HasPrefix(strVal, "0x") {
			val, err = strconv.ParseUint(strVal[2:], 16, 64)
		} else {
			val, err = strconv.ParseUint(strVal, 10, 64)
		}

		if err != nil {
			return fmt.Errorf("failed to convert to uint: %w", err)
		}

		// Check bounds based on field type
		switch field.Kind() {
		case reflect.Uint16:
			if val > math.MaxUint16 {
				return fmt.Errorf("value exceeds uint16 range")
			}
		case reflect.Uint8:
			if val > math.MaxUint8 {
				return fmt.Errorf("value exceeds uint8 range")
			}
		}

		field.SetUint(val)

	case reflect.Int:
		strVal := value.(string)
		var val int64
		var err error

		if strings.HasPrefix(strVal, "0x") {
			val, err = strconv.ParseInt(strVal[2:], 16, 64)
		} else {
			val, err = strconv.ParseInt(strVal, 10, 64)
		}

		if err != nil {
			return fmt.Errorf("failed to convert to int: %w", err)
		}
		field.SetInt(val)

	case reflect.Struct:
		switch field.Interface().(type) {
		case Address:
			strVal, ok := value.(string)
			if !ok || !strings.HasPrefix(strVal, "0x") {
				return fmt.Errorf("invalid address format")
			}

			decoded, err := hex.DecodeString(strVal[2:])
			if err != nil {
				return fmt.Errorf("failed to decode address: %w", err)
			}

			var addr [20]byte
			copy(addr[:], decoded)
			field.Set(reflect.ValueOf(Address{Addr: addr}))

		default:
			rawJSON, err := json.Marshal(value)
			if err != nil {
				return errors.New("error marshalling struct")
			}
			return json.Unmarshal(rawJSON, field.Addr().Interface())
		}

	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.Uint8 {
			strVal, ok := value.(string)
			if !ok || !strings.HasPrefix(strVal, "0x") {
				return fmt.Errorf("invalid hex string format")
			}

			decoded, err := hex.DecodeString(strVal[2:])
			if err != nil {
				return fmt.Errorf("failed to decode hex string: %w", err)
			}

			field.SetBytes(decoded)
			return nil
		}

		if field.Type().Elem().Kind() == reflect.Uint16 {
			sliceValue := reflect.ValueOf(value)
			newSlice := reflect.MakeSlice(field.Type(), sliceValue.Len(), sliceValue.Len())
			for i := 0; i < sliceValue.Len(); i++ {
				strVal := sliceValue.Index(i).Interface().(string)
				var val uint64
				var err error

				if strings.HasPrefix(strVal, "0x") {
					val, err = strconv.ParseUint(strVal[2:], 16, 16)
				} else {
					val, err = strconv.ParseUint(strVal, 10, 16)
				}

				if err != nil {
					return fmt.Errorf("failed to parse uint16 at index %d: %w", i, err)
				}
				newSlice.Index(i).SetUint(val)
			}
			field.Set(newSlice)
			return nil
		}

		// Handle slices of slices ([][]byte)
		if field.Type().Elem().Kind() == reflect.Slice {
			sliceValue := reflect.ValueOf(value)
			if sliceValue.Kind() != reflect.Slice {
				return fmt.Errorf("expected slice for 2D slice input")
			}

			newSlice := reflect.MakeSlice(field.Type(), sliceValue.Len(), sliceValue.Len())

			for i := 0; i < sliceValue.Len(); i++ {
				err := setFieldValue(newSlice.Index(i), sliceValue.Index(i).Interface())
				if err != nil {
					return fmt.Errorf("error setting slice element %d: %w", i, err)
				}
			}

			field.Set(newSlice)
			return nil
		}

		// For other slice types
		rawJSON, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("error marshalling slice: %w", err)
		}
		return json.Unmarshal(rawJSON, field.Addr().Interface())

	case reflect.Bool:
		switch v := value.(type) {
		case string:
			boolVal, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("failed to convert to bool: %w", err)
			}
			field.SetBool(boolVal)
		case bool:
			field.SetBool(v)
		default:
			return fmt.Errorf("unsupported type for bool conversion")
		}

	default:
		rawJSON, err := json.Marshal(value)
		if err != nil {
			return errors.New("error marshalling value")
		}
		return json.Unmarshal(rawJSON, field.Addr().Interface())
	}
	return nil
}



func marshalJSON(v interface{}) ([]byte, error) {
	serialized := make(map[string]interface{})

	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	typeOfValue := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := typeOfValue.Field(i)
		jsonTag := fieldType.Tag.Get("json")

		fieldName := jsonTag
		if fieldName == "" {
			fieldName = snakeCase(fieldType.Name)
		}

		val, err := getFieldValue(field)
		if err != nil {
			return nil, fmt.Errorf("error getting field %s value: %w", fieldName, err)
		}

		// Skip nil values
		if val == nil {
			continue
		}

		serialized[fieldName] = val
	}

	// Marshal the map to JSON
	return json.Marshal(serialized)
}

func getFieldValue(field reflect.Value) (interface{}, error) {
	if field.Kind() == reflect.Ptr && field.IsNil() {
		return nil, nil
	}

	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}

	switch field.Kind() {
	case reflect.Uint64, reflect.Uint16, reflect.Uint8:
		return fmt.Sprintf("0x%x", field.Uint()), nil

	case reflect.Int:
		return fmt.Sprintf("0x%x", field.Int()), nil

	case reflect.Struct:
		switch field.Interface().(type) {
		case Address:
			addr := field.Interface().(Address)
			return fmt.Sprintf("0x%s", hex.EncodeToString(addr.Addr[:])), nil
		}

		return marshalJSON(field.Interface())

	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.Uint8 {
			bytes := field.Bytes()
			if len(bytes) == 0 {
				return "0x", nil
			}
			return fmt.Sprintf("0x%s", hex.EncodeToString(bytes)), nil
		}

		sliceLen := field.Len()
		result := make([]interface{}, sliceLen)
		for i := 0; i < sliceLen; i++ {
			val, err := getFieldValue(field.Index(i))
			if err != nil {
				return nil, fmt.Errorf("error getting slice element %d: %w", i, err)
			}
			result[i] = val
		}
		return result, nil

	case reflect.Bool:
		return field.Bool(), nil

	default:
		return field.Interface(), nil
	}
}
