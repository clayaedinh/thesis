package src

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) StoreUserRSAPubkey(ctx contractapi.TransactionContextInterface, username string, b64pubkey string) error {
	pubkey, err := base64.StdEncoding.DecodeString(b64pubkey)
	if err != nil {
		return fmt.Errorf("base64 decoding of RSA pubkey failed: %v", err)
	}
	err = ctx.GetStub().PutPrivateData(collectionPubkeyRSA, username, pubkey)
	if err != nil {
		return fmt.Errorf("failed to store user RSA Pubkey: %v", err)
	}
	return nil
}

func (s *SmartContract) RetrieveUserRSAPubkey(ctx contractapi.TransactionContextInterface, username string) (string, error) {
	pubkey, err := ctx.GetStub().GetPrivateData(collectionPubkeyRSA, username)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user RSA Pubkey: %v", err)
	}
	if pubkey == nil {
		return "", fmt.Errorf("pubkey for user '%v' does not exist", username)
	}
	b64pubkey := base64.StdEncoding.EncodeToString(pubkey)
	return b64pubkey, nil
}
