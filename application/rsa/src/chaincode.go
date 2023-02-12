package src

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc/status"
)

func ChainStoreUserPubkey(contract *client.Contract, username string, pubkey []byte) error {
	b64pubkey := base64.StdEncoding.EncodeToString(pubkey)
	_, err := contract.SubmitTransaction("StoreUserRSAPubkey", username, b64pubkey)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainRetrieveUserPubkey(contract *client.Contract, username string) ([]byte, error) {
	evaluateResult, err := contract.EvaluateTransaction("RetrieveUserRSAPubkey", username)
	if err != nil {
		return nil, ChaincodeParseError(err)
	}
	if evaluateResult == nil {
		return nil, fmt.Errorf("Error: pubkey retrieved for user '%v' is nil", username)
	}
	return base64.StdEncoding.DecodeString(string(evaluateResult))
}

func ChainCreatePrescriptionSimple(contract *client.Contract, prescription *Prescription) error {
	//Get hash(prescription_id + userid_shared_to), used as key for ledger private collection
	tag_raw := sha256.Sum256([]byte(fmt.Sprintf("%v", prescription.Id) + getCurrentUser()))
	tag := base64.StdEncoding.EncodeToString(tag_raw[:])

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

	_, err = contract.SubmitTransaction("CreatePrescription", tag, b64encrypted)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil

}

func ChainReadPrescription(contract *client.Contract, prescriptionId string) (*Prescription, error) {
	// get hash(prescription_id + current_userid)
	tag_raw := sha256.Sum256([]byte(prescriptionId + getCurrentUser()))
	tag := base64.StdEncoding.EncodeToString(tag_raw[:])

	// retrieve from smart contract
	pdata, err := contract.EvaluateTransaction("ReadPrescription", tag)
	if err != nil {
		return nil, ChaincodeParseError(err)
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

func ChaincodeParseError(err error) error {
	var errorString string
	switch err := err.(type) {
	case *client.EndorseError:
		errorString += fmt.Sprintf("Endorse error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
	case *client.SubmitError:
		errorString += fmt.Sprintf("Submit error for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
	case *client.CommitStatusError:
		if errors.Is(err, context.DeadlineExceeded) {
			errorString += fmt.Sprintf("Timeout waiting for transaction %s commit status: %s", err.TransactionID, err)
		} else {
			errorString += fmt.Sprintf("Error obtaining commit status for transaction %s with gRPC status %v: %s\n", err.TransactionID, status.Code(err), err)
		}
	case *client.CommitError:
		errorString += fmt.Sprintf("Transaction %s failed to commit with status %d: %s\n", err.TransactionID, int32(err.Code), err)
	default:
		panic(fmt.Errorf("unexpected error type %T: %w", err, err))
	}

	// Any error that originates from a peer or orderer node external to the gateway will have its details
	// embedded within the gRPC status error. The following code shows how to extract that.
	statusErr := status.Convert(err)

	details := statusErr.Details()
	if len(details) > 0 {
		errorString += "Error Details:\n"

		for _, detail := range details {
			switch detail := detail.(type) {
			case *gateway.ErrorDetail:
				errorString += fmt.Sprintf("- address: %s, mspId: %s, message: %s\n", detail.Address, detail.MspId, detail.Message)
			}
		}
	}
	return fmt.Errorf("\033[0;31m%v\033[0m", errorString)
}

func ChainTestMethod(contract *client.Contract) error {
	return nil
}
