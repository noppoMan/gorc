package util

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func Base64Encode(raw_str string) string {
	encode_str := base64.StdEncoding.EncodeToString([]byte(raw_str))
	r := strings.NewReplacer("=", "-", "/", "_", "+", ".")
	y64_encode_str := r.Replace(encode_str)
	return y64_encode_str
}

func Base64Decode(y64_encode_str string) string {
	r := strings.NewReplacer("-", "=", "_", "/", ".", "+")
	b64_decode_str := r.Replace(y64_encode_str)
	data, err := base64.StdEncoding.DecodeString(b64_decode_str)
	if err != nil {
		fmt.Println("error: ", err)
		return ""
	}
	return string(data)
}
