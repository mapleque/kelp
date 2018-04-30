package jwt_test

import (
	"encoding/json"
	"fmt"

	"github.com/mapleque/kelp/jwt"
)

// This example shows how to use JWT
func Example_jwt() {
	// make a algorithm
	key := "a secret key"
	alg := jwt.HS256(key)

	// normal jwt with default value
	j := jwt.New(alg)

	// set claims
	j.Claims.SetSubject("sub")

	// sign
	token, err := j.Sign()
	if err == nil {
		fmt.Printf("sign token is %s\n", token)
	} else {
		fmt.Printf("sign error %v", err)
	}

	// verify
	if err := j.Verify(token); err == nil {
		fmt.Printf("verify pass\n")
	} else {
		fmt.Printf("verify error %v", err)
	}

	if err := j.Parse(token); err == nil {
		ret, _ := json.Marshal(j)
		fmt.Printf("claims is %v\n", string(ret))
	} else {
		fmt.Printf("parse error %v", err)
	}
	// Output:
	// sign token is eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJzdWIifQ.7FKVJPwdyL3lZ_BP3CBC1P-Ghoq7MRNAphUkYUXyUMU
	// verify pass
	// claims is {"Header":{"typ":"JWT","alg":"HS256"},"Claims":{"sub":"sub"}}
}
