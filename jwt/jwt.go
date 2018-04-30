package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

// JWT is a Json Web Token object.
type JWT struct {
	raw string
	alg Alg

	Header *Header
	Claims Claims
}

var (
	ErrorInvalidToken      = errors.New("invalid token")
	ErrorInvalidKey        = errors.New("invalid key")
	ErrorInvalidPublicKey  = errors.New("invalid public key")
	ErrorInvalidPrivateKey = errors.New("invalid private key")
	ErrorInvalidSign       = errors.New("invalid sign")
)

// New build an JWT entity with default value:
//     alg:    alg
//     Header: NewHeader(alg)
//     Claims: NewStdClaims()
func New(alg Alg) *JWT {
	return &JWT{
		alg:    alg,
		Header: NewHeader(alg),
		Claims: NewStdClaims(),
	}
}

// Parse read token and decode to JWT entity.
//
// The JWT entity should init header and claims for data witch to bind.
//
// Invalid token returns on error.
func (this *JWT) Parse(token string) error {
	this.raw = token
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ErrorInvalidToken
	}
	if part, err := base64Decode(parts[0]); err != nil {
		return err
	} else {
		if err := json.Unmarshal(part, &this.Header); err != nil {
			return err
		}
	}

	if part, err := base64Decode(parts[1]); err != nil {
		return err
	} else {
		if err := json.Unmarshal(part, &this.Claims); err != nil {
			return err
		}
	}
	return nil
}

// Sign for encode a token from claims by alg.
func (this *JWT) Sign() (string, error) {
	var token []byte
	if data, err := json.Marshal(this.Header); err != nil {
		return "", err
	} else {
		token = append(token, base64Encode(data)...)
	}

	if data, err := json.Marshal(this.Claims); err != nil {
		return "", err
	} else {
		token = append(token, '.')
		token = append(token, base64Encode(data)...)
	}

	if signData, err := this.alg.Sign(token); err != nil {
		return "", err
	} else {
		token = append(token, '.')
		token = append(token, base64Encode(signData)...)
	}

	this.raw = string(token)
	return this.raw, nil
}

// Verify check the token signature.
func (this *JWT) Verify(token string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ErrorInvalidToken
	}
	if sign, err := base64Decode(parts[2]); err != nil {
		return err
	} else {
		return this.alg.Verify([]byte(parts[0]+"."+parts[1]), sign)
	}
}

func base64Decode(s string) ([]byte, error) {
	return base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(s)
}

func base64Encode(s []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(s)
}
