package encrypttext

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	rand "math/rand"
	"time"

	//"golang.org/x/sync/errgroup"
	"hash"
	"io"
	"io/ioutil"
	"log"
	_ "fmt"
)

type EncryptResult struct {
	Ciphertext []byte
	EncryptedKey []byte
}

func genRandom() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 32
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf)
	return str
}

func encryptText(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(crand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func encryptKey(publicKey *rsa.PublicKey, sourceText, label []byte) (encryptedText []byte) {
	var err error
	var md5_hash hash.Hash
	md5_hash = md5.New()
	if encryptedText, err = rsa.EncryptOAEP(md5_hash, crand.Reader, publicKey, sourceText, label); err != nil {
		log.Fatal(err)
	}
	return
}

func decryptKey(publicKey *rsa.PrivateKey, sourceText, label []byte) (decryptedText []byte) {
	var err error
	var md5_hash hash.Hash
	md5_hash = md5.New()
	if decryptedText, err = rsa.DecryptOAEP(md5_hash, crand.Reader, publicKey, sourceText, label); err != nil {
		log.Fatal(err)
	}
	return
}

func EncryptItAll(pubKeyFile string, plaintext string) EncryptResult {
	var err error
	var publicKey *rsa.PublicKey
	var ciphertext, encryptedKey, label []byte

	message := []byte(plaintext)
	key := []byte(genRandom())
	ciphertext, err = encryptText(message, key)
	if err != nil {
		log.Fatalf("Error encrypting text: %s", err)
	}
	pubby, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		log.Fatalf("Error reading public key file: %s", err)
	}

	pubPem, _ := pem.Decode([]byte(pubby))

	parsedKey, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if err != nil {
		log.Fatalf("Error parsing PKIX: %s", err)
	}

	publicKey = parsedKey.(*rsa.PublicKey)
	encryptedKey = encryptKey(publicKey, key, label)
	resultStruct := EncryptResult{
		Ciphertext: ciphertext,
		EncryptedKey: encryptedKey,
	}

	return resultStruct

}

func DecryptItAll(privateKeyFile string, encryptedText string,encryptedKey string) string {

	var err error
	var privateKey *rsa.PrivateKey
	var decryptedKey, label []byte
	var 

	pubby, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		log.Fatalf("Error reading public key file: %s", err)
	}

	pubPem, _ := pem.Decode([]byte(pubby))

	parsedKey, err := x509.ParsePKCS8PrivateKey(pubPem.Bytes)
	if err != nil {
		log.Fatalf("Error parsing PKIX: %s", err)
	}

	privateKey = parsedKey.(*rsa.PrivateKey)
	decryptedKey = decryptKey(privateKey, encryptedKey, label)

	return string(decryptedKey)
}