package jwt_test

import (
	"encoding/json"
	"fmt"

	"github.com/mapleque/kelp/jwt"
)

// MyClaims extend jwt.StdClaims
// which implement jwt.Claims interface
type MyClaims struct {
	jwt.StdClaims
	OtherField string
}

// This example shows how to init claims
func Example_claims() {
	// make a algorithm
	key := "a secret key"
	alg := jwt.HS256(key)

	// normal jwt with default value
	j := jwt.New(alg)

	// set claims
	j.Claims.SetSubject("sub")
	j.Claims.SetAudience("aud")
	// ...
	if ret, err := json.Marshal(j); err == nil {
		fmt.Printf("claims is %s\n", string(ret))
	}

	// diy claims
	c := &MyClaims{}
	c.SetSubject("sub")
	c.OtherField = "value"
	j.Claims = c
	if ret, err := json.Marshal(j); err == nil {
		fmt.Printf("my claims is %s\n", string(ret))
	}

	// Output:
	// claims is {"Header":{"typ":"JWT","alg":"HS256"},"Claims":{"sub":"sub","aud":"aud"}}
	// my claims is {"Header":{"typ":"JWT","alg":"HS256"},"Claims":{"sub":"sub","OtherField":"value"}}
}
