package str_test

import (
	"fmt"
	"github.com/mapleque/kelp/str"
)

func Example_hash() {
	data := "Hello kelp hash!"

	fmt.Printf("md5: %s\n", str.Md5(data))
	fmt.Printf("sha1: %s\n", str.Sha1(data))
	fmt.Printf("sha256: %s\n", str.Sha256(data))
	fmt.Printf("sha512: %s\n", str.Sha512(data))
	// Output:
	// md5: 55bde95115f834705c52ac5c457782dd
	// sha1: 1d17377aa3dc77fe506d0fc35eb4774cd652bc41
	// sha256: 67669133b5e7b412152c86ea5aa566d62888cd49bb1e09dde73021e3c8d3cb2b
	// sha512: 938d88b8bd2322d6fb2dd4e7073aff7efbb384c6221dd0e446f4ae621a8241041fbe0ebda2bed88eade1f9efc576d18ef6469a439786208f1febaf0ba877b525
}
