package src

import (
	"crypto/rsa"
	"crypto/x509"
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

func ParsePubkeyBytes(arr []byte) (*rsa.PublicKey, error) {
	parseOut, _ := x509.ParsePKIXPublicKey(arr)
	if parseOut == nil {
		return nil, fmt.Errorf("Failed to parse bytes as x509 PKIX Public Key.")
	}
	key := parseOut.(*rsa.PublicKey)
	return key, nil
}

func ParsePrivkeyBytes(arr []byte) (*rsa.PrivateKey, error) {
	parseOut, _ := x509.ParsePKCS8PrivateKey(arr)
	if parseOut == nil {
		return nil, fmt.Errorf("Failed to parse bytes as x509 PKIX Public Key.")
	}
	key := parseOut.(*rsa.PrivateKey)
	return key, nil
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
