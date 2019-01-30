package str_test

import (
	"fmt"
	"github.com/mapleque/kelp/str"
)

func Example_rsa() {
	publicKey := `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCUHBMlK94kToOsaIE/WmAGVYlI
ULg/+rrgNGFqwLnTTns3XuyHefYmzqyWNMTRNDImOoPaY89EDrPZYep2Yio6ls0D
kJkFiM+JQZNM9HrKTyCgByoM24e+9Cljrbd8FVU/e7cS33t9Y/C3lgv5Puk3yl9D
3WoFbrE9NpU0Ov8dbQIDAQAB
-----END PUBLIC KEY-----
`
	privateKey := `
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCUHBMlK94kToOsaIE/WmAGVYlIULg/+rrgNGFqwLnTTns3XuyH
efYmzqyWNMTRNDImOoPaY89EDrPZYep2Yio6ls0DkJkFiM+JQZNM9HrKTyCgByoM
24e+9Cljrbd8FVU/e7cS33t9Y/C3lgv5Puk3yl9D3WoFbrE9NpU0Ov8dbQIDAQAB
AoGABSzd9mSMBJTBwRp9uar8w/vlKiO37HRkZ0UtSj+lvp51a7/jX/CBC2YZXb5G
SlEal39f8BegvG4PFr93I9/WPdw82UkX/4q2eJaVlLJ6blezyM+NwFDkD2M4yycP
gmkdmEAFgfT5D8Zn4oIrznUgjL+24EOa3zp19rQ5Mc0wjAECQQDKXDS4fTwgmbEZ
5P952MEN5d3zHcQQw9MJgVuZoh1nJXqwo51IrvUPFAy/V5lVaxHIn5oxpOY9mD1I
XRVe5ybtAkEAu16HFMds/R0a4z8KLLNiOZ0TiSQ42JXVMRBm5Jw/iGad7uBDAY6E
XdfetAOs8W3XRuCIFICc6zUKav0XqYaAgQJABmaeQEut0DYsVO5aamdBzAe+WodR
gVpAXaea1yQ6m92ioN28BuWJ2N1AffjuX7ZQTLFHtlRJ+B7NqXFQUL0tDQJBALhL
ANiKUwQfNYwhPFO9WTbL7iQtMZCux1QMCvh/SupR7LPBd4a3dDCNnKo5F0kcvesj
/BUWb8HVmNqk+DoxZoECQBgLw05kviJR7bvPCaokQ/0vw8nnOmlXrQiCZs+RX2LC
dOB6eZ2iyilVMWFn03P4VxUBNFhDOT0cQeyOGURHRyk=
-----END RSA PRIVATE KEY-----
`
	data := "Hello kelp rsa!"

	// rsa encrypt
	encryptData, _ := str.RsaEncrypt(publicKey, data)
	encodingEncryptData := str.Base64Encode(encryptData)
	fmt.Printf("encrypt data is alway different\n")

	// rsa decrypt
	decryptData := str.Base64Decode(encodingEncryptData)
	targetData, _ := str.RsaDecrypt(privateKey, decryptData)
	fmt.Printf("target data: %s\n", targetData)

	// Output:
	// encrypt data is alway different
	// target data: Hello kelp rsa!
}
