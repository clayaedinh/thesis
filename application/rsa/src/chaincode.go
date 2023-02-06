package src

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func CreatePrescription(contract *client.Contract) {
	//step 1: struct
	//step 2: encode
	//step 3: retrieve pubkey of other user
	//step 4: encrypt
	//optional step 5: base64
}

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
