package utility

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

const key = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCFdcxSnBJakGD5
BGEqaVmqacFdypHUhmtET4074ID0cZxPWOxFxFAwVb/0LyTYVx7MBwKnucmVaKT6
eoYsIJ0bWseSRC/ySTgxJldYz/Zu2HkNmPNEMox/dhn84haqx9Mtg9hwGOqIZifo
OcJ+E7zYKjhIjSicwkfgqecoaYyZCuaeObjtwx7K76DoAyWuxZMiDlTu+Ke5RTlp
8tpcEDqIgshGRO2vriS216SByF57bnz2nbS+J7HagYiFEHhZwJIBiEFkXagMxU16
SXZCoaphPQ0kOysc75+5xuHi/CLAgT+NObyeMebFRl1u/EW3plxFUsDkKoEOMyYS
HHqrswZHAgMBAAECggEACN0Mb1YXL/WTwyYDz/3EKzmv0mtQKRWuTRCdeCMOXW2o
LGri8jU6ACPJxk1VPJr4nCNBDWOc3DPhdRMrEwYtePIb+/5UUtqDBVyfA3J4Ut9E
lt8YFOjohNSSoEVhrQDtaQHvH27AMMPcFaO0Y4wrCA4xw7vAPTz36hdOl1P/NvC2
V5AwWh3DbBCRPfCk1fyx41zNI1wrohKDEB2ub/hd7464P6FAcsL6f8PvOu9s+Aur
+VYNL3VhPWZSBQRmHxSlXWiFhYD0lvkFKdeHe23pLo6jv1i7yLsHA5Kzjt5xKsnb
z/HtX37HRPtYJgQzLR06enG3cO6J59MqJ9CaEd3O5QKBgQC68Z8mx+FNQ5WwpHo8
lFIrxMq1EzKcilD/ODuFuVFRgO219rIQRSMENOlRAymYsmCN2jqswCsq2YIMIUQc
C/5ZMdySQaxlsPrENhAFoPVm41mMR6k4JoRRhKA6hMDTjDVEK+qFJIfbV6aWpd0O
qHP3YnRr9V6kyocUAJvvTa1wVQKBgQC2wntYG3b5vXafggjKmWAmC39pYnmrppQr
7HHsdlnAPqArfkYxUUk1G8UnmHYEqYU2mio+C036F4VyIsf4dqZnoEznr7T9E53T
0yrOEH1wXu6hmAy4rWE7wz1ZboYfSNKE9s69/mef++yzgH38qk1isygYQSufqWRB
pv/NIFSIKwKBgDNLJL4BTgJjLule1+NTVxCHWI9CizqEgSDmDv7sEDHqzE6HN+ha
7/axhesikQFCwFdrr3nC6JVDRPmLDyMa71kN41WGC4WDf+riYpcIyQzICMQCzZ2I
g/nSCBzGXBoveFYSLrEFivlWHXFsZTEma1tPel483xEcON/2ItMQXyxZAoGBAJKp
1fQp7juSuQxefRGhLhC572Cx/zQp9QSetfnuLC5j04OzzT6sndQ52ejhp+wr4lSk
OTwbNFN75sJmeRXCmd3VPYI8dkEWKfUgpFxDzXaNKHGTpLnboYklMCmB0a5vcUn1
CopcC+rOb/DJL9HBFWMcpRN50TlK5cLt8qA5zryLAoGAX9kGbjJ/7rjXESt43a7r
XXaZ+UaOIBjC2oE/KZwk7oQuSOIS9YvNDcTGuMjdm5RWkGG6nqiGHkbT1rUGl5DH
DSFlhunac6qzUevfIYGe0T7BR+JsKYLvuDixrVEAtncB7sjWd6SekpC30HrPFXOX
oX0TrJfDuNdEJ38TaCOyQ1o=
-----END PRIVATE KEY-----`

var rsaPrivateKey *rsa.PrivateKey

func initPrivateKey() error {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return errors.New("failed to Decode private key")
	}
	result, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	rsaPrivateKey = result.(*rsa.PrivateKey)
	return nil
}

func Sign(input string) (*string, error) {
	if rsaPrivateKey == nil {
		if e := initPrivateKey(); e != nil {
			return nil, e
		}
	}
	hashed := sha256.Sum256([]byte(input))
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	signedData := strings.ReplaceAll(strings.ReplaceAll(base64.URLEncoding.EncodeToString(signature), "_", "/"), "-", "+")
	return &signedData, nil
}