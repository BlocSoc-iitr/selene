package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// if we need to export the functions , just make their first letter capitalised
func Hex_str_to_bytes(s string) ([]byte){
	bytesArray, _ := hex.DecodeString(s[2:])
	return bytesArray
}

func Address_to_hex_string(addr common.Address) (string){
	bytesArray := addr.Bytes()
	return fmt.Sprintf("0x%x", hex.EncodeToString(bytesArray))
}

func U64_to_hex_string(val uint64) (string){
	return fmt.Sprintf("0x%x", val)
}

func Bytes_deserialize(data []byte) ([]byte, error) {
    var hexString string
    if err := json.Unmarshal(data, &hexString); err != nil {
        return nil, err
    }
    
    bytes, err := hex.DecodeString(hexString)
    if err != nil {
        return nil, err
    }

    return bytes, nil
}

func Bytes_serialize(bytes []byte) ([]byte, error) {
    hexString := hex.EncodeToString(bytes)
    return json.Marshal(hexString)
}