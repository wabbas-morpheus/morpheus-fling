package encrypttext

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	//"golang.org/x/sync/errgroup"
	"hash"
	"io"
	"io/ioutil"
	"log"
)

type EncryptResult struct {
	Ciphertext []byte
	EncryptedKey []byte
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
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func encryptKey(publicKey *rsa.PublicKey, sourceText, label []byte) (encryptedText []byte) {
	var err error
	var md5_hash hash.Hash
	md5_hash = md5.New()
	if encryptedText, err = rsa.EncryptOAEP(md5_hash, rand.Reader, publicKey, sourceText, label); err != nil {
		log.Fatal(err)
	}
	return
}

func EncryptItAll(pubKeyFile string, inputKey string, plaintext string) EncryptResult {
	var err error
	var publicKey *rsa.PublicKey
	var ciphertext, encryptedKey, label []byte

	message := []byte(plaintext)
	key := []byte(inputKey)
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