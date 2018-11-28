package http

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	mrand "math/rand"
	"strconv"
	"time"
)

// 非对称加密，用于双方传输消息，一方加密另一方解密的场景

func RsaEncrypt(publicKey, src []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("illigal public key")
	}
	inter, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := inter.(*rsa.PublicKey)
	dst, err := rsa.EncryptPKCS1v15(rand.Reader, pub, src)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func RsaDecrypt(privateKey, src []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("illigal private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, src)
}

// 对称加密，用于自己加密自己解密的场景
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

	data := src
	if len(data)%blockSize != 0 {
		return nil, errors.New("web/crypto.AesEcbDecrypt: data not full block")
	}

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
	blockMode := cipher.NewCBCDecrypter(block, iv)

	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)

	res := _PKCS5UnPadding(dst)
	return res, nil
}

func _PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
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

// 信息摘要加密，用于只加密不需要解密的场景

func Md5(src []byte) []byte {
	h := md5.New()
	h.Write(src)
	data := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return []byte(dst)
}

func RandMd5() string {
	timestamp := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	prefix := []byte(strconv.Itoa(mrand.Intn(10000)))
	surfix := []byte(strconv.Itoa(mrand.Intn(10000)))
	token := string(Md5(bytes.Join([][]byte{prefix, timestamp, surfix}, []byte(""))))
	return token
}

func Sha1(src []byte) []byte {
	h := sha1.New()
	h.Write(src)
	data := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return dst
}

// 编解码，用于需要转换字符集的场景

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

// 签名，密钥+数据+时间戳签名，用于信息接收方校验身份和数据是否可信

func Sha1Sign(key, data []byte) []byte {
	body := Base64Encode(data)
	timestamp := []byte(strconv.FormatInt(time.Now().Unix(), 10))

	// body|timestamp|key ----> sha1
	sign := Sha1(bytes.Join([][]byte{body, timestamp, key}, []byte(`|`)))
	return sign
}

func Sha1SignTimestamp(key, data []byte, timestamp int64) []byte {
	body := Base64Encode(data)
	stampByte := []byte(strconv.FormatInt(timestamp, 10))

	// body|timestamp|key ----> sha1
	sign := Sha1(bytes.Join([][]byte{body, stampByte, key}, []byte(`|`)))
	return sign
}

func Sha1Verify(key, data, sign []byte, maxDelaySecond int) bool {
	if maxDelaySecond < 0 {
		return false
	}
	now := time.Now().Unix()
	body := Base64Encode(data)
	for i := 0; i <= maxDelaySecond; i++ {
		timestamp := []byte(strconv.FormatInt(now-int64(i), 10))
		if bytes.Equal(sign, Sha1(bytes.Join([][]byte{body, timestamp, key}, []byte(`|`)))) {
			return true
		}
	}
	return false
}

func Sha1VerifyTimestamp(key, data, sign []byte, maxDelaySecond, timestamp int64) bool {
	if timestamp+maxDelaySecond < time.Now().Unix() {
		return false
	}
	body := Base64Encode(data)
	stampByte := []byte(strconv.FormatInt(timestamp, 10))
	if bytes.Equal(sign, Sha1(bytes.Join([][]byte{body, stampByte, key}, []byte(`|`)))) {
		return true
	}
	return false
}
