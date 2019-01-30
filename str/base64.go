package str

import (
	"encoding/base64"
)

func Base64UrlEncode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func Base64UrlDecode(src string) []byte {
	dst, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		log.Error(err)
		return nil
	}
	return dst
}

func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func Base64Decode(src string) []byte {
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		log.Error(err)
		return nil
	}
	return dst
}
