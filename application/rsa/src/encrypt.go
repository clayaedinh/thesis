package src

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	keyfolder    = "rsakeys/"
	pubFilename  = "/pubkey.pem"
	privFilename = "/privkey.pem"
)

// ===============================================
// Encryption Read Parse
// ===============================================

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

// ===============================================
// Encryption Read (bytes)
// ===============================================
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

// ===============================================
// Encryption Read (rsa key type)
// ===============================================
func localPubkey(usernameHash string) (*rsa.PublicKey, error) {
	pbytes, err := readKeyfile(keyfolder + usernameHash + privFilename)
	if err != nil {
		return nil, err
	}
	pubkey, err := parsePubkeyBytes(pbytes)
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}

func localPrivkey(usernameHash string) (*rsa.PrivateKey, error) {
	pbytes, err := readKeyfile(keyfolder + usernameHash + privFilename)
	if err != nil {
		return nil, err
	}
	privkey, err := parsePrivkeyBytes(pbytes)
	if err != nil {
		return nil, err
	}
	return privkey, nil
}

/*
func keyFromChainRetrieval(arr []byte) (*rsa.PublicKey, error) {
	encoded, err := base64.StdEncoding.DecodeString(string(arr))
	if err != nil {
		return nil, fmt.Errorf("Base64 decoding of key failed: %v", nil)
	}
	return parsePubkeyBytes(encoded)
}
*/

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

// ===============================================
// Hash Username
// Hashes the username so no info is revealed
// ===============================================

func HashUsername(username string) string {
	raw := sha256.Sum256([]byte(username))
	return base64.StdEncoding.EncodeToString(raw[:])
}

// ===============================================
// Key Generation
// Generates public & private keys
// ===============================================

// from https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a
func generateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(fmt.Errorf("Failed to generate key pair: %v", err))
	}
	return privkey, &privkey.PublicKey
}

func saveKeyToFile(keyPem, filename string) error {
	pemBytes := []byte(keyPem)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %v : %v", filename, err)
	}
	_, err = file.Write(pemBytes)
	if err != nil {
		return fmt.Errorf("failed to write key to file %v : %v", filename, err)
	}
	file.Close()
	return nil
}

// from https://levelup.gitconnected.com/a-guide-to-rsa-encryption-in-go-1a18d827f35d
func exportPubKeyAsPEMStr(pubkey *rsa.PublicKey) string {
	pubKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubkey),
		},
	))
	return pubKeyPem
}

// from https://levelup.gitconnected.com/a-guide-to-rsa-encryption-in-go-1a18d827f35d
func exportPrivKeyAsPEMStr(privkey *rsa.PrivateKey) string {
	privKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privkey),
		},
	))
	return privKeyPem
}

// used to generate a new pair of keys and their PEM files
func GenerateUserKeyFiles(username string) error {
	privkey, pubkey := generateKeyPair(2048)

	pubpem := exportPubKeyAsPEMStr(pubkey)
	privpem := exportPrivKeyAsPEMStr(privkey)

	usernameHash := HashUsername(username)
	err := saveKeyToFile(pubpem, keyfolder+usernameHash+pubFilename)
	if err != nil {
		return err
	}
	err = saveKeyToFile(privpem, keyfolder+usernameHash+privFilename)
	if err != nil {
		return err
	}
	return nil
}

// ===============================================
// Decoders
// For communication with the chaincode
// ===============================================

func decodePrescriptionSet(rawgob []byte) (map[string]string, error) {
	pset := make(map[string]string)
	enc := gob.NewDecoder(bytes.NewReader(rawgob))
	err := enc.Decode(&pset)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return pset, nil
}

func decodeStringSlice(rawgob []byte) ([]string, error) {
	var strings []string
	enc := gob.NewDecoder(bytes.NewReader(rawgob))
	err := enc.Decode(&strings)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return strings, nil
}

func encodePrescriptionSet(pset *map[string]string) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(pset)
	if err != nil {
		return nil, fmt.Errorf("Failed to gob the prescription set: %v", err)
	}
	return buf.Bytes(), nil
}
