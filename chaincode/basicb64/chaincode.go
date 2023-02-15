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

package basic

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
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

// Structure for Prescription Data
type Prescription struct {
	Id             uint64 `json:"Id"`
	Brand          string `json:"Brand"`
	Dosage         string `json:"Dosage"`
	PatientName    string `json:"PatientName"`
	PatientAddress string `json:"PatientAddress"`
	PrescriberName string `json:"PrescriberName"`
	PrescriberNo   uint32 `json:"PrescriberNo"`
	PiecesTotal    uint8  `json:"AmountTotal"`
	PiecesFilled   uint8  `json:"AmountFilled"` // in terms of percentage
}

type PrescriptionAccessList struct {
	PID     uint64   `json:"prescriptionId"`
	UserIds []string `json:"userIds"`
}

// String Constants for User Roles
const (
	USER_DOCTOR     string = "DOCTOR"
	USER_PATIENT           = "PATIENT"
	USER_PHARMACIST        = "PHARMA"
)

// ============================================================ //
// ENCODING
// ============================================================ //

func packageAccessList(access *PrescriptionAccessList) (string, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*access)
	if err != nil {
		return "", fmt.Errorf("error encoding data %v: %v", access.PID, err)
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

func packagePrescription(prescription *Prescription) (string, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*prescription)
	if err != nil {
		return "", fmt.Errorf("error encoding data %v: %v", prescription.Id, err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func unpackagePrescription(b64prescription string) (*Prescription, error) {
	encoded, err := base64.StdEncoding.DecodeString(b64prescription)
	if err != nil {
		return nil, err
	}
	prescription := Prescription{}
	enc := gob.NewDecoder(bytes.NewReader(encoded))
	err = enc.Decode(&prescription)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return &prescription, nil
}

// ============================================================ //
// HELPER FUNCTIONS
// ============================================================ //

func clientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	// Get Client Identity
	b64ID, err := ctx.GetClientIdentity().GetID()

	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	//Decode Identity
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

func clientIdentityB64(ctx contractapi.TransactionContextInterface) (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	return b64ID, nil
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
	b64clientId, err := clientIdentityB64(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get client id: %v", err)
	}

	//CONFIRM WITH ACCESS RECORD IF USER HAS ACCESS
	var matchingId bool = false
	for _, element := range access.UserIds {
		if element == b64clientId {
			matchingId = true
			break
		}
	}
	if matchingId == false {
		return fmt.Errorf("Permission Denied. Given User does not have access.")
	}
	return nil
}

// ============================================================ //
// GET MY ID
// ============================================================ //
func (s *SmartContract) GetMyID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientId, err := clientIdentityB64(ctx)
	if err != nil {
		return "", fmt.Errorf("Failed to get client id: %v", err)
	}
	return clientId, nil
}

// ============================================================ //
// Generate Prescription Id
// ============================================================ //
func genPrescriptionId() uint64 {
	rand.Seed(time.Now().UnixNano())
	pid := rand.Uint64()
	log.Printf("Generated prescription id: %v", pid)
	return pid
}

// ============================================================ //
// Create Prescription
// ============================================================ //
func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface, pid string, b64prescription string) error {
	// Verify if no prescription already exists with the given id
	existing, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid) //get the asset from chaincode state
	if err != nil {
		return fmt.Errorf("Failed to verify if prescription Id is already used: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("Prescription creation failed: ID already in use :")
	}

	// Verify if current user is a Patient
	err = ctx.GetClientIdentity().AssertAttributeValue("role", USER_PATIENT)
	if err != nil {
		return err
	}

	//Check if client has access to this prescription
	err = checkClientAccess(ctx, pid)
	if err != nil {
		return fmt.Errorf("Client access check failed: %v", err)
	}

	// Add Asset to Chain
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64prescription))
	if err != nil {
		return fmt.Errorf("Failed to add prescription %v to private data: %v", pid, err)
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
		return fmt.Errorf("Failed to verify if prescription Id is already used: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("Prescription update failed: ID not in use :")
	}

	// Verify if current user is a Doctor
	err = ctx.GetClientIdentity().AssertAttributeValue("role", USER_DOCTOR)
	if err != nil {
		return err
	}

	//Check if client has access to this prescription
	err = checkClientAccess(ctx, pid)
	if err != nil {
		return fmt.Errorf("Client access check failed: %v", err)
	}

	// Add Asset to Chain
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64prescription))
	if err != nil {
		return fmt.Errorf("Failed to add prescription %v to private data: %v", pid, err)
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
		return fmt.Errorf("Failed to verify if prescription Id is already used: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("Prescription update failed: ID not in use :")
	}

	// Verify if current user is a Pharmacist
	err = ctx.GetClientIdentity().AssertAttributeValue("role", USER_PHARMACIST)
	if err != nil {
		return err
	}

	//Check if client has access to this prescription
	err = checkClientAccess(ctx, pid)
	if err != nil {
		return fmt.Errorf("Client access check failed: %v", err)
	}

	// Add Asset to Chain
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64prescription))
	if err != nil {
		return fmt.Errorf("Failed to add prescription %v to private data: %v", pid, err)
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
		return fmt.Errorf("Client access check failed: %v", err)
	}

	//Delete data propoer
	err = ctx.GetStub().DelPrivateData(collectionPrescription, pid)
	if err != nil {
		return fmt.Errorf("Error in deleting prescription data: %v", err)
	}

	//delete access lists
	err = ctx.GetStub().DelPrivateData(collectionAccessList, pid)
	if err != nil {
		return fmt.Errorf("Error in deleting access list data: %v", err)
	}
	return nil
}

func (s *SmartContract) ReadPrescription(ctx contractapi.TransactionContextInterface, pid string) (string, error) {
	//Check if client has access to this prescription
	err := checkClientAccess(ctx, pid)
	if err != nil {
		return "", fmt.Errorf("Client access check failed: %v", err)
	}
	b64prescription, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid) //get the asset from chaincode state
	if err != nil {
		return "", fmt.Errorf("failed to read prescription: %v", err)
	}
	if b64prescription == nil {
		return "", fmt.Errorf("%v does not exist in collection.", pid)
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
		return fmt.Errorf("Current user must be patient: %v", err)
	}

	//Check if client has access to this prescription
	err = checkClientAccess(ctx, pid)
	if err != nil {
		return fmt.Errorf("Client access check failed: %v", err)
	}

	// Get Access List
	b64access, err := ctx.GetStub().GetPrivateData(collectionAccessList, pid)
	if err != nil {
		return fmt.Errorf("failed to read prescription: %v", err)
	}
	if b64access == nil {
		return fmt.Errorf("%v does not exist in access list.", pid)
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

// ============================================================ //
// READ ALL PRESCRIPTIONS (for testing purpose only)
// DELETE ME when finished testing
// ============================================================ //
/*
// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllPrescriptions(ctx contractapi.TransactionContextInterface) ([]*Prescription, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(prescriptionCollection, "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Prescription
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Prescription
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (s *SmartContract) GetAllAccessLists(ctx contractapi.TransactionContextInterface) ([]*PrescriptionAccessList, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(accessListCollection, "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*PrescriptionAccessList
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset PrescriptionAccessList
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// ============================================================ //
// NON-ESSENTIAL: CREATE SAMPLE ASSETS
// ============================================================ //
func (s *SmartContract) CreateSamples(ctx contractapi.TransactionContextInterface) error {
	log.Printf("CREATING SAMPLE PRESCRIPTIONS")
	//Sample Prescription Data
	samples := []Prescription{
		{DrugBrand: "Paracetamol",
			DrugDoseSched:  "1 tablet / day",
			DrugPrice:      20.00,
			Id:             "prescription_1",
			Notes:          "",
			PatientName:    "Juan de la Cruz",
			PatientAddress: "Katipunan Ave, QC",
			PrescriberName: "Doctor Doctor",
			PrescriberNo:   "12345678"},

		{DrugBrand: "Advil",
			DrugDoseSched:  "1 pill / day",
			DrugPrice:      200.00,
			Id:             "prescription_2",
			Notes:          "do not buy generic brand",
			PatientName:    "Cruz de la Juan",
			PatientAddress: "Las Pinas City",
			PrescriberName: "Maria Clara, PhD",
			PrescriberNo:   "87654321"},

		{DrugBrand: "Vitamin Tablet",
			DrugDoseSched:  "2 tablets / day",
			DrugPrice:      150.00,
			Id:             "prescription_3",
			Notes:          "eat after meals",
			PatientName:    "Jan Delachoo",
			PatientAddress: "Loyola Heights, QC",
			PrescriberName: "Doctor A.",
			PrescriberNo:   "12312312"},
	}
	// Create Prescriptions & Access Lists
	for _, prescription := range samples {
		log.Printf("Prescription id: %v", prescription.Id)
		err := createPrescriptionFromAsset(ctx, &prescription)
		if err != nil {
			return err
		}
	}
	return nil
}
*/

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating thesis chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting thesis chaincode: %v", err)
	}
}
