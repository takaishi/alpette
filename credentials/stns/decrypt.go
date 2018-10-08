package stns

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

func (tc *stnsTC) Decrypt(s string, key *rsa.PrivateKey) (string, error) {
	log.Printf("[DEBUG] s: %s\n", s)
	log.Printf("[DEBUG] key: %#v\n", key)
	in, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, key, in)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func (tc *stnsTC) readPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("invalid private key data")
	}

	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		return nil, errors.New(fmt.Sprintf("invalid key type: %s", block.Type))
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}
