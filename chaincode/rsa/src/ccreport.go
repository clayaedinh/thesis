package src

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) RegisterReportReader(ctx contractapi.TransactionContextInterface, username string) error {
	err := ctx.GetStub().PutPrivateData(collectionReportReaders, username, []byte(username))
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) GetAllReportReaders(ctx contractapi.TransactionContextInterface) ([]string, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionReportReaders, "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var readers []string
	for resultsIterator.HasNext() {
		reader, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		readers = append(readers, string(reader.Value))
	}

	return readers, nil
}

func (s *SmartContract) GetAllPrescriptions(ctx contractapi.TransactionContextInterface) ([]byte, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionReportReaders, "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets [][]byte
	for resultsIterator.HasNext() {
		asset, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset.Value)
	}

	//I have no idea if this will work
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(assets)
	if err != nil {
		return nil, fmt.Errorf("Failed to package prescriptions for transport: %v", err)
	}

	return buf.Bytes(), nil
}
