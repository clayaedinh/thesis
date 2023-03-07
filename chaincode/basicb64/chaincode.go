/*
THESIS CHAINCODE
Chan, Clay, Tan

Contains the following functions:

1. Create Prescription
2. View Prescription
3. Share Prescription
4. Edit Prescription
5. Delete Prescription
6. Set / Fill Prescription

Partly based on Hyperledger Fabric's sample chaincode.


When more nodes are created, PLEASE update the collections_config.json in the chaincode.
*/

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ============================================================ //
// DATA
// ============================================================ //

const collectionPrescription = "collectionPrescription"
const collectionAccessList = "collectionAccessList"

type SmartContract struct {
	contractapi.Contract
}

type PrescriptionAccessList struct {
	UserIds []string `json:"userIds"`
}

// String Constants for User Roles
const (
	USER_DOCTOR     = "DOCTOR"
	USER_PATIENT    = "PATIENT"
	USER_PHARMACIST = "PHARMA"
	USER_READER     = "READER"
)

// ============================================================ //
// ACCESS LIST
// ============================================================ //

func packageAccessList(access *PrescriptionAccessList) (string, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*access)
	if err != nil {
		return "", fmt.Errorf("error encoding data: %v", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func unpackageAccessList(b64access string) (*PrescriptionAccessList, error) {
	encoded, err := base64.StdEncoding.DecodeString(b64access)
	if err != nil {
		return nil, err
	}
	access := PrescriptionAccessList{}
	enc := gob.NewDecoder(bytes.NewReader(encoded))
	err = enc.Decode(&access)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return &access, nil
}

func createPackagedAccessList(obscuredName string) (string, error) {
	// Create Access List for Prescription
	accessList := PrescriptionAccessList{
		UserIds: []string{obscuredName},
	}
	return packageAccessList(&accessList)
}

// ============================================================ //
// HELPER FUNCTIONS
// ============================================================ //
func obscureName(username string) string {
	raw := sha256.Sum256([]byte(username))
	return hex.EncodeToString(raw[:])
}
func genPrescriptionId() string {
	rand.Seed(time.Now().UnixNano())
	pid := rand.Uint64()
	return fmt.Sprintf("%v", pid)
}
func packageStringSlice(strings *[]string) (string, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*strings)
	if err != nil {
		return "", fmt.Errorf("failed to gob the string slice: %v", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
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

func checkClientAccess(ctx contractapi.TransactionContextInterface, pid string) error {
	// Get Access List
	b64access, err := ctx.GetStub().GetPrivateData(collectionAccessList, pid)
	if err != nil {
		return fmt.Errorf("failed to read access list: %v", err)
	}
	if b64access == nil {
		return fmt.Errorf("%v does not exist in collection %v", pid, collectionAccessList)
	}

	// Unpackage access list
	access, err := unpackageAccessList(string(b64access))
	if err != nil {
		return err
	}

	// GET CLIENT ID
	obscuredName := clientObscuredName(ctx)

	//CONFIRM WITH ACCESS RECORD IF USER HAS ACCESS
	var matchingId bool = false
	for _, element := range access.UserIds {
		if element == obscuredName {
			matchingId = true
			break
		}
	}
	if !matchingId {
		return fmt.Errorf("permission denied. given user does not have access")
	}
	return nil
}

// ============================================================ //
// Create Prescription
// ============================================================ //
func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface, b64prescription string) (string, error) {
	// Generate ID, and check if no prescription already exists with the given id
	pid := genPrescriptionId()
	existing, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid) //get the asset from chaincode state
	if err != nil {
		return "", fmt.Errorf("failed to verify if prescription Id is already used: %v", err)
	}
	if existing != nil {
		return "", fmt.Errorf("prescription creation failed: ID already in use :%v", pid)
	}
	// get current client's obscured name
	currentUser := clientObscuredName(ctx)

	// Create access list
	b64access, err := createPackagedAccessList(currentUser)
	if err != nil {
		return "", err
	}
	// Add to the prescription collection
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64prescription))
	if err != nil {
		return "", fmt.Errorf("failed to add prescription %v to private data: %v", pid, err)
	}
	//Add the user to the access list
	err = ctx.GetStub().PutPrivateData(collectionAccessList, pid, []byte(b64access))
	if err != nil {
		return "", err
	}
	return pid, nil
}

// ============================================================ //
// Read Prescription
// ============================================================ //
func (s *SmartContract) ReadPrescription(ctx contractapi.TransactionContextInterface, pid string) (string, error) {
	err := checkClientAccess(ctx, pid)
	if err != nil {
		return "", fmt.Errorf("client access check failed: %v", err)
	}
	b64prescription, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid) //get the asset from chaincode state
	if err != nil {
		return "", fmt.Errorf("failed to read prescription: %v", err)
	}
	if b64prescription == nil {
		return "", fmt.Errorf("%v does not exist in collection", pid)
	}
	return string(b64prescription), nil
}

// ============================================================ //
// Share Prescription
// ============================================================ //
func (s *SmartContract) SharePrescription(ctx contractapi.TransactionContextInterface, pid string, newUser string) error {
	//Check if client has access to this prescription
	err := checkClientAccess(ctx, pid)
	if err != nil {
		return fmt.Errorf("client access check failed: %v", err)
	}
	// Get Access List
	b64access, err := ctx.GetStub().GetPrivateData(collectionAccessList, pid)
	if err != nil {
		return fmt.Errorf("failed to read prescription: %v", err)
	}
	if b64access == nil {
		return fmt.Errorf("%v does not exist in access list", pid)
	}
	//Unpack access list
	access, err := unpackageAccessList(string(b64access))
	if err != nil {
		return err
	}
	//Add the given user to the access list
	access.UserIds = append(access.UserIds, newUser)
	//Package data
	b64new, err := packageAccessList(access)
	if err != nil {
		return err
	}
	//Reinsert data
	err = ctx.GetStub().PutPrivateData(collectionAccessList, pid, []byte(b64new))
	if err != nil {
		return err
	}
	return nil
}

// ============================================================ //
// Update Prescription
// ============================================================ //
func (s *SmartContract) UpdatePrescription(ctx contractapi.TransactionContextInterface, pid string, b64prescription string) error {
	// Verify if a prescription already exists with the given id
	existing, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid) //get the asset from chaincode state
	if err != nil {
		return fmt.Errorf("failed to verify if prescription Id is already used: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("prescription update failed: ID not in use : %v", pid)
	}
	//Check if client has access to this prescription
	err = checkClientAccess(ctx, pid)
	if err != nil {
		return fmt.Errorf("client access check failed: %v", err)
	}

	// Add Asset to Chain
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64prescription))
	if err != nil {
		return fmt.Errorf("failed to add prescription %v to private data: %v", pid, err)
	}
	return nil
}

// ============================================================ //
// Setfill Prescription
// ============================================================ //
func (s *SmartContract) SetfillPrescription(ctx contractapi.TransactionContextInterface, pid string, b64prescription string) error {
	// Verify if a prescription already exists with the given id
	existing, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid) //get the asset from chaincode state
	if err != nil {
		return fmt.Errorf("failed to verify if prescription Id is already used: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("prescription update failed: ID not in use : %v", pid)
	}
	//Check if client has access to this prescription
	err = checkClientAccess(ctx, pid)
	if err != nil {
		return fmt.Errorf("client access check failed: %v", err)
	}
	// Add Asset to Chain
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64prescription))
	if err != nil {
		return fmt.Errorf("failed to add prescription %v to private data: %v", pid, err)
	}
	return nil
}

// ============================================================ //
// Delete Prescription
// ============================================================ //
func (s *SmartContract) DeletePrescription(ctx contractapi.TransactionContextInterface, pid string) error {
	//Check if client has access to this prescription
	err := checkClientAccess(ctx, pid)
	if err != nil {
		return fmt.Errorf("client access check failed: %v", err)
	}
	//Delete data propoer
	err = ctx.GetStub().DelPrivateData(collectionPrescription, pid)
	if err != nil {
		return fmt.Errorf("error in deleting prescription data: %v", err)
	}
	//delete access lists
	err = ctx.GetStub().DelPrivateData(collectionAccessList, pid)
	if err != nil {
		return fmt.Errorf("error in deleting access list data: %v", err)
	}
	return nil
}

// ============================================================ //
// Report Reading
// ============================================================ //
func (s *SmartContract) GetPrescriptionReport(ctx contractapi.TransactionContextInterface) (string, error) {
	// Create Iterator
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionPrescription, "", "")
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()
	// Create set of reports (all prescriptions, encrypted for given user)
	var reportSet []string
	for resultsIterator.HasNext() {
		// get each prescription
		entry, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}
		// get the prescription with the given username
		reportSet = append(reportSet, string(entry.Value))
	}
	// package the reportset
	b64reports, err := packageStringSlice(&reportSet)
	if err != nil {
		return "", err
	}
	return b64reports, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating thesis chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting thesis chaincode: %v", err)
	}
}
