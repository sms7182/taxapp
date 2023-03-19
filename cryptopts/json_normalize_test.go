package cryptops

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestSortKeys(t *testing.T) {
	j := []byte(`
	{
		"k2" : "v1",
		"k4" : "v2",
		"k3" : {
			"k1" : "v4",
			"k5" : "v5"
		}
	}`)

	obj := make(map[string]interface{})
	json.Unmarshal(j, &obj)

	vals := sortJsonMap(obj, "")
	fmt.Println(vals)
	fmt.Println(normalizeVals(vals))
}

func TestBl(t *testing.T) {
	fmt.Println(float64String(2.3))
	fmt.Println(float64String(6.0))
	fmt.Println(float64String(7.45695))
	fmt.Println(float64String(0.897))
}

var (
	defaultPubKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAscetx8F1Q7H67ZSgIcTw
zQfCf919iACii2o5sh+1l7N62nE9zBpSx3OEgNv64l8v4OchXMU8gKk28piExpdQ
kvzDW5VK7STmEuIZ7IqWKsZge1YmGDsyIFw74V9Uslhc05t7VKYhWWFPAKfouPPM
3ZKe5ZiALAjvLVIEUYYnQ452H2RJGuGYJeKvPiNtOwKSLA/ROwvE/1I+0S+gq1hd
+GbrYPJLfj77pkZJKnf/ye3rgbQglfBQzSHSDKuwC6xNZEWMR4DBzraE0MeKrNhN
4PFxKpkyVRPftrahbiTA6ohvoBsSyD+RdT1dRde4qbJGXuW6AQ2DQYNPWqTdVDNo
kwIDAQAB
-----END PUBLIC KEY-----
`
	defaultPriKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAscetx8F1Q7H67ZSgIcTwzQfCf919iACii2o5sh+1l7N62nE9
zBpSx3OEgNv64l8v4OchXMU8gKk28piExpdQkvzDW5VK7STmEuIZ7IqWKsZge1Ym
GDsyIFw74V9Uslhc05t7VKYhWWFPAKfouPPM3ZKe5ZiALAjvLVIEUYYnQ452H2RJ
GuGYJeKvPiNtOwKSLA/ROwvE/1I+0S+gq1hd+GbrYPJLfj77pkZJKnf/ye3rgbQg
lfBQzSHSDKuwC6xNZEWMR4DBzraE0MeKrNhN4PFxKpkyVRPftrahbiTA6ohvoBsS
yD+RdT1dRde4qbJGXuW6AQ2DQYNPWqTdVDNokwIDAQABAoIBACnS6iVGdAn7AyeF
ga6wIF575twiBXhLffICiZRINXZ8+PgPEBTGVJcrrA6MshczgZYNiiHDHRq/tHea
PhJiYshRwrv3AWuM9LuYibTGXdGuXeBmQgwNURuf106MGObkNuJpf7hIZSwb4nQr
DGsGoDm4Vr15BR5W873bv7xWLUKNCpPwo65pGJHTCTjBm6AC0doQ/WbN+V9ly2B6
uz2Afl9wrQlTZUUHLFuvO9IjukCCj0ZclPocsURA0j3TF47kXmZxhYT5wCerzZRO
tgfJXd0sGvoBZTS5OpEVG3ef/EFbyDLy76QwsItaphJlgCCyXwFOTDxDQQWd6uPA
3V9/lWECgYEA6p74OccfFythzHZN/SOk4bC5bUKWw6umE0h0t3lQN7bCtddDu2RH
xR4SFlIA8fL6Vp4BkseydRg3mMHz+UZ5E5EaRPWWQHU0Spn9qnrO0QTuWeQJlpN0
6f0Am2pZeh6voRbbzAC3yKGI5frdDLtM2p2k/gbGlWTB5e0L2UGAUNkCgYEAwfrF
7UWIrCqx2abgsGNfHp4omwhfv8jpD4CGpXKGrHvnagGfLYABngbmIo0GLHyUR0Rm
wE2qfeDp+64vvNj+RV4lRME1PNFsWxaJ8eMUHr06lDO51Cy1lhTWymT4NXj+Esys
dFJvCElfwxbZjflyNf8hfkSa24Rfo6WoI9jV4UsCgYEA4HJZlrRVms2mjnmym/LI
Xhu5F7v3DJMdmh7bgVWtls7gsCKRqigBvKHKvc2PF+bQ86HOcYNWxkv3i8wnwJVZ
aI2MauHh7iHxd1ifYcKALVchSZ8sSP8hfmLJfOQdWwUWEO4UMLGTH3zgwNnfM7nO
iOj8mQMUYIB2OaYuipTt0ukCgYEAl4qRHAdJea81GCNtv38ybVoDwPIu00ZjBNBU
4GXzXkbCCCfSMhqhqNIc8fsYSqLcuDxwxWUnf4W5ZfyzoKYpJwogtXD3ZVb6fsLB
662KJ2WPoP4z+9Ud22zWTHHLEwM+AnPRemJ4CZJA9MkiFu88UYDKqrlv/XSRvugI
zlB07rcCgYEAueo9hE02p0iSqxXWru8zu7PxY8Gy2+tksMZb4PWB5C732BMr3ryP
lz5UUW+5iBe/z54HOdmBbVdd3G+fRlkCm9XUex0GlwaN3g45k8rcyJi/8iRexIpF
2c3olpk+wO+d7ciK+7Qc8uHYyZlnBxQu6FIRDTE/Y8QOkU97/BDSkYQ=
-----END RSA PRIVATE KEY-----
`
)

func TestEncrypt(t *testing.T) {
	text := "text to encrypt"

	block, _ := pem.Decode([]byte(defaultPubKey))
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		os.Stderr.Write([]byte("Failed to parse public key\n"))
		panic(err)
	}
	pub := pubInterface.(*rsa.PublicKey)

	c, err := EncryptRsaOAEPSHA256([]byte(text), pub)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(c))

	block, _ = pem.Decode([]byte(defaultPriKey))
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pri, c, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(plaintext))
}
