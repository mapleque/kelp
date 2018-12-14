package bytes

import (
	"encoding/base64"
)

func Base64Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.URLEncoding.Encode(dst, src)
	return dst
}

func Base64Decode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	base64.URLEncoding.Decode(dst, src)
	for dst[len(dst)-1] == 0 {
		dst = dst[:len(dst)-1]
	}
	return dst
}

func Base64StdEncode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func Base64StdDecode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	base64.StdEncoding.Decode(dst, src)
	for dst[len(dst)-1] == 0 {
		dst = dst[:len(dst)-1]
	}
	return dst
}
