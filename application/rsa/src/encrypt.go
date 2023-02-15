package src

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

const (
	keyfolder    = "rsakeys"
	pubFilename  = "pubkey.pem"
	privFilename = "privkey.pem"
)

// ===============================================
// Encryption Read (rsa key type)
// ===============================================
func readLocalPubkey(obscureName string) (*rsa.PublicKey, error) {
	pbytes, err := readLocalKey(obscureName, pubFilename)
	if err != nil {
		return nil, err
	}
	pubkey, err := parsePubkey(pbytes)
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}

func readLocalPrivkey(obscureName string) (*rsa.PrivateKey, error) {
	pbytes, err := readLocalKey(obscureName, privFilename)
	if err != nil {
		return nil, err
	}
	privkey, err := parsePrivkey(pbytes)
	if err != nil {
		return nil, err
	}
	return privkey, nil
}

// =====================================================
// RSA Encryption and Decryption
// =====================================================

// from https://stackoverflow.com/questions/62348923/rs256-message-too-long-for-rsa-public-key-size-error-signing-jwt
func encryptBytes(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	msgLen := len(msg)
	hash := sha256.New()
	step := pub.Size() - 2*hash.Size() - 2
	var encrypted []byte
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlock, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg[start:finish], nil)
		if err != nil {
			return nil, fmt.Errorf("error encrypting bytes : %v", err)
		}
		encrypted = append(encrypted, encryptedBlock...)
	}

	return encrypted, nil
}

// from https://stackoverflow.com/questions/62348923/rs256-message-too-long-for-rsa-public-key-size-error-signing-jwt
func decryptBytes(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	msgLen := len(ciphertext)
	hash := sha256.New()
	step := priv.PublicKey.Size()
	var decrypted []byte
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		decryptedBlock, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext[start:finish], nil)
		if err != nil {
			return nil, fmt.Errorf("error decyrypting bytes : %v", err)
		}
		decrypted = append(decrypted, decryptedBlock...)
	}
	return decrypted, nil
}

// ====================================================
// Obscure Username
// Hashes and hexes the username so no info is revealed
// ====================================================

func obscureName(username string) string {
	raw := sha256.Sum256([]byte(username))
	return hex.EncodeToString(raw[:])
}

// ===============================================
// Key Generation
// Generates public & private keys
// ===============================================

// used to generate a new pair of keys and their PEM files
func GenerateUserKeyFiles(username string) {
	privkey, pubkey := generateKeyPair(2048)

	obscureName := obscureName(username)

	err := savePubkey(pubkey, obscureName)
	if err != nil {
		panic(err)
	}
	err = savePrivKey(privkey, obscureName)
	if err != nil {
		panic(err)
	}
}

// from https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a
func generateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(fmt.Errorf("Failed to generate key pair: %v", err))
	}
	return privkey, &privkey.PublicKey
}

func savePubkey(pubkey *rsa.PublicKey, obscureName string) error {
	data, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return err
	}
	saveLocalKey(data, obscureName, pubFilename)
	return nil
}

func savePrivKey(privkey *rsa.PrivateKey, obscureName string) error {
	data, err := x509.MarshalPKCS8PrivateKey(privkey)
	if err != nil {
		return err
	}
	saveLocalKey(data, obscureName, privFilename)
	return nil
}

func saveLocalKey(keyPem []byte, obscureName string, filename string) {
	os.Mkdir(keyfolder, 0777)
	os.Mkdir(filepath.Join(keyfolder, obscureName), 0777)

	file, err := os.Create(filepath.Join(keyfolder, obscureName, filename))
	if err != nil {
		fmt.Print(err)
	}
	defer file.Close()
	_, err = file.Write(keyPem)
	if err != nil {
		fmt.Print(err)
	}
}

// ===============================================
// Encoders/Decoders
// For communication with the chaincode
// ===============================================

/*
func gobEncode(data interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("GOBENCODE Failed to gob the prescription set: %v", err)
	}
	return buf.Bytes(), nil
}

func decodePrescriptionSet(rawgob []byte) (map[string]string, error) {
	pset := make(map[string]string)
	enc := gob.NewDecoder(bytes.NewReader(rawgob))
	err := enc.Decode(&pset)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return pset, nil
}
*/

func unpackageStringSlice(b64slice string) ([]string, error) {
	gobslice, err := base64.StdEncoding.DecodeString(b64slice)
	if err != nil {
		return nil, err
	}
	var strings []string
	enc := gob.NewDecoder(bytes.NewReader(gobslice))
	err = enc.Decode(&strings)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return strings, nil
}

func packagePrescriptionSet(pset map[string]string) (string, error) {

	// STEP 1: Gob-Encode
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(pset)
	if err != nil {
		return "", fmt.Errorf("Failed to gob the prescription set: %v", err)
	}

	// STEP 2: Base-64 it
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func unpackagePrescriptionSet(packaged string) (map[string]string, error) {
	rawgob, err := base64.StdEncoding.DecodeString(packaged)
	if err != nil {
		return nil, err
	}
	pset := make(map[string]string)
	enc := gob.NewDecoder(bytes.NewReader(rawgob))
	err = enc.Decode(&pset)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return pset, nil
}

// ===============================================
// Encryption Read Parse
// ===============================================

func parsePubkey(arr []byte) (*rsa.PublicKey, error) {
	parseOut, _ := x509.ParsePKIXPublicKey(arr)
	if parseOut == nil {
		return nil, fmt.Errorf("Failed to parse bytes as x509 PKIX Public Key.")
	}
	key := parseOut.(*rsa.PublicKey)
	return key, nil
}

func parsePrivkey(arr []byte) (*rsa.PrivateKey, error) {
	parseOut, _ := x509.ParsePKCS8PrivateKey(arr)
	if parseOut == nil {
		return nil, fmt.Errorf("Failed to parse bytes as x509 PKIX Public Key.")
	}
	key := parseOut.(*rsa.PrivateKey)
	return key, nil
}

// ===============================================
// Encryption Read (bytes)
// ===============================================
func readLocalKey(username, filename string) ([]byte, error) {
	data, err := os.ReadFile(filepath.Join(keyfolder, username, filename))
	if err != nil {
		return nil, fmt.Errorf("failed to read file %v : %v", filename, err)
	}
	return data, nil
}
