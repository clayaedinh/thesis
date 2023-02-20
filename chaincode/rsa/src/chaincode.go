package src

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"time"

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
	USER_READER     = "READER"
)

type SmartContract struct {
	contractapi.Contract
}

func obscureName(username string) string {
	raw := sha256.Sum256([]byte(username))
	return hex.EncodeToString(raw[:])
}

// Verify if the given prescription set has a prescription encrypted with the current client
func verifyClientAccess(ctx contractapi.TransactionContextInterface, b64pset string) error {

	//get current user
	currentUser, err := clientObscuredName(ctx)
	if err != nil {
		return err
	}

	// Unpackage Prescription Set
	pset, err := unpackagePrescriptionSet(string(b64pset))
	if err != nil {
		return err
	}

	_, exists := (*pset)[currentUser]
	if !exists {
		return fmt.Errorf("client does not have access to the given prescription set")
	}
	return nil
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

func genPrescriptionId() string {
	rand.Seed(time.Now().UnixNano())
	pid := rand.Uint64()
	return fmt.Sprintf("%v", pid)
}
