package src

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) StoreUserRSAPubkey(ctx contractapi.TransactionContextInterface, usernameHash string, b64pubkey string) error {
	pubkey, err := base64.StdEncoding.DecodeString(b64pubkey)
	if err != nil {
		return fmt.Errorf("Base64 decoding of RSA pubkey failed: %v", err)
	}
	err = ctx.GetStub().PutPrivateData(collectionPubkeyRSA, usernameHash, pubkey)
	if err != nil {
		return fmt.Errorf("Failed to store user RSA Pubkey: %v", err)
	}
	return nil
}

func (s *SmartContract) RetrieveUserRSAPubkey(ctx contractapi.TransactionContextInterface, usernameHash string) (string, error) {
	pubkey, err := ctx.GetStub().GetPrivateData(collectionPubkeyRSA, usernameHash)
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve user RSA Pubkey: %v", err)
	}
	if pubkey == nil {
		return "", fmt.Errorf("Pubkey for user '%v' does not exist.", usernameHash)
	}
	b64pubkey := base64.StdEncoding.EncodeToString(pubkey)
	return b64pubkey, nil
}
