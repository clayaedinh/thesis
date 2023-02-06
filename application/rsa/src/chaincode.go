package src

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func ChainStoreUserPubkey(contract *client.Contract, username string, pubkey []byte) error {
	b64pubkey := base64.StdEncoding.EncodeToString(pubkey)
	_, err := contract.SubmitTransaction("StoreUserRSAPubkey", username, b64pubkey)
	if err != nil {
		return err
	}
	return nil
}

func ChainRetrieveUserPubkey(contract *client.Contract, username string) ([]byte, error) {
	evaluateResult, err := contract.EvaluateTransaction("RetrieveUserRSAPubkey", username)
	if err != nil {
		return nil, err
	}
	if evaluateResult == nil {
		return nil, fmt.Errorf("Error: pubkey retrieved for user '%v' is nil", username)
	}
	return evaluateResult, nil
}

/*
func ChainCreatePrescriptionSimple(contract *client.Contract, prescription *Prescription) error {
	// Encode Prescription
	encoded, err := encodePrescription(prescription)
	if err != nil {
		return fmt.Errorf("Failed to encode prescription: ", err)
	}
	//Get pubkey of current user

	//step 3: retrieve pubkey of other user
	//step 4: encrypt
	//optional step 5: base64
}
*/
