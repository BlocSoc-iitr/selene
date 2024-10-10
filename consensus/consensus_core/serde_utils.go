package consensus_core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func (f *FinalityUpdate) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, f)
}

func (o *OptimisticUpdate) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, o)
}

func (b *Bootstrap) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, b)
}



func (f *Forks) UnmarshalJSON(data []byte) error {
	var serialized map[string]interface{}
	if err := json.Unmarshal(data, &serialized); err != nil {
		return fmt.Errorf("error unmarshalling into map: %w", err)
	}

	v := reflect.ValueOf(f).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Tag.Get("json")
		if fieldName == "" {
			fieldName = snakeCase(fieldType.Name)
		}

		if forkData, ok := serialized[fieldName]; ok {
			forkJSON, err := json.Marshal(forkData)
			if err != nil {
				return fmt.Errorf("error marshalling %s: %w", fieldName, err)
			}
			if err := json.Unmarshal(forkJSON, field.Addr().Interface()); err != nil {
				return fmt.Errorf("error unmarshalling %s: %w", fieldName, err)
			}
		}
	}

	return nil
}

func (h *Header) UnmarshalJSON(data []byte) error {
	// Define a temporary map to hold JSON data
	var serialized map[string]interface{}
	if err := json.Unmarshal(data, &serialized); err != nil {
		return fmt.Errorf("error unmarshalling into map: %w", err)
	}

	// Extract the "beacon" field from the serialized map
	beaconData, ok := serialized["beacon"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("beacon field missing or invalid")
	}

	v := reflect.ValueOf(h).Elem()
	t := v.Type()

	// Iterate through each field of the Header struct using reflection
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		jsonTag := fieldType.Tag.Get("json")

		// Convert the struct field name or its json tag to snake_case for matching
		fieldName := jsonTag
		if fieldName == "" {
			fieldName = snakeCase(fieldType.Name)
		}

		if value, ok := beaconData[fieldName]; ok {
			switch field.Interface().(type) {
			case uint64:
				val, err := Str_to_uint64(value.(string))
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(val))
			case Bytes32:
				val, err := Hex_str_to_Bytes32(value.(string))
				if err != nil {
					continue
				}
				field.Set(reflect.ValueOf(val))
			}
		} else {
			continue
		}
	}

	return nil
}

func (s *SyncAggregate) UnmarshalJSON(data []byte) error {
	// Define a temporary map to hold JSON data
	var serialized map[string]string
	if err := json.Unmarshal(data, &serialized); err != nil {
		return fmt.Errorf("error unmarshalling into map: %w", err)
	}

	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		snakeFieldName := snakeCase(fieldName)

		switch snakeFieldName {
		case "sync_committee_bits":
			if bits, ok := serialized[snakeFieldName]; ok {
				decodedBits := hexDecode(bits)
				if decodedBits == nil {
					return fmt.Errorf("error decoding sync_committee_bits")
				}
				field.SetBytes(decodedBits)
			} else {
				return fmt.Errorf("%s field missing or invalid", snakeFieldName)
			}
		case "sync_committee_signature":
			if signature, ok := serialized[snakeFieldName]; ok {
				decodedSignature := hexDecode(signature)
				if decodedSignature == nil {
					return fmt.Errorf("error decoding sync_committee_signature")
				}
				if len(decodedSignature) != len(s.SyncCommitteeSignature) {
					return fmt.Errorf("decoded sync_committee_signature length mismatch")
				}
				copy(s.SyncCommitteeSignature[:], decodedSignature)
			} else {
				return fmt.Errorf("%s field missing or invalid", snakeFieldName)
			}
		default:
			return fmt.Errorf("unsupported field %s", fieldName)
		}
	}

	return nil
}

func (s *SyncCommittee) UnmarshalJSON(data []byte) error {
	var serialized map[string]interface{}
	if err := json.Unmarshal(data, &serialized); err != nil {
		return fmt.Errorf("error unmarshalling into map: %w", err)
	}

	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		snakeFieldName := snakeCase(fieldName)

		switch snakeFieldName {
		case "pubkeys":
			if pubkeys, ok := serialized[snakeFieldName].([]interface{}); ok {
				s.Pubkeys = make([]BLSPubKey, len(pubkeys))
				for j, pubkey := range pubkeys {
					if pubkeyStr, ok := pubkey.(string); ok {
						decoded := hexDecode(pubkeyStr)
						if decoded == nil {
							return fmt.Errorf("error decoding pubkey at index %d", j)
						}
						s.Pubkeys[j] = BLSPubKey(decoded)
					} else {
						return fmt.Errorf("invalid pubkey format at index %d", j)
					}
				}
			} else {
				return fmt.Errorf("%s field missing or invalid", snakeFieldName)
			}
		case "aggregate_pubkey":
			if aggPubkey, ok := serialized[snakeFieldName].(string); ok {
				decoded := hexDecode(aggPubkey)
				if decoded == nil {
					return fmt.Errorf("error decoding aggregate pubkey")
				}
				s.AggregatePubkey = BLSPubKey(decoded)
			} else {
				return fmt.Errorf("%s field missing or invalid", snakeFieldName)
			}
		default:
			return fmt.Errorf("unsupported field %s", fieldName)
		}
	}

	return nil
}

func (u *Update) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, u)
}

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

func (b *BeaconBlock) UnmarshalJSON(data []byte) error {
	// Define a temporary structure to handle both string and number slots.
	var tmp struct {
		Slot json.RawMessage `json:"slot"`
	}

	// Unmarshal into the temporary struct.
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("error unmarshalling into temporary struct: %w", err)
	}

	// If it wasn't a string, try unmarshalling as a number.
	var slotNum uint64
	if err := json.Unmarshal(tmp.Slot, &slotNum); err == nil {
		b.Slot = slotNum
		return nil
	}

	return fmt.Errorf("slot field is not a valid string or number")
}

func setFieldValue(field reflect.Value, value interface{}) error {
	switch field.Kind() {
	case reflect.Struct:
		// Marshal and unmarshal for struct types
		rawJSON, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return json.Unmarshal(rawJSON, field.Addr().Interface())

	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.Uint8 {
			// For Bytes32 conversion
			dataSlice := value.([]interface{})
			for _, b := range dataSlice {
				field.Set(reflect.Append(field, reflect.ValueOf(Bytes32(hexDecode(b.(string)))))) // Adjust hexDecode as needed
			}
			return nil
		}
		// For other slice types, marshal and unmarshal
		rawJSON, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return json.Unmarshal(rawJSON, field.Addr().Interface())

	case reflect.String:
		field.SetString(value.(string))
	case reflect.Uint64:
		val, err := Str_to_uint64(value.(string))
		if err != nil {
			return err
		}
		field.SetUint(val)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Type())
	}
	return nil
}




// if we need to export the functions , just make their first letter capitalised
func Hex_str_to_bytes(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")

	bytesArray, err := hex.DecodeString(s)

	if err != nil {
		return nil, err
	}

	return bytesArray, nil
}

func Hex_str_to_Bytes32(s string) (Bytes32, error) {
	bytesArray, err := Hex_str_to_bytes(s)
	if err != nil {
		return Bytes32{}, err
	}

	var bytes32 Bytes32
	copy(bytes32[:], bytesArray)
	return bytes32, nil
}

func Address_to_hex_string(addr common.Address) string {
	bytesArray := addr.Bytes()
	return fmt.Sprintf("0x%x", hex.EncodeToString(bytesArray))
}

func U64_to_hex_string(val uint64) string {
	return fmt.Sprintf("0x%x", val)
}

func Str_to_uint64(s string) (uint64, error) {
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
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

// hexDecode is a utility function to decode hex strings to byte slices.
func hexDecode(input string) []byte {
	data, _ := Hex_str_to_bytes(input)
	return data
}

func (b *Bytes32) UnmarshalJSON(data []byte) error {
	var hexString string
	if err := json.Unmarshal(data, &hexString); err != nil {
		return err
	}

	copy(b[:], hexDecode(hexString)[:])

	return nil
}
