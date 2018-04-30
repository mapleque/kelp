package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

// hmacsha implemnt Alg interface
type hmacsha struct {
	name string
	hash func() hash.Hash
	key  []byte
}

// HS256 is an crypto alogorithm using HMAC and SHA256
func HS256(key string) Alg {
	return &hmacsha{
		"HS256",
		sha256.New,
		[]byte(key),
	}
}

// HS384 is an crypto alogorithm using HMAC and SHA384
func HS384(key string) Alg {
	return &hmacsha{
		"HS384",
		sha512.New384,
		[]byte(key),
	}
}

// HS512 is an crypto alogorithm using HMAC and SHA512
func HS512(key string) Alg {
	return &hmacsha{
		"HS512",
		sha512.New,
		[]byte(key),
	}
}

func (this *hmacsha) Name() string {
	return this.name
}

func (this *hmacsha) Sign(data []byte) ([]byte, error) {
	if len(this.key) == 0 {
		return nil, ErrorInvalidKey
	}

	h := hmac.New(this.hash, this.key)
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (this *hmacsha) Verify(data, sign []byte) error {
	tar, err := this.Sign(data)
	if err != nil {
		return err
	}
	if !hmac.Equal(tar, sign) {
		return ErrorInvalidSign
	}
	return nil
}
