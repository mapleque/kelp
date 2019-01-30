package str_test

import (
	"fmt"
	"github.com/mapleque/kelp/str"
)

const (
	encodingKey = "bVpxYlg0MzNNN05XQ3ZYY09nUGJGR1JtRENVbDh2a24="
	encodingIv  = "AAAAAAAAAAAAAAAAAAAAAA=="
)

func Example_aesEcb() {
	encodingKey := "bVpxYlg0MzNNN05XQ3ZYY09nUGJGR1JtRENVbDh2a24="
	data := "Hello kelp aes ecb!"

	// decode key
	key := str.Base64Decode(encodingKey)
	fmt.Printf("key: %s\n", key)

	// ecb encrypt
	encryptData, _ := str.AesEcbEncrypt(key, data)
	encodingEncryptData := str.Base64Encode(encryptData)
	fmt.Printf("encrypt data: %s\n", encodingEncryptData)

	// ecb decrypt
	decryptData := str.Base64Decode(encodingEncryptData)
	targetData, _ := str.AesEcbDecrypt(key, decryptData)
	fmt.Printf("target data: %s\n", targetData)

	// Output:
	// key: mZqbX433M7NWCvXcOgPbFGRmDCUl8vkn
	// encrypt data: pzkRXscGTq+YRUNG4wMkDMYIwnbvMeiU369bJxypWIs=
	// target data: Hello kelp aes ecb!
}

func Example_aesCbc() {
	encodingKey = "bVpxYlg0MzNNN05XQ3ZYY09nUGJGR1JtRENVbDh2a24="
	encodingIv = "AAAAAAAAAAAAAAAAAAAAAA=="
	data := "Hello kelp aes cbc!"

	// decode key
	key := str.Base64Decode(encodingKey)
	fmt.Printf("key: %s\n", key)

	// decode iv
	iv := str.Base64Decode(encodingIv)
	fmt.Printf("iv: %b\n", iv)

	// cbc encrypt
	encryptData, _ := str.AesCbcEncrypt(key, iv, data)
	encodingEncryptData := str.Base64Encode(encryptData)
	fmt.Printf("encrypt data: %s\n", encodingEncryptData)

	// cbc decrypt
	decryptData := str.Base64Decode(encodingEncryptData)
	targetData, _ := str.AesCbcDecrypt(key, iv, decryptData)
	fmt.Printf("target data: %s\n", targetData)

	// Output:
	// key: mZqbX433M7NWCvXcOgPbFGRmDCUl8vkn
	// iv: [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	// encrypt data: bJhrmxa94RDahqEtCpvPsrg9ZDlGMHYXHO66tXODhGM=
	// target data: Hello kelp aes cbc!
}
