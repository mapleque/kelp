package str

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func AesEcbEncrypt(key []byte, src string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	data := PKCS5Padding([]byte(src), blockSize)
	dst := make([]byte, 0)
	tmp := make([]byte, block.BlockSize())
	for len(data) > 0 {
		block.Encrypt(tmp, data[:blockSize])
		data = data[blockSize:]
		dst = append(dst, tmp...)
	}

	return dst, nil
}

func AesEcbDecrypt(key, src []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	dst := make([]byte, 0)
	tmp := make([]byte, blockSize)

	data := src
	if len(data)%blockSize != 0 {
		return "", errors.New("[kelp.str] aes ecb descrypt faild: data not full block")
	}

	for len(data) > 0 {
		block.Decrypt(tmp, data[:blockSize])
		data = data[blockSize:]
		dst = append(dst, tmp...)
	}

	res := PKCS5UnPadding(dst)
	return string(res), nil
}

func AesCbcEncrypt(key, iv []byte, src string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)

	blockSize := block.BlockSize()
	data := PKCS5Padding([]byte(src), blockSize)

	dst := make([]byte, len(data))
	blockMode.CryptBlocks(dst, data)

	return dst, nil
}

func AesCbcDecrypt(key, iv, src []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)

	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)

	res := PKCS5UnPadding(dst)
	return string(res), nil
}

func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func PKCS5UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	datalen := length - unpadding
	if datalen < 0 || datalen > len(data) {
		return nil
	}
	return data[:datalen]
}
