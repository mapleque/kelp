package str

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func RsaEncrypt(publicKey, src string) ([]byte, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("illigal public key")
	}
	inter, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := inter.(*rsa.PublicKey)
	dst, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(src))
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func RsaDecrypt(privateKey string, src []byte) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("illigal private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	res, err := rsa.DecryptPKCS1v15(rand.Reader, priv, src)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
