package utils

import (
	"encoding/hex"
	"strings"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

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

func Bytes_deserialize(data []byte) ([]byte, error) {
	var hexString string
	if err := json.Unmarshal(data, &hexString); err != nil {
		return nil, err
	}

	return Hex_str_to_bytes(hexString)
}

func Bytes_serialize(bytes []byte) ([]byte, error) {
	if bytes == nil {
		return json.Marshal(nil)
	}
	hexString := hex.EncodeToString(bytes)
	return json.Marshal(hexString)
}
