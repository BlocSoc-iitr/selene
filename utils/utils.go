package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BlocSoc-iitr/selene/common"
)

// if we need to export the functions , just make their first letter capitalised
func Hex_str_to_bytes(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")

	bytesArray, err := hex.DecodeString(s)

	if err != nil {
		return nil, err
	}

	return bytesArray, nil
}

func Address_to_hex_string(addr common.Address) string {
	bytesArray := addr.Bytes()
	return fmt.Sprintf("0x%x", hex.EncodeToString(bytesArray))
}

func U64_to_hex_string(val uint64) string {
	return fmt.Sprintf("0x%x", val)
}

func BytesSerialise(bytes []byte) ([]byte, error) {

	if bytes == nil {
		return json.Marshal(nil)
	}
	bytesString := hex.EncodeToString(bytes)
	result, err := json.Marshal(bytesString)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func BytesDeserialise(data []byte) ([]byte, error) {
	var bytesOpt *string
	if err := json.Unmarshal(data, &bytesOpt); err != nil {
		return nil, err
	}

	if bytesOpt == nil {
		return nil, nil
	} else {
		bytes, err := common.Hex_str_to_bytes(*bytesOpt)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	}

}
