package bytes

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha1(src []byte) []byte {
	h := sha1.New()
	h.Write(src)
	data := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return dst
}
