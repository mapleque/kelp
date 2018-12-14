package bytes

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

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
