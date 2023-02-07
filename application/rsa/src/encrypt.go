package src

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

func readKeyfile(filename string) ([]byte, error) {
	rawfile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %v : %v", filename, err)
	}
	block, _ := pem.Decode(rawfile)
	if block == nil {
		return nil, fmt.Errorf("failed to parse file %v as .pem", filename)
	}
	return block.Bytes, nil
}

func ReadUserPubkey(username string) ([]byte, error) {
	pathname := "rsakeys/" + username + "/pubkey.pem"
	return readKeyfile(pathname)
}
func ReadUserPrivkey(username string) ([]byte, error) {
	pathname := "rsakeys/" + username + "/privkey.pem"
	return readKeyfile(pathname)
}

func parsePubkeyBytes(arr []byte) (*rsa.PublicKey, error) {
	parseOut, _ := x509.ParsePKIXPublicKey(arr)
	if parseOut == nil {
		return nil, fmt.Errorf("Failed to parse bytes as x509 PKIX Public Key.")
	}
	key := parseOut.(*rsa.PublicKey)
	return key, nil
}

func parsePrivkeyBytes(arr []byte) (*rsa.PrivateKey, error) {
	parseOut, _ := x509.ParsePKCS8PrivateKey(arr)
	if parseOut == nil {
		return nil, fmt.Errorf("Failed to parse bytes as x509 PKIX Public Key.")
	}
	key := parseOut.(*rsa.PrivateKey)
	return key, nil
}

func keyFromChainRetrieval(arr []byte) (*rsa.PublicKey, error) {
	encoded, err := base64.StdEncoding.DecodeString(string(arr))
	if err != nil {
		return nil, fmt.Errorf("Base64 decoding of key failed: %v", nil)
	}
	return parsePubkeyBytes(encoded)
}

// From https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a
func encryptBytes(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("error encrypting bytes : %v", err)
	}
	return ciphertext, nil
}

// From https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a
func decryptBytes(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.New()
	msg, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error decyrypting bytes : %v", err)
	}
	return msg, nil
}

/*
func RetrieveLocalUserPubkey(username string) (*rsa.PublicKey, error) {
	keybytes, err := readKeyfile(username)
	if err != nil {
		return nil, err
	}
	return pubkeyFromBytes(keybytes)
}

func RetrieveLocalUserPrivkey(username string) (*rsa.PrivateKey, error) {
	keybytes, err := readKeyfile(username)
	if err != nil {
		return nil, err
	}
	return privkeyFromBytes(keybytes)
}
*/
