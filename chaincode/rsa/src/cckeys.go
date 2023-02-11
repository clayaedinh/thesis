package src

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) StoreUserRSAPubkey(ctx contractapi.TransactionContextInterface, username string, b64pubkey string) error {
	pubkey, err := base64.StdEncoding.DecodeString(b64pubkey)
	if err != nil {
		return fmt.Errorf("Base64 decoding of RSA pubkey failed: %v", err)
	}
	err = ctx.GetStub().PutPrivateData(collectionPubkeyRSA, username, pubkey)
	if err != nil {
		return fmt.Errorf("Failed to store user RSA Pubkey: %v", err)
	}
	return nil
}

func (s *SmartContract) RetrieveUserRSAPubkey(ctx contractapi.TransactionContextInterface, username string) (string, error) {
	pubkey, err := ctx.GetStub().GetPrivateData(collectionPubkeyRSA, username)
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve user RSA Pubkey: %v", err)
	}
	if pubkey == nil {
		return "", fmt.Errorf("Pubkey for user '%v' does not exist.", username)
	}
	b64pubkey := base64.StdEncoding.EncodeToString(pubkey)
	return b64pubkey, nil
}

func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface, tag string, pdata string) error {

	decrypted, err := base64.StdEncoding.DecodeString(pdata)
	if err != nil {
		return fmt.Errorf("Base64 failed to decrypt given prescription: %v", err)
	}

	// Add Prescription to Chain
	err = ctx.GetStub().PutPrivateData(collectionPrescription, tag, decrypted)
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
	return base64.StdEncoding.EncodeToString(pdata), nil
}
