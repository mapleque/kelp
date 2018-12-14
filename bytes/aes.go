package bytes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func AesEcbEncrypt(key, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	data := _PKCS5Padding(src, blockSize)
	dst := make([]byte, 0)
	tmp := make([]byte, block.BlockSize())
	for len(data) > 0 {
		block.Encrypt(tmp, data[:blockSize])
		data = data[blockSize:]
		dst = append(dst, tmp...)
	}

	return dst, nil
}

func AesEcbDecrypt(key, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	dst := make([]byte, 0)
	tmp := make([]byte, blockSize)

	data := _PKCS5Padding(src, blockSize)

	for len(data) > 0 {
		block.Decrypt(tmp, data[:blockSize])
		data = data[blockSize:]
		dst = append(dst, tmp...)
	}

	res := _PKCS5UnPadding(dst)

	return res, nil
}

func AesCbcEncrypt(key, iv, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)

	blockSize := block.BlockSize()
	data := _PKCS5Padding(src, blockSize)

	dst := make([]byte, len(data))
	blockMode.CryptBlocks(dst, data)

	return dst, nil
}

func AesCbcDecrypt(key, iv, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv = _PKCS5Padding(iv, block.BlockSize())
	blockMode := cipher.NewCBCDecrypter(block, iv)

	src = _PKCS5Padding(src, block.BlockSize())

	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)

	res := _PKCS5UnPadding(dst)
	return res, nil
}

func _PKCS5Padding(data []byte, blockSize int) []byte {
	paddingNumber := len(data) % blockSize
	if paddingNumber == 0 {
		return data
	}
	padding := blockSize - paddingNumber
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func _PKCS5UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	datalen := length - unpadding
	if datalen < 0 || datalen > len(data) {
		return []byte{}
	}
	return data[:datalen]
}
