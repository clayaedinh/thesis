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

func genPrescriptionId() string {
	rand.Seed(time.Now().UnixNano())
	pid := rand.Uint64()
	return fmt.Sprintf("%v", pid)
}

func checkIfUserPubkeyExists(ctx contractapi.TransactionContextInterface, obscuredName string) error {
	//verify if user has public key information (i.e. if the user exists properly)
	val, err := ctx.GetStub().GetPrivateData(collectionPubkeyRSA, obscuredName)
	if err != nil {
		return fmt.Errorf("failed to verify if user has RSA public key information :%v", err)
	}
	if val == nil {
		return fmt.Errorf("given user does not have RSA public key information")
	}
	return nil
}

// ============================================================ //
// Unpackage & Check Access
// unpackages a set of prescriptions, checks if current user
// has access to any of the prescriptions inside of it
// ============================================================ //
func unpackageAndCheckAccess(ctx contractapi.TransactionContextInterface, b64pset string, obscureName string) (*map[string]string, error) {
	pset, err := unpackagePrescriptionSet(string(b64pset))
	if err != nil {
		return nil, err
	}
	_, exists := (*pset)[obscureName]
	if !exists {
		return nil, fmt.Errorf("client does not have access to the given prescription set")
	}
	return pset, nil
}

// ============================================================ //
// CLIENT IDENTITY
// From certificate, get identity, obscured name, and check
// if user if allowed to access chaincode
// ============================================================ //
func clientObscuredName(ctx contractapi.TransactionContextInterface) string {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		panic(fmt.Errorf("failed to read clientID: %v", err))
	}
	identity, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		panic(fmt.Errorf("failed to base64 decode clientID: %v", err))
	}
	cnregex, err := regexp.Compile(`CN=(\w*),`)
	if err != nil {
		panic(err)
	}
	username := cnregex.FindStringSubmatch(string(identity))[1]
	return obscureName(username)
}
