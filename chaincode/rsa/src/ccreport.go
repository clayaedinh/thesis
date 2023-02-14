package src

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) RegisterReportReader(ctx contractapi.TransactionContextInterface, username string) error {
	//verify if real user by querying the key collection
	val, err := ctx.GetStub().GetPrivateData(collectionPubkeyRSA, username)
	if err != nil {
		return fmt.Errorf("Failed to verify whether user exists :%v", err)
	}
	if val == nil {
		return fmt.Errorf("Given user does not exist.")
	}

	err = ctx.GetStub().PutPrivateData(collectionReportReaders, username, []byte(username))
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) RemoveReportReader(ctx contractapi.TransactionContextInterface, username string) error {
	err := ctx.GetStub().DelPrivateData(collectionReportReaders, username)
	if err != nil {
		return fmt.Errorf("Error in removing report reader: %v", err)
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
	rgob, err := encodeStringSlice(readers)
	if err != nil {
		return "", err
	}
	b64gob := base64.StdEncoding.EncodeToString(rgob)

	return b64gob, nil
}

func (s *SmartContract) ReportGenerate(ctx contractapi.TransactionContextInterface, pid string, b64gob string) error {
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
	err = ctx.GetStub().PutPrivateData(collectionReports, pid, newgob)
	if err != nil {
		return fmt.Errorf("Failed to add prescription to private data: %v", err)
	}
	return nil
}

func (s *SmartContract) GetAllPrescriptions(ctx contractapi.TransactionContextInterface) (string, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionReports, "", "")
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var assets [][]byte
	for resultsIterator.HasNext() {
		asset, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}
		assets = append(assets, asset.Value)
	}

	//I have no idea if this will work
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(assets)
	if err != nil {
		return "", fmt.Errorf("Failed to package prescriptions for transport: %v", err)
	}

	//base64 it
	b64all := base64.StdEncoding.EncodeToString(buf.Bytes())

	return b64all, nil
}

/*
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

*/
