package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

// rsasha implement Alg interface
type rsasha struct {
	name       string
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	hash       crypto.Hash
}

// RS256 is an crypto algorithm using RSA and SHA-256
func RS256(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) Alg {
	return &rsasha{
		"RS256",
		publicKey,
		privateKey,
		crypto.SHA256,
	}
}

// RS384 is an crypto algorithm using RSA and SHA-384
func RS384(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) Alg {
	return &rsasha{
		"RS384",
		publicKey,
		privateKey,
		crypto.SHA384,
	}
}

// RS512 is an crypto algorithm using RSA and SHA-512
func RS512(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) Alg {
	return &rsasha{
		"RS512",
		publicKey,
		privateKey,
		crypto.SHA512,
	}
}

func (this *rsasha) Name() string {
	return this.name
}

func (this *rsasha) Sign(data []byte) ([]byte, error) {
	if this.privateKey == nil {
		return nil, ErrorInvalidPrivateKey
	}

	h := this.hash.New()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}

	sign, err := rsa.SignPKCS1v15(rand.Reader, this.privateKey, this.hash, h.Sum(nil))
	if err != nil {
		return nil, err
	}
	return sign, nil

}

func (this *rsasha) Verify(data, sign []byte) error {
	if this.publicKey == nil {
		return ErrorInvalidPublicKey
	}
	h := this.hash.New()
	if _, err := h.Write(data); err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(this.publicKey, this.hash, h.Sum(nil), sign)
}
