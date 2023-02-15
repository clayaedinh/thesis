package src

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TODO: Add Access Controls.

func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface, pid string, obscuredName string, b64prescription string) error {
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
	pset[obscuredName] = b64prescription

	b64pset, err := packagePrescriptionSet(&pset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64pset))
	if err != nil {
		return fmt.Errorf("Failed to add prescription to private data: %v", err)
	}
	return nil
}

func (s *SmartContract) UpdatePrescription(ctx contractapi.TransactionContextInterface, pid string, b64pset string) error {
	//Upload the update
	err := ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64pset))
	if err != nil {
		return fmt.Errorf("Failed to add prescription to private data: %v", err)
	}
	return nil
}

func (s *SmartContract) SharePrescription(ctx contractapi.TransactionContextInterface, pid string, obscuredName string, b64prescription string) error {
	// Get Prescription Set
	b64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if b64pset == nil {
		return fmt.Errorf("Cannot create prescription %v as it does not exist", pid)
	}

	// Unpackage gob
	pset, err := unpackagePrescriptionSet(string(b64pset))
	if err != nil {
		return fmt.Errorf("Failed to unpack prescription set: %v", err)
	}

	// Insert hash of the name of creating user
	(*pset)[obscuredName] = b64prescription

	// Repackage
	b64updatedpset, err := packagePrescriptionSet(pset)
	if err != nil {
		return err
	}
	// Save to Private Data
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64updatedpset))
	if err != nil {
		return fmt.Errorf("Failed to add prescription to private data: %v", err)
	}
	return nil
}

func (s *SmartContract) ReadPrescription(ctx contractapi.TransactionContextInterface, pid string, obscuredName string) (string, error) {
	// Get Prescription Set
	b64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return "", fmt.Errorf("Failed to read prescription: %v", err)
	}
	if b64pset == nil {
		return "", fmt.Errorf("No prescription set to read with given pid: %v", err)
	}

	// Unpackage Prescription Set
	pset, err := unpackagePrescriptionSet(string(b64pset))
	if err != nil {
		return "", err
	}

	// Return b64prescription if it has been encrypted for the given user
	b64prescription, exists := (*pset)[obscuredName]

	if exists {
		return b64prescription, nil
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
	b64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return "", fmt.Errorf("Failed to read prescription: %v", err)
	}
	if b64pset == nil {
		return "", fmt.Errorf("No prescription with given pid: %v", err)
	}
	// Unpackage the data
	pset, err := unpackagePrescriptionSet(string(b64pset))
	if err != nil {
		return "", fmt.Errorf("Failed to unpack prescription set: %v", err)
	}
	// Get list of keys in map
	strslice := make([]string, len(*pset))
	i := 0
	for key := range *pset {
		strslice[i] = key
		i++
	}
	// Encode this list of map keys
	b64slice, err := packageStringSlice(&strslice)
	if err != nil {
		return "", err
	}
	return b64slice, nil

}
