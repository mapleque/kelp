package bytes

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

func Md5(src []byte) []byte {
	h := md5.New()
	h.Write(src)
	data := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return []byte(dst)
}

func RandMd5() []byte {
	timestamp := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	prefix := []byte(strconv.Itoa(rand.Intn(10000)))
	surfix := []byte(strconv.Itoa(rand.Intn(10000)))
	return Md5(bytes.Join([][]byte{prefix, timestamp, surfix}, []byte("")))
}
