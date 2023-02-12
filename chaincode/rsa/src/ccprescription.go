package src

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface, tag string, pdata string) error {
	// Add Prescription to Chain
	err := ctx.GetStub().PutPrivateData(collectionPrescription, tag, []byte(pdata))
	if err != nil {
		return fmt.Errorf("Failed to add prescription to private data: %v", err)
	}
	return nil
}

func (s *SmartContract) ReadPrescription(ctx contractapi.TransactionContextInterface, tag string) (string, error) {
	pdata, err := ctx.GetStub().GetPrivateData(collectionPrescription, tag)
	if err != nil {
		return "", fmt.Errorf("Failed to read presscription: %v", err)
	}
	if pdata == nil {
		return "", fmt.Errorf("No prescription to read with given tag: %v", err)
	}
	return string(pdata), nil
}
