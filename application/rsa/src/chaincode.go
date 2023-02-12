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
		return nil, fmt.Errorf("Error: failed to retrieve pubkey of user %v: %v", username, err)
	}
	if evaluateResult == nil {
		return nil, fmt.Errorf("Error: pubkey retrieved for user '%v' is nil", username)
	}
	return base64.StdEncoding.DecodeString(string(evaluateResult))
}

func ChainCreatePrescriptionSimple(contract *client.Contract, prescription *Prescription) error {
	fmt.Printf("prescription: %v\n", prescription)

	//Get hash(prescription_id + userid_shared_to), used as key for ledger private collection
	tag_raw := sha256.Sum256([]byte(fmt.Sprintf("%v", prescription.Id) + getCurrentUser()))
	tag := string(tag_raw[:])

	// Encode Prescription to Bytes
	encoded, err := encodePrescription(prescription)
	if err != nil {
		return fmt.Errorf("Failed to encode prescription: %v", err)
	}

	//Obtain Public Key of current user
	rawkey, err := ChainRetrieveUserPubkey(contract, getCurrentUser())
	if err != nil {
		return fmt.Errorf("Failed to retrieve public key of user %v: %v", getCurrentUser(), err)
	}
	//pubkey, err := keyFromChainRetrieval(rawbytes)
	pubkey, err := parsePubkeyBytes(rawkey)
	if err != nil {
		return fmt.Errorf("Failed to parse public key: %v", err)
	}

	//Encrypt data with current user's public key
	encrypted, err := encryptBytes(encoded, pubkey)
	if err != nil {
		return fmt.Errorf("Failed to encrypt prescription: %v", err)
	}

	//Encode data as base64
	b64encrypted := base64.StdEncoding.EncodeToString(encrypted)
	_, err = contract.SubmitTransaction("CreatePrescriptionSimple", tag, b64encrypted)

	_, err = contract.SubmitTransaction("CreatePrescriptionSimple", tag, b64encrypted)

	if err != nil {
		return fmt.Errorf("CreatePrescriptionSimple smart contract failed: %v", err)
	}
	return nil

}

func ChainReadPrescription(contract *client.Contract, prescriptionId string) (*Prescription, error) {
	// get hash(prescription_id + current_userid)
	tag_raw := sha256.Sum256([]byte(prescriptionId + getCurrentUser()))
	tag := string(tag_raw[:])

	// retrieve from smart contract
	pdata, err := contract.EvaluateTransaction("ReadPrescription", tag)
	if err != nil {
		return nil, fmt.Errorf("ReadPrescription smart contract failed: %v", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(pdata))
	if err != nil {
		return nil, fmt.Errorf("Base64 failed to decrypt prescription: %v", err)
	}

	//Obtain Private Key of current user
	rawkey, err := ReadUserPrivkey(getCurrentUser())
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve private key of user %v: %v", getCurrentUser(), err)
	}
	//pubkey, err := keyFromChainRetrieval(rawkey)
	privkey, err := parsePrivkeyBytes(rawkey)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key of user %v: %v", getCurrentUser(), err)
	}

	// Decrypt data with current user's public key
	decrypted, err := decryptBytes(decoded, privkey)
	if err != nil {
		return nil, fmt.Errorf("Failed to encrypt prescription: %v", err)
	}

	// Decode Prescription to Bytes
	prescription, err := decodePrescription(decrypted)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode prescription: %v", err)
	}
	return prescription, nil
}
