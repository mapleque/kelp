package str

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(src string) string {
	h := md5.New()
	h.Write([]byte(src))
	data := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return string(dst)
}
