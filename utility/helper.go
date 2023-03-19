package utility

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/gofrs/uuid"
)

type Detail struct {
	Quantity int `json:"quantity"`
	Price    int `json:"price"`
}
type DataToEncrypt struct {
	Detail []Detail `json:"detail"`

	Description string `json:"description"`
}

func JsonNormalize() {

}
func Decrypt(st string, key string) []byte {
	tagst := []byte("itismysecuretag")
	cipherTxt, err := hex.DecodeString(st)
	iv := []byte("d3fbd5bcbcd8")
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	tag := []byte(tagst)
	plaing, err := aesgcm.Open(nil, iv, cipherTxt, tag)
	if err != nil {
		panic(err.Error())
	}

	xordec := xorrox(plaing, []byte(key))

	return xordec
}

func EncryptWithIV() string {
	data := DataToEncrypt{
		Description: "test encryption",
	}
	data.Detail = append(data.Detail, Detail{
		Quantity: 2,
		Price:    250,
	})

	idwithhyphen, _ := uuid.NewV4()
	id := strings.Replace(idwithhyphen.String(), "-", "", -1)
	// iv := []byte(id)
	key := []byte(id)
	// key := make([]byte, 32)
	// if _, err := rand.Read(key); err != nil {
	// 	fmt.Printf("has error %s", err.Error())
	// // }
	// iv := make([]byte, 12)

	// if _, err := rand.Read(iv); err != nil {
	// 	fmt.Printf("has error %s", err.Error())

	// }
	iv := []byte("d3fbd5bcbcd8")

	block, err := aes.NewCipher(key)
	js, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("marshal has error")
	}
	jsxor := xorrox(js, key)
	// tag := make([]byte, 16)
	tag := []byte("itismysecuretag")
	aesgcm, err := cipher.NewGCM(block)
	result := aesgcm.Seal(nil, iv, jsxor, tag)
	return fmt.Sprintf("%x", result)

}
func xor(left []byte, right []byte) []byte {
	leftsize := len(left)
	rightsize := len(right)
	var min int
	var size int
	if leftsize > rightsize {
		size = leftsize
		min = rightsize
	} else {
		size = rightsize
		min = leftsize
	}
	val := make([]byte, size)
	for i := 0; i < min; i++ {
		val[i] = (left[i] ^ right[i])
	}
	return val
}
func xorrox(input, key []byte) (output []byte) {
	val := make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		val[i] = (input[i] ^ key[i%len(key)])
	}

	return val
}

func NormalizeJson(data interface{}) (map[string]interface{}, error) {

	// bytes, err := json.Marshal(data)
	// if err != nil {

	// 	return nil, err
	// }
	// jsonS := string(bytes)

	// fmt.Printf(jsonS)
	// x := make(map[string]interface{})
	// result:=make(map[string]interface{})
	// json.Unmarshal([]byte(jsonS), &x)
	// keys := maps.Keys(x)
	// //sort.Strings(keys)
	// for k := range keys {
	// 	if x[k] != nil && reflect.TypeOf(x[k]) == reflect.TypeOf((*Packet)(nil)) {
	// 	   nested,err:= json.Marshal(x[k])
	// 	  // if

	// 	}
	// }
	return nil, nil
}
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
		svl := fmt.Sprintf("%v", vl)

		if vl != nil && svl != "" {
			if i == 0 {
				str = fmt.Sprintf("%v", vl)
			} else {
				str = fmt.Sprintf("%s#%v", str, vl)
			}
		} else {
			str = fmt.Sprintf("%s##", str)
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
	private, err := ioutil.ReadFile("sign.key")

	if err != nil {
		log.Print("read private key file has error")
		return nil, err
	}
	privatePem, _ := pem.Decode(private)
	if privatePem.Type != "PRIVATE KEY" {
		log.Print("RSA PrivateKey is of the wrong type")
		return nil, err
	}
	var privatePemBytes []byte
	privatePemBytes = privatePem.Bytes
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS8PrivateKey(privatePemBytes); err != nil {

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
