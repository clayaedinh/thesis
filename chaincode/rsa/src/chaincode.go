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

type SmartContract struct {
	contractapi.Contract
}

func obscureName(username string) string {
	raw := sha256.Sum256([]byte(username))
	return hex.EncodeToString(raw[:])
}

func submittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	// Get Client Identity
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	//Decode Identity
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

func clientObscuredName(ctx contractapi.TransactionContextInterface) (string, error) {
	identity, err := submittingClientIdentity(ctx)
	if err != nil {
		return "", err
	}

	cnregex, err := regexp.Compile(`CN=(\w*),`)
	if err != nil {
		panic(err)
	}

	username := cnregex.FindStringSubmatch(identity)[1]
	return obscureName(username), nil
}
