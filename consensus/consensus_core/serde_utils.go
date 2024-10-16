package consensus_core

import (
	"encoding/hex"
	"encoding/json"
	"errors"
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
	b.Hash = data
	return unmarshalJSON(data, b)
}

func (b *BeaconBlockBody) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, b)
}

func (e *Eth1Data) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, e)
}

func (a *Attestation) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, a)
}

func (a *AttesterSlashing) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, a)
}

func (i *IndexedAttestation) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, i)
}

func (p *ProposerSlashing) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, p)
}

func (s *SignedBeaconBlockHeader) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, s)
}

func (c *Checkpoint) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, c)
}

func (b *Bitlist) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, b)
}

func (d *Deposit) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, d)
}

func (d *DepositData) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, d)
}

func (d *SignedVoluntaryExit) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, d)
}

func (v *VoluntaryExit) UnmarshalJSON(data []byte) error {

	return unmarshalJSON(data, v)
}

func (e *ExecutionPayload) UnmarshalJSON(data []byte) error {

	return unmarshalJSON(data, e)
}

func (s *SignedBlsToExecutionChange) UnmarshalJSON(data []byte) error {

	return unmarshalJSON(data, s)
}

func (a *AttestationData) UnmarshalJSON(data []byte) error {

	return unmarshalJSON(data, a)
}

func (a *Address) UnmarshalJSON(data []byte) error {
	address, err := Hex_str_to_bytes(string(data))
	if err != nil {
		return err
	}
	*a = Address(address[:])
	return nil
}

func (w *Withdrawal) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, w)

}

func (b *BlsToExecutionChange) UnmarshalJSON(data []byte) error {
	return unmarshalJSON(data, b)
}

func setFieldValue(field reflect.Value, value interface{}) error {
	if field.Kind() == reflect.Ptr {
		field.Set(reflect.New(field.Type().Elem()))
		field = field.Elem()
	}

	switch field.Kind() {
	case reflect.Struct:
		rawJSON, err := json.Marshal(value)
		if err != nil {
			return errors.New("error marshalling struct")
		}
		return json.Unmarshal(rawJSON, field.Addr().Interface())

	case reflect.Slice:
		// Handling 2D slices (e.g., [][]byte)
		if field.Type().Elem().Kind() == reflect.Slice {
			sliceValue := reflect.ValueOf(value)
			if sliceValue.Kind() != reflect.Slice {
				return fmt.Errorf("expected slice for 2D slice input, got %T", value)
			}

			// Ensure the length matches
			if sliceValue.Len() != field.Len() {
				field.Set(reflect.MakeSlice(field.Type(), sliceValue.Len(), sliceValue.Len()))
			}

			// Iterate over each element and set values
			for i := 0; i < sliceValue.Len(); i++ {
				innerValue := sliceValue.Index(i).Interface()
				err := setFieldValue(field.Index(i), innerValue)
				if err != nil {
					return fmt.Errorf("error setting value for element %d: %w", i, err)
				}
			}
			return nil
		}

		// For slices of bytes (e.g., []byte)
		if field.Type().Elem().Kind() == reflect.Uint8 {
			dataSlice := value.(string)
			bytesArray, err := Hex_str_to_bytes(dataSlice)
			if err != nil {
				return errors.New("error decoding byte slice")
			}
			field.SetBytes(bytesArray)
			return nil
		}

		if field.Type().Elem().Kind() == reflect.Bool {
			bitlist := Str_to_boolArray(value.(string))
			field.Set(reflect.ValueOf(bitlist))
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
		// Handle both string and numeric input
		switch v := value.(type) {
		case string:
			val, err := Str_to_uint64(v)
			if err != nil {
				return errors.New("error decoding uint64 from string")
			}
			field.SetUint(val)
		case float64: // JSON unmarshalling often converts numbers to float64
			field.SetUint(uint64(v))
		default:
			return fmt.Errorf("expected string or float64 for uint64, got %T", v)
		}
	case reflect.Array:
		elemKind := field.Type().Elem().Kind()

		// Handle 2D arrays (e.g., [[32]byte])
		if elemKind == reflect.Array {
			sliceValue := reflect.ValueOf(value)
			if sliceValue.Kind() != reflect.Slice {
				return fmt.Errorf("expected slice for 2D array input, got %T", value)
			}

			if sliceValue.Len() != field.Len() {
				return fmt.Errorf("input slice length %d does not match the required length %d", sliceValue.Len(), field.Len())
			}

			for i := 0; i < field.Len(); i++ {
				innerValue := sliceValue.Index(i).Interface()
				err := setFieldValue(field.Index(i), innerValue)
				if err != nil {
					return fmt.Errorf("error setting value for element %d: %w", i, err)
				}
			}
			return nil
		}

		// Handling 1D arrays (e.g., [32]byte)
		if field.Type().Len() == 32 {
			val, err := Hex_str_to_Bytes32(value.(string))
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(val))
		} else if field.Type().Len() == 48 {
			val, err := Hex_str_to_bytes(value.(string))
			if err != nil {
				return err
			}
			if len(val) != 48 {
				return errors.New("input byte slice length does not match the required length of 48")
			}
			var fixedArray [48]byte
			copy(fixedArray[:], val)
			field.Set(reflect.ValueOf(fixedArray))
		} else if field.Type().Len() == 96 {
			val, err := Hex_str_to_bytes(value.(string))
			if err != nil {
				return err
			}
			if len(val) != 96 {
				return errors.New("input byte slice length does not match the required length of 96")
			}
			var fixedArray SignatureBytes
			copy(fixedArray[:], val)
			field.Set(reflect.ValueOf(fixedArray))
		} else if field.Type().Len() == 256 {
			val, err := Hex_str_to_bytes(value.(string))
			if err != nil {
				return err
			}
			if len(val) != 256 {
				return errors.New("input byte slice length does not match the required length of 256")
			}
			var fixedArray LogsBloom
			copy(fixedArray[:], val)
			field.Set(reflect.ValueOf(fixedArray))
		} else if field.Type().Len() == 1073741824 {
			val, err := Hex_str_to_bytes(value.(string))
			if err != nil {
				return err
			}
			if len(val) != 1073741824 {
				return errors.New("input byte slice length does not match the required length of 1073741824")
			}
			var fixedArray Transaction
			copy(fixedArray[:], val)
			field.Set(reflect.ValueOf(fixedArray))
		} else if field.Type().Len() == 20 {
			val, err := Hex_str_to_bytes(value.(string))
			if err != nil {
				return errors.New("error decoding address")
			}
			var fixedArray Address
			copy(fixedArray[:], val)
			field.Set(reflect.ValueOf(fixedArray))
		} else {
			return errors.New("unsupported array length")

		}
	}
	return nil
}

// if we need to export the functions , just make their first letter capitalised
func Hex_str_to_bytes(s string) ([]byte, error) {
	s = stripQuotes(s)
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

func Str_to_boolArray(s string) Bitlist {

	data, _ := Hex_str_to_bytes(s)
	var bitlist Bitlist
	for _, b := range data {
		if b == 0x01 {
			bitlist = append(bitlist, true)
		}
	}
	return bitlist
}

func stripQuotes(s string) string {
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
