package jwt

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"math/big"
)

// ecdsasha implement Alg interface
type ecdsasha struct {
	name       string
	publicKey  *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
	hash       func() hash.Hash
}

// ES256 is an crypto algorithm using ECDSA and SHA-256
func ES256(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) Alg {
	return &ecdsasha{
		"ES256",
		publicKey,
		privateKey,
		sha256.New,
	}
}

// ES384 is an crypto algorithm using ECDSA and SHA-384
func ES384(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) Alg {
	return &ecdsasha{
		"ES384",
		publicKey,
		privateKey,
		sha512.New384,
	}
}

// ES512 is an crypto algorithm using ECDSA and SHA-512
func ES512(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) Alg {
	return &ecdsasha{
		"ES512",
		publicKey,
		privateKey,
		sha512.New,
	}
}

func (this *ecdsasha) Name() string {
	return this.name
}

func (this *ecdsasha) Sign(data []byte) ([]byte, error) {
	if this.privateKey == nil {
		return nil, ErrorInvalidPrivateKey
	}

	h := this.hash()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}

	r, s, err := ecdsa.Sign(rand.Reader, this.privateKey, h.Sum(nil))
	if err != nil {
		return nil, err
	}

	byteSize := this.byteSize(this.privateKey.Params().BitSize)
	rbytes := r.Bytes()
	rsign := make([]byte, byteSize)
	copy(rsign[byteSize-len(rbytes):], rbytes)
	sbytes := s.Bytes()
	ssign := make([]byte, byteSize)
	copy(ssign[byteSize-len(sbytes):], sbytes)
	return append(rsign, ssign...), nil
}

func (this *ecdsasha) Verify(data, sign []byte) error {
	if this.publicKey == nil {
		return ErrorInvalidPublicKey
	}
	byteSize := this.byteSize(this.privateKey.Params().BitSize)
	if len(sign) != byteSize*2 {
		return ErrorInvalidSign
	}
	r := big.NewInt(0).SetBytes(sign[:byteSize])
	s := big.NewInt(0).SetBytes(sign[byteSize:])
	h := this.hash()
	if _, err := h.Write(data); err != nil {
		return err
	}
	if !ecdsa.Verify(this.publicKey, h.Sum(nil), r, s) {
		return ErrorInvalidSign
	}
	return nil
}

func (this *ecdsasha) byteSize(bitSize int) int {
	byteSize := bitSize / 8
	if bitSize%8 > 0 {
		return byteSize + 1
	}
	return byteSize
}
