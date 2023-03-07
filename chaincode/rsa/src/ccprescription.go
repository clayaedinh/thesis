package src

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*
ACCESS CONTROLS

Patient - CreatePrescription, SharePrescription, Delete Prescription
Doctor - Update Prescription
Pharmacist - SetFill Prescription
All - Read Prescription
*/
// ============================================================ //
// Create Prescription
// ============================================================ //
func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface, b64prescription string) (string, error) {
	// Generate ID, and check if no prescription already exists with the given id
	pid := genPrescriptionId()
	prev, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return "", err
	}
	if prev != nil {
		return "", fmt.Errorf("cannot create prescription %v as it already exists", pid)
	}
	// Get requesting user
	currentUser := clientObscuredName(ctx)

	// Make map pset, consisting of all different encryptions of the same prescription
	pset := make(map[string]string)
	// Insert given b64-encoded & encrypted prescription into the pset where key = current user's name
	pset[currentUser] = b64prescription
	// Package the pset
	b64pset, err := packagePrescriptionSet(&pset)
	if err != nil {
		return "", err
	}
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64pset))
	if err != nil {
		return "", fmt.Errorf("failed to add prescription to private data: %v", err)
	}
	return pid, nil
}

// ============================================================ //
// Read Prescription
// ============================================================ //
func (s *SmartContract) ReadPrescription(ctx contractapi.TransactionContextInterface, pid string) (string, error) {
	// Get Prescription Set
	b64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return "", fmt.Errorf("failed to read prescription: %v", err)
	}
	if b64pset == nil {
		return "", fmt.Errorf("no prescription set to read with given pid: %v", err)
	}
	obscureName := clientObscuredName(ctx)
	pset, err := unpackageAndCheckAccess(ctx, string(b64pset), obscureName)
	if err != nil {
		return "", err
	}
	return (*pset)[obscureName], nil
}

// ============================================================ //
// Share Prescription
// ============================================================ //
func (s *SmartContract) SharePrescription(ctx contractapi.TransactionContextInterface, pid string, shareToUser string, b64prescription string) error {
	// Get Prescription Set
	b64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if b64pset == nil {
		return fmt.Errorf("cannot create prescription %v as it does not exist", pid)
	}
	// Unpackage prescription set
	pset, err := unpackageAndCheckAccess(ctx, string(b64pset), clientObscuredName(ctx))
	if err != nil {
		return err
	}
	// Insert prescription with key=user shared to
	(*pset)[shareToUser] = b64prescription
	// Repackage the prescription set
	b64updatedpset, err := packagePrescriptionSet(pset)
	if err != nil {
		return err
	}
	// Save to Private Data
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64updatedpset))
	if err != nil {
		return fmt.Errorf("failed to add prescription to private data: %v", err)
	}
	return nil
}

// ============================================================ //
// Users this prescriptin is Shared TO
// ============================================================ //
func (s *SmartContract) PrescriptionSharedTo(ctx contractapi.TransactionContextInterface, pid string) (string, error) {
	// Get Prescription Set
	b64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return "", fmt.Errorf("failed to read prescription: %v", err)
	}
	if b64pset == nil {
		return "", fmt.Errorf("no prescription with given pid: %v", err)
	}
	// Unpackage the data
	pset, err := unpackagePrescriptionSet(string(b64pset))
	if err != nil {
		return "", fmt.Errorf("failed to unpack prescription set: %v", err)
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

// ============================================================ //
// Update Prescription
// ============================================================ //
func (s *SmartContract) UpdatePrescription(ctx contractapi.TransactionContextInterface, pid string, b64pset string) error {
	// Verify if a prescription already exists with the given id
	oldb64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if oldb64pset == nil {
		return fmt.Errorf("cannot create prescription %v as it does not exist", pid)
	}
	// Confirm that user has access to this prescription in particular
	_, err = unpackageAndCheckAccess(ctx, string(oldb64pset), clientObscuredName(ctx))
	if err != nil {
		return err
	}
	//Upload the update
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64pset))
	if err != nil {
		return fmt.Errorf("failed to add prescription to private data: %v", err)
	}
	return nil
}

// ============================================================ //
// Setfill Prescription
// ============================================================ //
func (s *SmartContract) SetfillPrescription(ctx contractapi.TransactionContextInterface, pid string, b64pset string) error {
	// Verify if a prescription already exists with the given id
	oldb64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if oldb64pset == nil {
		return fmt.Errorf("cannot create prescription %v as it does not exist", pid)
	}

	// Confirm that user has access to this prescription in particular
	_, err = unpackageAndCheckAccess(ctx, string(oldb64pset), clientObscuredName(ctx))
	if err != nil {
		return err
	}
	//Upload the update
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64pset))
	if err != nil {
		return fmt.Errorf("failed to add prescription to private data: %v", err)
	}
	return nil
}

// ============================================================ //
// Delete Prescription
// ============================================================ //
func (s *SmartContract) DeletePrescription(ctx contractapi.TransactionContextInterface, pid string) error {
	// Get Old Prescription Set
	oldb64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if oldb64pset == nil {
		return fmt.Errorf("cannot delete prescription %v as it does not exist", pid)
	}
	// Confirm that user has access to this prescription in particular
	_, err = unpackageAndCheckAccess(ctx, string(oldb64pset), clientObscuredName(ctx))
	if err != nil {
		return err
	}
	// Delete data
	err = ctx.GetStub().DelPrivateData(collectionPrescription, pid)
	if err != nil {
		return fmt.Errorf("error in deleting prescription data: %v", err)
	}
	return nil
}
