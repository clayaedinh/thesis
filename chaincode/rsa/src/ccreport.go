package src

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) RegisterReportReader(ctx contractapi.TransactionContextInterface, obscuredName string) error {
	//verify if real user by querying the key collection
	val, err := ctx.GetStub().GetPrivateData(collectionPubkeyRSA, obscuredName)
	if err != nil {
		return fmt.Errorf("failed to verify whether user exists :%v", err)
	}
	if val == nil {
		return fmt.Errorf("given user does not exist")
	}

	//Add the user to the report reading collection
	err = ctx.GetStub().PutPrivateData(collectionReportReaders, obscuredName, []byte(obscuredName))
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) RemoveReportReader(ctx contractapi.TransactionContextInterface, username string) error {
	// remove the given report reader
	err := ctx.GetStub().DelPrivateData(collectionReportReaders, username)
	if err != nil {
		return fmt.Errorf("error in removing report reader: %v", err)
	}
	return nil
}

func (s *SmartContract) GetAllReportReaders(ctx contractapi.TransactionContextInterface) (string, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionReportReaders, "", "")
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var readers []string
	for resultsIterator.HasNext() {
		reader, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}
		readers = append(readers, string(reader.Value))
	}

	// Encode this list of map keys
	b64slice, err := packageStringSlice(&readers)
	if err != nil {
		return "", err
	}

	return b64slice, nil
}

func (s *SmartContract) ReportGenerate(ctx contractapi.TransactionContextInterface, pid string, b64reports string) error {
	// unpack the b64 reports
	reportset, err := unpackagePrescriptionSet(b64reports)
	if err != nil {
		return err
	}

	// get the current pset for the given pid
	b64pset, err := ctx.GetStub().GetPrivateData(collectionPrescription, pid)
	if err != nil {
		return err
	}

	// Unpackage pset
	pset, err := unpackagePrescriptionSet(string(b64pset))
	if err != nil {
		return fmt.Errorf("failed to unpack prescription set: %v", err)
	}

	// Iterate over the report set and add each prescription inside to the pset
	for key, value := range *reportset {
		(*pset)[key] = value
	}

	// repackage pset
	b64psetNew, err := packagePrescriptionSet(pset)
	if err != nil {
		return err
	}

	// update private data
	err = ctx.GetStub().PutPrivateData(collectionPrescription, pid, []byte(b64psetNew))
	if err != nil {
		return fmt.Errorf("failed to add prescription to private data: %v", err)
	}
	return nil

}

func (s *SmartContract) GetPrescriptionReport(ctx contractapi.TransactionContextInterface, username string) (string, error) {
	// Create Iterator
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionPrescription, "", "")
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	// Create set of reports (all prescriptions, encrypted for given user)
	reportset := make(map[string]string)
	for resultsIterator.HasNext() {
		// get each prescription set
		b64pset, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}
		// unpackage the prescription set
		pset, err := unpackagePrescriptionSet(string(b64pset.Value))
		if err != nil {
			return "", err
		}
		// get the prescription with the given username
		reportset[string(b64pset.Key)] = (*pset)[username]
	}
	// package the reportset
	b64reports, err := packagePrescriptionSet(&reportset)
	if err != nil {
		return "", err
	}
	return b64reports, nil
}
