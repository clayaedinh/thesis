package application

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

func pubkeyFromBytes(arr []byte) (*rsa.PublicKey, error) {
	parseOut, _ := x509.ParsePKIXPublicKey(arr)
	if parseOut == nil {
		return nil, fmt.Errorf("Failed to parse bytes as x509 PKIX Public Key.")
	}
	key := parseOut.(*rsa.PublicKey)
	return key, nil
}

func privkeyFromBytes(arr []byte) (*rsa.PrivateKey, error) {
	parseOut, _ := x509.ParsePKCS8PrivateKey(arr)
	if parseOut == nil {
		return nil, fmt.Errorf("Failed to parse bytes as x509 PKIX Public Key.")
	}
	key := parseOut.(*rsa.PrivateKey)
	return key, nil
}
