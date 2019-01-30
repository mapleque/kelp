package str

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

func Sha1(src string) string {
	return Sha([]byte(src), sha1.New())
}

func Sha256(src string) string {
	return Sha([]byte(src), sha256.New())
}

func Sha512(src string) string {
	return Sha([]byte(src), sha512.New())
}

func Sha(src []byte, h hash.Hash) string {
	h.Write([]byte(src))
	data := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return string(dst)
}
