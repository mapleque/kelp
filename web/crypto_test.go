package web

import (
	"testing"
	"time"
)

var (
	publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCUHBMlK94kToOsaIE/WmAGVYlI
ULg/+rrgNGFqwLnTTns3XuyHefYmzqyWNMTRNDImOoPaY89EDrPZYep2Yio6ls0D
kJkFiM+JQZNM9HrKTyCgByoM24e+9Cljrbd8FVU/e7cS33t9Y/C3lgv5Puk3yl9D
3WoFbrE9NpU0Ov8dbQIDAQAB
-----END PUBLIC KEY-----
`)
	privateKey = []byte(`
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
`)
	aesEcbKey = []byte(`mZqbX433M7NWCvXcOgPbFGRmDCUl8vkn`)
	signKey   = []byte(`mZqbX433M7NWCvXcOgPbFGRmDCUl8vkn`)
)

func TestBase64(t *testing.T) {
	src := []byte(`this is a message!`)
	dst := Base64Encode(src)
	t.Log(string(dst))
	if string(src) != string(Base64Decode(dst)) {
		t.Error("base64 faild")
	}
}

func TestRSA(t *testing.T) {
	src := []byte(`this is a message!`)
	cipher, err := RsaEncrypt(publicKey, src)
	if err != nil {
		t.Error(err)
	}
	dest, err := RsaDecrypt(privateKey, cipher)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(cipher))
	if string(dest) != string(src) {
		t.Errorf("dest %s not equal src %s", dest, src)
	}
}

func TestMultiRsa(t *testing.T) {
	for i := 0; i < 10000; i++ {
		src := []byte(RandMd5())
		cipher, err := RsaEncrypt(publicKey, src)
		if err != nil {
			t.Error(err)
		}
		dest, err := RsaDecrypt(privateKey, cipher)
		if err != nil {
			t.Error(err)
		}
		t.Log(string(cipher))
		if string(dest) != string(src) {
			t.Errorf("dest %s not equal src %s", dest, src)
		}
	}
}

func TestAesEcb(t *testing.T) {
	src := []byte(`{}`)
	cipher, err := AesEcbEncrypt(aesEcbKey, src)
	if err != nil {
		t.Error(err)
	}
	dest, err := AesEcbDecrypt(aesEcbKey, cipher)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(cipher))
	if string(dest) != string(src) {
		t.Errorf("dest %s not equal src %s", dest, src)
	}
}

func TestSha1Sign(t *testing.T) {
	src := []byte(`this is a message!`)
	sign := Sha1Sign(signKey, src)
	if !Sha1Verify(signKey, src, sign, 0) {
		t.Error("sign verify should pass")
	}
	t.Log(string(sign))
	time.Sleep(1 * time.Second)
	if !Sha1Verify(signKey, src, sign, 1) {
		t.Error("sign verify should pass")
	}
	if Sha1Verify(signKey, src, sign, 0) {
		t.Error("sign verify should no pass")
	}
	if Sha1Verify([]byte(`aaa`), src, sign, 0) {
		t.Error("sign verify should not pass")
	}

	timestamp := time.Now().Unix()
	signT := Sha1SignTimestamp(signKey, src, timestamp)
	if !Sha1VerifyTimestamp(signKey, src, signT, 1, timestamp) {
		t.Error("sign timestamp verify should pass")
	}
	time.Sleep(1 * time.Second)
	if Sha1VerifyTimestamp(signKey, src, signT, 0, timestamp) {
		t.Error("sign verify should not pass")
	}
}

func TestRandMd5(t *testing.T) {
	for i := 0; i < 10000; i++ {
		m := RandMd5()
		t.Log(m)
		if len(m) != 32 {
			t.Error("md5 wrong")
		}
	}
}
