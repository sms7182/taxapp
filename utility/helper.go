package utility

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"sort"
	"strings"
)

func Normalize(obj interface{}) (*string, error) {
	t := reflect.TypeOf(obj)
	if kind := t.Kind(); kind != reflect.Struct {
		log.Fatalf("This program expects to work on a struct; we got a %v instead.", kind)
	}
	maps := make(map[string]interface{})
	fields := TraverseObject(t)

	sort.Strings(fields)
	value := reflect.ValueOf(&obj).Elem().Elem()

	for i := range fields {
		field := fields[i]
		splited := strings.Split(field, ".")
		if len(splited) > 1 {
			temp := value
			for j := 0; j < len(splited); j++ {
				tempValue := temp.FieldByName(splited[j])
				tempInterface := tempValue.Interface()
				if j+1 < len(splited) {
					temp = reflect.ValueOf(&tempInterface).Elem().Elem()
				} else {
					obj := reflect.ValueOf(&tempInterface).Elem()

					maps[field] = obj
				}
			}

		} else {
			fieldValue := value.FieldByName(field)

			maps[field] = fieldValue.Interface()
		}
	}
	str := ""
	for i := range fields {
		vl := maps[fields[i]]
		if vl != nil {
			if i == 0 {
				str = fmt.Sprintf("%v", vl)
			} else {
				str = fmt.Sprintf("%s#%v", str, vl)
			}
		} else {
			str = fmt.Sprintf("%s###", str)
		}
	}
	return &str, nil

}
func traversObjectWithoutReflection(obj interface{}) {
	s := reflect.ValueOf(obj)
	for _, k := range s.MapKeys() {
		fmt.Println(s.MapIndex(k))
	}

	objType := reflect.ValueOf(obj).Type()
	//var fields []string
	value := reflect.ValueOf(&obj).Elem().Elem()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)

		if field.Type.Kind() == reflect.Struct {
			tempValue := value.FieldByName(field.Name)
			traversObjectWithoutReflection(tempValue)

		} else if field.Type.Kind() == reflect.Interface {

		} else {

		}
	}
}
func TraverseObject(rType reflect.Type) []string {
	if kind := rType.Kind(); kind != reflect.Struct {
		log.Fatalf("expects to a struct type %v", kind)
	}
	var fields []string
	for i := 0; i < rType.NumField(); i++ {

		field := rType.Field(i)
		if field.Type.Kind() == reflect.Struct {

			nested_fields := TraverseObject(field.Type)
			for j := range nested_fields {
				nested := nested_fields[j]
				fields = append(fields, fmt.Sprint(field.Name, ".", nested))
			}

		} else {

			fields = append(fields, field.Name)
		}
	}
	return fields
}
func SignAndVerify(stringToBeSigned *string) (*string, error) {

	data := []byte(*stringToBeSigned)
	h := sha256.New()
	h.Write([]byte(data))
	sum := h.Sum(nil)
	signature, err := SignString(stringToBeSigned)
	if err != nil {
		log.Print("Sign has error")
	}
	publicKey, err := ioutil.ReadFile("sign.pub")
	pubPem, _ := pem.Decode(publicKey)
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {
		log.Printf("Unable to parse RSA public key")
		return nil, err
	}
	var ok bool
	var pubKey *rsa.PublicKey
	if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, err
	}
	verify(pubKey, sum, signature)
	var result string
	result = base64.StdEncoding.EncodeToString(signature)
	return &result, nil
}
func SignString(stringToBeSigned *string) ([]byte, error) {
	rng := rand.Reader
	data := []byte(*stringToBeSigned)
	h := sha256.New()
	h.Write([]byte(data))
	sum := h.Sum(nil)
	private, err := ioutil.ReadFile("sign")
	if err != nil {
		log.Print("read private key file has error")
		return nil, err
	}
	privatePem, _ := pem.Decode(private)
	if privatePem.Type != "RSA PRIVATE KEY" {
		log.Print("RSA PrivateKey is of the wrong type")
		return nil, err
	}
	var privatePemBytes []byte
	privatePemBytes = privatePem.Bytes
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privatePemBytes); err != nil {

		log.Printf("Unable to parse RSA private key, generating a temp one :%s", err.Error())
		return nil, err

	}
	var priKey *rsa.PrivateKey
	var ok bool
	priKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		log.Printf("Unable to parse RSA private key, generating a temp one : %s", err.Error())
		return nil, err
	}

	return rsa.SignPKCS1v15(rng, priKey, crypto.SHA256, sum)
}
func verify(publicKey *rsa.PublicKey, requestHashSum, signature []byte) {
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, requestHashSum, signature)
	if err != nil {
		log.Print("verify has error")
		panic(err)
	}
}
