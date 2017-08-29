package util

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

func MD5(tar string) string {
	h := md5.New()
	h.Write([]byte(tar))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

func AesEcbEnc(key, tar string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Error("key error", err.Error())
		return ""
	}
	blockSize := block.BlockSize()
	src := PKCS5Padding([]byte(tar), blockSize)
	dst := make([]byte, 0)
	tmp := make([]byte, block.BlockSize())
	for len(src) > 0 {
		block.Encrypt(tmp, src[:blockSize])
		src = src[blockSize:]
		dst = append(dst, tmp...)
	}

	res := base64.URLEncoding.EncodeToString(dst)
	return res
}

func AesEcbDec(key, tar string) string {
	src, err := base64.URLEncoding.DecodeString(tar)
	if err != nil {
		return ""
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Error("key error", err.Error())
		return ""
	}
	blockSize := block.BlockSize()
	dst := make([]byte, 0)
	tmp := make([]byte, blockSize)

	if len(src) <= blockSize {
		return ""
	}

	for len(src) > 0 {
		block.Decrypt(tmp, src[:blockSize])
		src = src[blockSize:]
		dst = append(dst, tmp...)
	}

	res := PKCS5UnPadding(dst)

	return string(res)
}

func PKCS5Padding(tar []byte, blockSize int) []byte {
	padding := blockSize - len(tar)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(tar, padtext...)
}
func PKCS5UnPadding(tar []byte) []byte {
	length := len(tar)
	unpadding := int(tar[length-1])
	tarlen := length - unpadding
	if tarlen < 0 || tarlen > len(tar) {
		return []byte{}
	}
	return tar[:tarlen]
}
