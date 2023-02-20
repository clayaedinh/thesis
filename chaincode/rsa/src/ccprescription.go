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

func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface, b64prescription string) (string, error) {
	// Generate PID
	pid := genPrescriptionId()

	// Get requesting user
	currentUser, err := clientObscuredName(ctx)
	if err != nil {
		return "", err
	}

	// Verify if current user is a Patient
	err = ctx.GetClientIdentity().AssertAttributeValue("role", USER_PATIENT)
	if err != nil {
		return "", err
	}

	// Get Prescription Set
	prev, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return "", err
	}
	if prev != nil {
		return "", fmt.Errorf("cannot create prescription %v as it already exists", pid)
	}

	// Make map, consisting of all different encryptions of the same prescription
	pset := make(map[string]string)
	// Insert hash of the name of creating user
	pset[currentUser] = b64prescription

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

func (s *SmartContract) UpdatePrescription(ctx contractapi.TransactionContextInterface, pid string, b64pset string) error {

	// Verify if current user is a Doctor
	err := ctx.GetClientIdentity().AssertAttributeValue("role", USER_DOCTOR)
	if err != nil {
		return err
	}

	// Get Old Prescription Set
	oldpset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if oldpset == nil {
		return fmt.Errorf("cannot create prescription %v as it does not exist", pid)
	}

	// Confirm that user has access to this prescription in particular
	err = verifyClientAccess(ctx, string(oldpset))
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

func (s *SmartContract) SetfillPrescription(ctx contractapi.TransactionContextInterface, pid string, b64pset string) error {

	// Verify if current user is a Pharmacist
	err := ctx.GetClientIdentity().AssertAttributeValue("role", USER_PHARMACIST)
	if err != nil {
		return err
	}

	// Get Old Prescription Set
	oldpset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if oldpset == nil {
		return fmt.Errorf("cannot create prescription %v as it does not exist", pid)
	}

	// Confirm that user has access to this prescription in particular
	err = verifyClientAccess(ctx, string(oldpset))
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

func (s *SmartContract) SharePrescription(ctx contractapi.TransactionContextInterface, pid string, shareToUser string, b64prescription string) error {
	// Verify if current user is a Patient
	err := ctx.GetClientIdentity().AssertAttributeValue("role", USER_PATIENT)
	if err != nil {
		return err
	}

	// Get Prescription Set
	b64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if b64pset == nil {
		return fmt.Errorf("cannot create prescription %v as it does not exist", pid)
	}

	// Unpackage prescription set
	pset, err := unpackagePrescriptionSet(string(b64pset))
	if err != nil {
		return fmt.Errorf("failed to unpack prescription set: %v", err)
	}

	// Confirm that user has access to this prescription in particular
	err = verifyClientAccess(ctx, string(b64pset))
	if err != nil {
		return err
	}

	// Insert hash of the name of the user shared to
	(*pset)[shareToUser] = b64prescription

	// Repackage
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

func (s *SmartContract) ReadPrescription(ctx contractapi.TransactionContextInterface, pid string) (string, error) {
	// Get requesting user
	currentUser, err := clientObscuredName(ctx)
	if err != nil {
		return "", err
	}
	// Get Prescription Set
	b64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return "", fmt.Errorf("failed to read prescription: %v", err)
	}
	if b64pset == nil {
		return "", fmt.Errorf("no prescription set to read with given pid: %v", err)
	}

	// Unpackage Prescription Set
	pset, err := unpackagePrescriptionSet(string(b64pset))
	if err != nil {
		return "", err
	}

	// Return b64prescription if it has been encrypted for the given user
	b64prescription, exists := (*pset)[currentUser]

	if exists {
		return b64prescription, nil
	}

	return "", fmt.Errorf("given user does not have access to prescription")

}

func (s *SmartContract) DeletePrescription(ctx contractapi.TransactionContextInterface, pid string) error {

	// Verify if current user is a Patient
	err := ctx.GetClientIdentity().AssertAttributeValue("role", USER_PATIENT)
	if err != nil {
		return err
	}

	// Get Old Prescription Set
	oldpset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}
	if oldpset == nil {
		return fmt.Errorf("cannot create prescription %v as it does not exist", pid)
	}

	// Confirm that user has access to this prescription in particular
	err = verifyClientAccess(ctx, string(oldpset))
	if err != nil {
		return err
	}

	err = ctx.GetStub().DelPrivateData(collectionPrescription, pid)
	if err != nil {
		return fmt.Errorf("error in deleting prescription data: %v", err)
	}
	return nil
}

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
