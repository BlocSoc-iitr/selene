package config

import (
	"encoding/hex"
	"encoding/json"

	"github.com/BlocSoc-iitr/selene/common"
)

func bytes_serialise(bytes []byte) ([]byte, error) {

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

func bytes_deserialise(data []byte) ([]byte, error) {
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
