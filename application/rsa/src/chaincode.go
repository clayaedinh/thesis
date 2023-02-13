package src

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc/status"
)

// ========================================
// RSA Key Management
// ========================================

func ChainStoreLocalPubkey(contract *client.Contract, username string) error {
	obscureName := obscureName(username)
	pubkey, err := readLocalKey(obscureName, pubFilename)
	if err != nil {
		return err
	}

	b64pubkey := base64.StdEncoding.EncodeToString(pubkey)

	_, err = contract.SubmitTransaction("StoreUserRSAPubkey", obscureName, b64pubkey)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainRetrievePubkey(contract *client.Contract, username string) (*rsa.PublicKey, error) {
	obscureName := obscureName(username)
	evaluateResult, err := contract.EvaluateTransaction("RetrieveUserRSAPubkey", obscureName)
	if err != nil {
		return nil, ChaincodeParseError(err)
	}
	if evaluateResult == nil {
		return nil, fmt.Errorf("Error: pubkey retrieved for user '%v' is nil", obscureName)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(evaluateResult))
	if err != nil {
		return nil, fmt.Errorf("Base64 decoding failed on retrieved pubkey: %v", err)
	}
	return parsePubkey(decoded)
}

// =============================================
// Create Prescription
// Creates a completely BLANK prescription
// =============================================

func ChainCreatePrescription(contract *client.Contract) error {
	// Create a new, blank prescription
	var prescription Prescription

	// Assign an id to the prescription
	prescription.Id = genPrescriptionId()

	obscureName := currentUserObscure()
	// Get current user pubkey
	pubkey, err := readLocalPubkey(obscureName)
	if err != nil {
		return err
	}
	// package the prescription
	b64encrypted, err := packagePrescription(pubkey, &prescription)
	if err != nil {
		return err
	}

	_, err = contract.SubmitTransaction("CreatePrescription", fmt.Sprintf("%v", prescription.Id), obscureName, b64encrypted)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil

}

func ChainSharedToList(contract *client.Contract, pid string) ([]string, error) {
	// Get list of all users that the prescription was shared to
	pgob, err := contract.EvaluateTransaction("PrescriptionSharedTo", pid)
	if err != nil {
		return nil, ChaincodeParseError(err)
	}
	// base64 decode
	decoded, err := base64.StdEncoding.DecodeString(string(pgob))
	if err != nil {
		return nil, err
	}
	sharedto, err := decodeStringSlice(decoded)
	return sharedto, nil
}

func ChainUpdatePrescription(contract *client.Contract, update *Prescription) error {
	pid := fmt.Sprintf("%v", update.Id)
	usernames, err := ChainSharedToList(contract, pid)
	if err != nil {
		return err
	}
	pset := make(map[string]string)

	//Encrypt for each username
	for _, username := range usernames {
		pubkey, err := readLocalPubkey(currentUserObscure())
		if err != nil {
			return err
		}
		b64encrypted, err := packagePrescription(pubkey, update)
		if err != nil {
			return err
		}
		pset[username] = b64encrypted
	}

	pgob, err := gobEncode(&pset)
	if err != nil {
		return err
	}

	b64gob := base64.StdEncoding.EncodeToString(pgob)
	_, err = contract.SubmitTransaction("UpdatePrescription", pid, b64gob)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainReadPrescription(contract *client.Contract, pid string) (*Prescription, error) {
	// Retrieve from smart contract
	pdata, err := contract.EvaluateTransaction("ReadPrescription", pid, currentUserObscure())
	if err != nil {
		return nil, ChaincodeParseError(err)
	}
	// Unpackage and return the prescription
	return unpackagePrescription(string(pdata))
}

func ChainSharePrescription(contract *client.Contract, pid string, username string) error {
	obscureName := obscureName(username)
	//Retrieve prescription with current user credentials
	prescription, err := ChainReadPrescription(contract, pid)
	if err != nil {
		return err
	}
	//Request pubkey from username to share to
	otherPubkey, err := ChainRetrievePubkey(contract, username)
	if err != nil {
		return err
	}
	//Re-encrypt the prescription with the new user credentials
	b64encrypted, err := packagePrescription(otherPubkey, prescription)
	if err != nil {
		return err
	}

	//Save prescription with tag
	_, err = contract.SubmitTransaction("SharePrescription", pid, obscureName, b64encrypted)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainDeletePrescription(contract *client.Contract, pid string) error {
	_, err := contract.SubmitTransaction("DeletePrescription", pid)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainSetfillPrescription(contract *client.Contract, pid string, newfill uint8) error {
	prescription, err := ChainReadPrescription(contract, pid)
	if err != nil {
		return fmt.Errorf("Failed to read prescription %v : %v", pid, err)
	}
	prescription.PiecesFilled = newfill
	return ChainUpdatePrescription(contract, prescription)
}

func ChainReportGetReaders(contract *client.Contract) ([]byte, error) {
	readers, err := contract.EvaluateTransaction("GetAllReportReaders")
	if err != nil {
		return nil, err
	}
	return readers, nil
}

// Encrypts the given prescription at the given pid for all report readers
// Probably has to be called every time prescription is updated
// ^ this is why role-based encryption is better

func ChainReportEncrypt(contract *client.Contract, pid string) error {
	readers, err := ChainReportGetReaders(contract)
	if err != nil {
		return err
	}
	for _, reader := range readers {
		username := string(reader)
		err = ChainSharePrescription(contract, pid, username)
		if err != nil {
			return err
		}
	}
	return nil
}

// If this works I will be amazed
func ChainReportView(contract *client.Contract) error {
	pdatas, err := contract.EvaluateTransaction("GetAllPrescriptions")
	if err != nil {
		return err
	}

	//I have no idea if this will work
	prescriptions := [][]byte{}
	enc := gob.NewDecoder(bytes.NewReader(pdatas))
	err = enc.Decode(&prescriptions)
	if err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("Failed to package prescriptions for transport: %v", err)
	}

	for _, pdata := range prescriptions {
		prescription, err := unpackagePrescription(string(pdata))
		if err != nil {
			return err
		}
		fmt.Printf("prescription: %v\n", prescription)
	}
	return nil
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
