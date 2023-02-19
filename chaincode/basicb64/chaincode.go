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

// SmartContract provides functions for managing Assets such as Prescription

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
)

// ============================================================ //
// PRESCRIPTIONS
// ============================================================ //

type Prescription struct {
	Brand          string `json:"Brand"`
	Dosage         string `json:"Dosage"`
	PatientName    string `json:"PatientName"`
	PatientAddress string `json:"PatientAddress"`
	PrescriberName string `json:"PrescriberName"`
	PrescriberNo   uint32 `json:"PrescriberNo"`
	PiecesTotal    uint8  `json:"AmountTotal"`
	PiecesFilled   uint8  `json:"AmountFilled"`
}

func genPrescriptionId() string {
	rand.Seed(time.Now().UnixNano())
	pid := rand.Uint64()
	return fmt.Sprintf("%v", pid)
}

func createPackagedPrescription() (string, error) {
	var prescription Prescription
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(prescription)
	if err != nil {
		return "", fmt.Errorf("error encoding prescription data: %v", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

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

// ============================================================ //
// CLIENT ACCESS CHECK
// Verifies that the current client identity has access
// ============================================================ //
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
	obscuredName, err := clientObscuredName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

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
// GET MY ID
// ============================================================ //
func (s *SmartContract) GetMyID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientId, err := submittingClientIdentity(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get client id: %v", err)
	}
	return clientId, nil
}

// ============================================================ //
// Create Prescription
// ============================================================ //
func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface) (string, error) {

	//Generate new prescription id
	pid := genPrescriptionId()

	// Verify if no prescription already exists with the given id
	existing, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid) //get the asset from chaincode state
	if err != nil {
		return "", fmt.Errorf("failed to verify if prescription Id is already used: %v", err)
	}
	if existing != nil {
		return "", fmt.Errorf("prescription creation failed: ID already in use :%v", pid)
	}

	// Verify if current user is a Patient
	err = ctx.GetClientIdentity().AssertAttributeValue("role", USER_PATIENT)
	if err != nil {
		return "", err
	}

	//Generate new prescription
	b64prescription, err := createPackagedPrescription()
	if err != nil {
		return "", fmt.Errorf("failed to generate new prescription: %v", err)
	}
	// get current client's obscured name
	clientName, err := clientObscuredName(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve client name: %v", err)
	}

	// Create access list
	b64access, err := createPackagedAccessList(clientName)
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

	// Verify if current user is a Doctor
	err = ctx.GetClientIdentity().AssertAttributeValue("role", USER_DOCTOR)
	if err != nil {
		return err
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

	// Verify if current user is a Pharmacist
	err = ctx.GetClientIdentity().AssertAttributeValue("role", USER_PHARMACIST)
	if err != nil {
		return err
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
// Delete PRESCRIPTION
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

func (s *SmartContract) ReadPrescription(ctx contractapi.TransactionContextInterface, pid string) (string, error) {
	//Check if client has access to this prescription
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
// SHARE PRESCRIPTION
// ============================================================ //

func (s *SmartContract) SharePrescription(ctx contractapi.TransactionContextInterface, pid string, newUser string) error {

	// Verify if current user is a Patient
	err := ctx.GetClientIdentity().AssertAttributeValue("role", USER_PATIENT)
	if err != nil {
		return fmt.Errorf("current user must be patient: %v", err)
	}

	//Check if client has access to this prescription
	err = checkClientAccess(ctx, pid)
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

	// Transfer asset in private data collection to new owner
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

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating thesis chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting thesis chaincode: %v", err)
	}
}
