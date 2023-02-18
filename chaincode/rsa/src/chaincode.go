package src

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	collectionPrescription  = "collectionPrescription"
	collectionPubkeyRSA     = "collectionPubkeyRSA"
	collectionReportReaders = "collectionReportReaders"
)

// String Constants for User Roles
const (
	USER_DOCTOR     = "DOCTOR"
	USER_PATIENT    = "PATIENT"
	USER_PHARMACIST = "PHARMA"
)

type SmartContract struct {
	contractapi.Contract
}

func obscureName(username string) string {
	raw := sha256.Sum256([]byte(username))
	return hex.EncodeToString(raw[:])
}

func clientObscuredName(ctx contractapi.TransactionContextInterface) (string, error) {
	// Get Client Identity
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	//Decode Identity
	identity, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}

	cnregex, err := regexp.Compile(`CN=(\w*),`)
	if err != nil {
		panic(err)
	}

	username := cnregex.FindStringSubmatch(string(identity))[1]
	return obscureName(username), nil
}
