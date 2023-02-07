package src

import (
	"crypto/sha256"
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

func ChainCreatePrescriptionSimple(contract *client.Contract, prescription *Prescription) error {

	//Get hash(prescription_id + userid_shared_to), used as key for ledger private collection
	tag_raw := sha256.Sum256([]byte(prescription.Id + getCurrentUser()))
	tag := string(tag_raw[:])

	// Encode Prescription to Bytes
	encoded, err := encodePrescription(prescription)
	if err != nil {
		return fmt.Errorf("Failed to encode prescription: ", err)
	}

	//Obtain Public Key of current user
	rawbytes, err := ChainRetrieveUserPubkey(contract, getCurrentUser())
	if err != nil {
		return fmt.Errorf("Failed to retrieve public key of user %v: %v", getCurrentUser(), err)
	}
	pubkey, err := keyFromChainRetrieval(rawbytes)
	if err != nil {
		return fmt.Errorf("Failed to parse public key: %v", err)
	}

	//Encrypt data with current user's public key
	encrypted, err := encryptBytes(encoded, pubkey)
	if err != nil {
		return fmt.Errorf("Failed to encrypt prescription: %v", err)
	}
	// Encode data as base64
	b64encrypted := base64.StdEncoding.EncodeToString(encrypted)

	_, err = contract.SubmitTransaction("CreatePrescriptionSimple", tag, b64encrypted)
	return nil

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
