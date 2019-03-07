package tools

import (
	"encoding/base64"
)

func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func Base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}
