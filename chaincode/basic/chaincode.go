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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ============================================================ //
// DATA
// ============================================================ //

const prescriptionCollection = "prescriptionCollection"
const accessListCollection = "accessListCollection"

// SmartContract provides functions for managing Assets such as Prescription

type SmartContract struct {
	contractapi.Contract
}

// Structure for Prescription Data
type Prescription struct {
	DrugBrand      string  `json:"DrugBrand"`
	DrugDoseSched  string  `json:"DrugDoseSched"`
	DrugPrice      float64 `json:"DrugPrice"`
	Id             string  `json:"Id"`
	Notes          string  `json:"Notes"`
	PatientName    string  `json:"PatientName"`
	PatientAddress string  `json:"PatientAddress"`
	PrescriberName string  `json:"PrescriberName"`
	PrescriberNo   string  `json:"PrescriberNo"`
	FilledAmount   string  `json:"FilledAmount"`
}

type PrescriptionAccessList struct {
	PrescriptionId string   `json:"prescriptionId"`
	UserIds        []string `json:"userIds"`
}

// String Constants for User Roles
const (
	USER_DOCTOR     string = "DOCTOR"
	USER_PATIENT           = "PATIENT"
	USER_PHARMACIST        = "PHARMA"
)

// ============================================================ //
// HELPER FUNCTIONS
// ============================================================ //

func submittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
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

// ============================================================ //
// CLIENT ACCESS CHECK
// Verifies that the current client identity has access
// ============================================================ //
func checkClientAccess(ctx contractapi.TransactionContextInterface, prescriptionId string) error {
	// Get Access List
	accessJSON, err := ctx.GetStub().GetPrivateData(accessListCollection, prescriptionId)
	if err != nil {
		return fmt.Errorf("failed to read access list: %v", err)
	}
	if accessJSON == nil {
		return fmt.Errorf("%v does not exist in collection %v", prescriptionId, accessListCollection)
	}

	// Unmarshall JSON
	var accessList *PrescriptionAccessList
	err = json.Unmarshal(accessJSON, &accessList)
	if err != nil {
		return fmt.Errorf("failed to unmarshal access list: %v", err)
	}

	// GET CLIENT ID
	clientId, err := submittingClientIdentity(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get client id: %v", err)
	}

	//CONFIRM WITH ACCESS RECORD IF USER HAS ACCESS
	var matchingId bool = false
	for _, element := range accessList.UserIds {
		if element == clientId {
			matchingId = true
			break
		}
	}
	if matchingId == false {
		return fmt.Errorf("Permission Denied. Given User does not have access: %v", clientId)
	}
	return nil
}

// ============================================================ //
// GET MY ID
// ============================================================ //
func (s *SmartContract) GetMyID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientId, err := submittingClientIdentity(ctx)
	if err != nil {
		return "", fmt.Errorf("Failed to get client id: %v", err)
	}
	return clientId, nil
}

// ============================================================ //
// SAVE PRESCRIPTION
// Private method
// Used to save prescription data (newly created or update) to ledger
// ============================================================ //
func savePrescription(ctx contractapi.TransactionContextInterface, asset *Prescription) error {
	// Marshall Prescription Asset
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("Failed to marshall prescription asset to JSON: %v", err)
	}

	// Add Asset to Chain
	err = ctx.GetStub().PutPrivateData(prescriptionCollection, asset.Id, assetJSON)
	if err != nil {
		return fmt.Errorf("Failed to add prescription %v to private data: %v", asset.Id, err)
	}

	return nil
}

// ============================================================ //
// CREATE PRESCRIPTION FROM ASSET
// Private method
// Creates prescription with parameter as Prescription struct
// ============================================================ //
func createPrescriptionFromAsset(ctx contractapi.TransactionContextInterface, prescription *Prescription) error {
	//Save new prescription data
	err := savePrescription(ctx, prescription)
	if err != nil {
		return fmt.Errorf("Error saving prescription: %v", err)
	}

	// Get Client Identity
	clientId, err := submittingClientIdentity(ctx)
	if err != nil {
		return fmt.Errorf("Failed to get client identity: %v", err)
	}

	// Create Access List for Prescription
	accessList := PrescriptionAccessList{
		PrescriptionId: prescription.Id,
		UserIds:        []string{clientId},
	}
	// Marshall Access List
	accessJSON, err := json.Marshal(accessList)
	if err != nil {
		return fmt.Errorf("failed to marshal prescription access list into JSON: %v", err)
	}

	// Add Access List to Chain
	err = ctx.GetStub().PutPrivateData(accessListCollection, prescription.Id, accessJSON)
	if err != nil {
		return fmt.Errorf("failed to add prescription access list %v to private data: %v", prescription.Id, err)
	}

	return nil
}

// ============================================================ //
// CREATE PRESCRIPTION
// May only be called by doctors
// ============================================================ //
func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface,
	drugbrand string, drugdosesched string, drugprice float64, id string, notes string,
	patientname string, patientaddress string, prescribername string, prescriberno string) error {

	// Verify if no prescription already exists with the given id
	prescriptionJSON, err := ctx.GetStub().GetPrivateData(prescriptionCollection, id) //get the asset from chaincode state
	if err != nil {
		return fmt.Errorf("Failed to verify if prescription Id is already used: %v", err)
	}
	if prescriptionJSON != nil {
		return fmt.Errorf("Prescription creation failed: ID already in use :")
	}

	// Verify if current user is a Doctor
	err = ctx.GetClientIdentity().AssertAttributeValue("role", USER_DOCTOR)
	if err != nil {
		return err
	}

	//Check if client has access to this prescription
	err = checkClientAccess(ctx, id)
	if err != nil {
		return fmt.Errorf("Client access check failed: %v", err)
	}

	// Create Prescription Asset
	prescription := Prescription{
		DrugBrand:      drugbrand,
		DrugDoseSched:  drugdosesched,
		DrugPrice:      drugprice,
		Id:             id,
		Notes:          notes,
		PatientName:    patientname,
		PatientAddress: patientaddress,
		PrescriberName: prescribername,
		PrescriberNo:   prescriberno,
		FilledAmount:   "None",
	}

	return createPrescriptionFromAsset(ctx, &prescription)
}

// ============================================================ //
// UPDATE PRESCRIPTION
// May only be called by doctors
// ============================================================ //
func (s *SmartContract) UpdatePrescription(ctx contractapi.TransactionContextInterface,
	drugbrand string, drugdosesched string, drugprice float64, id string, notes string,
	patientname string, patientaddress string, prescribername string, prescriberno string) error {

	// Verify if prescription with given id exists, so it can be updated
	prescriptionJSON, err := ctx.GetStub().GetPrivateData(prescriptionCollection, id) //get the asset from chaincode state
	if err != nil {
		return fmt.Errorf("Failed to verify if prescription with provided Id exists: %v", err)
	}
	if prescriptionJSON == nil {
		return fmt.Errorf("Prescription update failed: no prescription exists with given Id")
	}

	// Verify if current user is a Doctor
	err = ctx.GetClientIdentity().AssertAttributeValue("role", USER_DOCTOR)
	if err != nil {
		return fmt.Errorf("Current user must be doctor: %v", err)
	}

	//Check if client has access to this prescription
	err = checkClientAccess(ctx, id)
	if err != nil {
		return fmt.Errorf("Client access check failed: %v", err)
	}

	//Unmarshal prescription JSON (so that FilledAmount may be obtained)
	var old_prescription *Prescription
	err = json.Unmarshal(prescriptionJSON, &old_prescription)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Prescription Record JSON: %v", err)
	}

	// Create Prescription Asset
	prescription := Prescription{
		DrugBrand:      drugbrand,
		DrugDoseSched:  drugdosesched,
		DrugPrice:      drugprice,
		Id:             id,
		Notes:          notes,
		PatientName:    patientname,
		PatientAddress: patientaddress,
		PrescriberName: prescribername,
		PrescriberNo:   prescriberno,
		FilledAmount:   old_prescription.FilledAmount,
	}

	return savePrescription(ctx, &prescription)
}

// ============================================================ //
// SETFILL PRESCRIPTION
// May only be called by pharmacists
// ============================================================ //
func (s *SmartContract) SetFillPrescription(ctx contractapi.TransactionContextInterface, prescriptionId string, newfill string) error {
	// Verify if prescription with given id exists, so it can be updated
	prescriptionJSON, err := ctx.GetStub().GetPrivateData(prescriptionCollection, prescriptionId) //get the asset from chaincode state
	if err != nil {
		return fmt.Errorf("Failed to verify if prescription with provided Id exists: %v", err)
	}
	if prescriptionJSON == nil {
		return fmt.Errorf("Prescription setfill failed: no prescription exists with given Id")
	}

	// Verify if current user is a Pharmacist
	err = ctx.GetClientIdentity().AssertAttributeValue("role", USER_PHARMACIST)
	if err != nil {
		return fmt.Errorf("Current user must be pharmacist: %v", err)
	}

	//Check if client has access to this prescription
	err = checkClientAccess(ctx, prescriptionId)
	if err != nil {
		return fmt.Errorf("Client access check failed: %v", err)
	}

	//Unmarshal old prescription JSON
	var old *Prescription
	err = json.Unmarshal(prescriptionJSON, &old)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Prescription Record JSON: %v", err)
	}

	// Create Prescription Asset
	prescription := Prescription{
		DrugBrand:      old.DrugBrand,
		DrugDoseSched:  old.DrugDoseSched,
		DrugPrice:      old.DrugPrice,
		Id:             old.Id,
		Notes:          old.Notes,
		PatientName:    old.PatientName,
		PatientAddress: old.PatientAddress,
		PrescriberName: old.PrescriberName,
		PrescriberNo:   old.PrescriberNo,
		FilledAmount:   newfill,
	}

	return savePrescription(ctx, &prescription)
}

// ============================================================ //
// Delete PRESCRIPTION
// ============================================================ //
func (s *SmartContract) DeletePrescription(ctx contractapi.TransactionContextInterface, prescriptionId string) error {

	//Check if client has access to this prescription
	err := checkClientAccess(ctx, prescriptionId)
	if err != nil {
		return fmt.Errorf("Client access check failed: %v", err)
	}

	//delete prescription data proper
	prescriptionJSON, err := ctx.GetStub().GetPrivateData(prescriptionCollection, prescriptionId)
	if err != nil {
		return fmt.Errorf("Failed to verify if prescription with provided Id exists: %v", err)
	}
	if prescriptionJSON == nil {
		return fmt.Errorf("Prescription deletion failed: no prescription exists with given Id")
	}
	err = ctx.GetStub().DelPrivateData(prescriptionCollection, prescriptionId)
	if err != nil {
		return fmt.Errorf("Error in deleting prescription data: %v", err)
	}

	//delete access lists
	accessJSON, err := ctx.GetStub().GetPrivateData(accessListCollection, prescriptionId)
	if err != nil {
		return fmt.Errorf("Failed to verify if access list with provided Id exists: %v", err)
	}
	if accessJSON == nil {
		return fmt.Errorf("Prescription deletion failed: no access list exists with given Id")
	}
	err = ctx.GetStub().DelPrivateData(accessListCollection, prescriptionId)
	if err != nil {
		return fmt.Errorf("Error in deleting access list data: %v", err)
	}
	return nil
}

// ============================================================ //
// READ PRESCRIPTION
// ============================================================ //

func (s *SmartContract) ReadPrescription(ctx contractapi.TransactionContextInterface, prescriptionId string) (*Prescription, error) {

	//Check if client has access to this prescription
	err := checkClientAccess(ctx, prescriptionId)
	if err != nil {
		return nil, fmt.Errorf("Client access check failed: %v", err)
	}

	//READ PRESCRIPTION AT ID
	log.Printf("ReadAsset: collection %v, Id %v", prescriptionCollection, prescriptionId)
	prescriptionJSON, err := ctx.GetStub().GetPrivateData(prescriptionCollection, prescriptionId) //get the asset from chaincode state
	if err != nil {
		return nil, fmt.Errorf("failed to read prescription: %v", err)
	}

	//No Prescription found, return empty response
	if prescriptionJSON == nil {
		return nil, fmt.Errorf("%v does not exist in collection %v", prescriptionId, prescriptionCollection)
	}

	//CONVERT TO READABLE FORMAT
	var prescription *Prescription
	err = json.Unmarshal(prescriptionJSON, &prescription)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Prescription Record JSON: %v", err)
	}

	return prescription, nil
}

// ============================================================ //
// SHARE PRESCRIPTION
// ============================================================ //

func (s *SmartContract) SharePrescription(ctx contractapi.TransactionContextInterface, prescriptionId string, newUserId string) error {

	// Verify if current user is a Patient
	err := ctx.GetClientIdentity().AssertAttributeValue("role", USER_PATIENT)
	if err != nil {
		return fmt.Errorf("Current user must be patient: %v", err)
	}

	//Check if client has access to this prescription
	err = checkClientAccess(ctx, prescriptionId)
	if err != nil {
		return fmt.Errorf("Client access check failed: %v", err)
	}

	// Get Access List
	accessJSON, err := ctx.GetStub().GetPrivateData(accessListCollection, prescriptionId)
	if err != nil {
		return fmt.Errorf("failed to read prescription: %v", err)
	}
	if accessJSON == nil {
		return fmt.Errorf("%v does not exist in collection %v", prescriptionId, accessListCollection)
	}

	// Unmarshall JSON
	var accessList *PrescriptionAccessList
	err = json.Unmarshal(accessJSON, &accessList)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Transfer asset in private data collection to new owner
	accessList.UserIds = append(accessList.UserIds, newUserId)

	//Marshal data
	assetJSONasBytes, err := json.Marshal(accessList)
	if err != nil {
		return fmt.Errorf("failed marshalling asset %v: %v", prescriptionId, err)
	}

	//Reinsert data
	err = ctx.GetStub().PutPrivateData(accessListCollection, prescriptionId, assetJSONasBytes) //rewrite the asset
	if err != nil {
		return err
	}
	return nil
}

// ============================================================ //
// READ ALL PRESCRIPTIONS (for testing purpose only)
// DELETE ME when finished testing
// ============================================================ //

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

func main() {
	assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating thesis chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting thesis chaincode: %v", err)
	}
}
