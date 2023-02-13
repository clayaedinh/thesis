package src

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TODO: Add Access Controls.

func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface, pid string, usernameHash string, pdata string) error {
	// Get Prescription Set
	prev, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if prev != nil {
		return fmt.Errorf("Cannot create prescription %v as it already exists", pid)
	}

	// Make map, consisting of all different encryptions of the same prescription
	pset := make(map[string]string)

	// Insert hash of the name of creating user
	pset[usernameHash] = pdata

	// Encode to gob for easier transport
	gobdata, err := encodePrescriptionSet(&pset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, gobdata)
	if err != nil {
		return fmt.Errorf("Failed to add prescription to private data: %v", err)
	}
	return nil
}

func (s *SmartContract) UpdatePrescription(ctx contractapi.TransactionContextInterface, pid string, b64gob string) error {
	//decode
	newgob, err := base64.StdEncoding.DecodeString(b64gob)
	if err != nil {
		return fmt.Errorf("Invalid Base64 encoding: %v", err)
	}

	// Get Prescription Set
	prev, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if prev == nil {
		return fmt.Errorf("Cannot update prescription %v as it does not exist", pid)
	}

	// Decode gob to verify if it is valid daata
	_, err = decodePrescriptionSet(newgob)
	if err != nil {
		return fmt.Errorf("Invalid update data: %v", err)
	}

	//Upload the update
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, newgob)
	if err != nil {
		return fmt.Errorf("Failed to add prescription to private data: %v", err)
	}
	return nil
}

func (s *SmartContract) SharePrescription(ctx contractapi.TransactionContextInterface, pid string, usernameHash string, pdata string) error {
	// Get Prescription Set
	pgob, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if pgob == nil {
		return fmt.Errorf("Cannot create prescription %v as it does not exist", pid)
	}

	// Decode gob
	pset, err := decodePrescriptionSet(pgob)
	if err != nil {
		return fmt.Errorf("Failed to unpack prescription set: %v", err)
	}

	// Insert hash of the name of creating user
	pset[usernameHash] = pdata

	// Encode to gob
	gobdata, err := encodePrescriptionSet(&pset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, gobdata)
	if err != nil {
		return fmt.Errorf("Failed to add prescription to private data: %v", err)
	}
	return nil
}

func (s *SmartContract) ReadPrescription(ctx contractapi.TransactionContextInterface, pid string, usernameHash string) (string, error) {
	// Get Prescription Set
	pgob, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return "", fmt.Errorf("Failed to read prescription: %v", err)
	}
	if pgob == nil {
		return "", fmt.Errorf("No prescription set to read with given pid: %v", err)
	}

	// Decode gob
	pset, err := decodePrescriptionSet(pgob)
	if err != nil {
		return "", fmt.Errorf("Failed to unpack prescription set: %v", err)
	}

	// Return pdata if it has been encrypted for the given user
	pdata, exists := pset[usernameHash]

	if exists {
		return pdata, nil
	}

	return "", fmt.Errorf("Given user does not have permission to this prescription")

}

func (s *SmartContract) DeletePrescription(ctx contractapi.TransactionContextInterface, pid string) error {
	err := ctx.GetStub().DelPrivateData(collectionPrescription, pid)
	if err != nil {
		return fmt.Errorf("Error in deleting prescription data: %v", err)
	}
	return nil
}

func (s *SmartContract) PrescriptionSharedTo(ctx contractapi.TransactionContextInterface, pid string) (string, error) {
	// Get Prescription Set
	pgob, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return "", fmt.Errorf("Failed to read prescription: %v", err)
	}
	if pgob == nil {
		return "", fmt.Errorf("No prescription with given pid: %v", err)
	}
	// Decode gob
	pset, err := decodePrescriptionSet(pgob)
	if err != nil {
		return "", fmt.Errorf("Failed to unpack prescription set: %v", err)
	}
	// Get list of keys in map
	mapkeys := make([]string, len(pset))
	i := 0
	for key := range pset {
		mapkeys[i] = key
		i++
	}
	// Encode this list of map keys
	keygob, err := encodeStringSlice(mapkeys)
	if err != nil {
		return "", err
	}
	b64gob := base64.StdEncoding.EncodeToString(keygob)
	return b64gob, nil

}

/*
func (s *SmartContract) PrescriptionAllEncryptions(ctx contractapi.TransactionContextInterface, pid string) ([]byte, error) {
	// Get Prescription Set
	pgob, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return nil, fmt.Errorf("Failed to read prescription: %v", err)
	}
	if pgob == nil {
		return nil, fmt.Errorf("No prescription with given pid: %v", err)
	}
	b64gob := base64.StdEncoding.EncodeToString(pgob)
	return []byte(b64gob), nil
}
*/
