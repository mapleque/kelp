package jwt

// Claims is an interface of claims
type Claims interface {
	Valid() error
	SetIssuer(iss string)
	SetSubject(sub string)
	SetAudience(aud string)
	SetExpireationTime(exp int64)
	SetNotBefore(nbf int64)
	SetJWTID(jti string)
	GetIssuer() string
	GetSubject() string
	GetAudience() string
	GetExpireationTime() int64
	GetNotBefore() int64
	GetIssuedAt() int64
	GetJWTID() string
}

// StdClaims implement Claims interface
// which include all standard properties.
//
// Extend this to add more public value if you need.
type StdClaims struct {
	Issuer          string `json:"iss,omitempty"`
	Subject         string `json:"sub,omitempty"`
	Audience        string `json:"aud,omitempty"`
	ExpireationTime int64  `json:"exp,omitempty"`
	NotBefore       int64  `json:"nbf,omitempty"`
	IssuedAt        int64  `json:"iat,omitempty"`
	JWTID           string `json:"jti,omitempty"`
}

// NewStdClaims returns a StdClaims entity with default value
func NewStdClaims() Claims {
	return &StdClaims{}
}

// Valid
//
// TODO implement valid method
func (this *StdClaims) Valid() error {
	return nil
}

func (this *StdClaims) SetIssuer(iss string) {
	this.Issuer = iss
}
func (this *StdClaims) SetSubject(sub string) {
	this.Subject = sub
}
func (this *StdClaims) SetAudience(aud string) {
	this.Audience = aud
}
func (this *StdClaims) SetExpireationTime(exp int64) {
	this.ExpireationTime = exp
}
func (this *StdClaims) SetNotBefore(nbf int64) {
	this.NotBefore = nbf
}
func (this *StdClaims) SetJWTID(jti string) {
	this.JWTID = jti
}
func (this *StdClaims) GetIssuer() string {
	return this.Issuer
}
func (this *StdClaims) GetSubject() string {
	return this.Subject
}
func (this *StdClaims) GetAudience() string {
	return this.Audience
}
func (this *StdClaims) GetExpireationTime() int64 {
	return this.ExpireationTime
}
func (this *StdClaims) GetNotBefore() int64 {
	return this.NotBefore
}
func (this *StdClaims) GetIssuedAt() int64 {
	return this.IssuedAt
}
func (this *StdClaims) GetJWTID() string {
	return this.JWTID
}
