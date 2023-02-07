package src

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) CreatePrescriptionSimple(ctx contractapi.TransactionContextInterface, prescription_id string, b64encrypted string) error {
	// Decode base64-encoded encrypted data
	encrypted, err := base64.StdEncoding.DecodeString(b64encrypted)
	if err != nil {
		return fmt.Errorf("Base64 decoding of RSA pubkey failed: %v", err)
	}

	// Add Prescription
	err = ctx.GetStub().PutPrivateData(collectionPrescription, prescription_id, encrypted)
	if err != nil {
		return fmt.Errorf("failed to add prescription access list %v to private data: %v", prescription_id, err)
	}
	return nil
}
