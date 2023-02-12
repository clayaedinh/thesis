package src

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
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
	pubkey, err := localPrivkeyBytes(username)
	if err != nil {
		return err
	}

	b64pubkey := base64.StdEncoding.EncodeToString(pubkey)
	_, err = contract.SubmitTransaction("StoreUserRSAPubkey", username, b64pubkey)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainRetrievePubkey(contract *client.Contract, username string) (*rsa.PublicKey, error) {
	evaluateResult, err := contract.EvaluateTransaction("RetrieveUserRSAPubkey", username)
	if err != nil {
		return nil, ChaincodeParseError(err)
	}
	if evaluateResult == nil {
		return nil, fmt.Errorf("Error: pubkey retrieved for user '%v' is nil", username)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(evaluateResult))
	if err != nil {
		return nil, fmt.Errorf("Base64 decoding failed on retrieved pubkey: %v", err)
	}
	return parsePubkeyBytes(decoded)
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

	// Generate a tag for the prescription
	tag := genPrescriptionTag(fmt.Sprintf("%v", prescription.Id), getCurrentUser())

	// Get current user pubkey
	pubkey, err := localPubkey(getCurrentUser())
	if err != nil {
		return err
	}
	// package the prescription
	b64encrypted, err := packagePrescription(pubkey, &prescription)
	if err != nil {
		return err
	}

	_, err = contract.SubmitTransaction("SavePrescription", tag, b64encrypted)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil

}

func ChainUpdatePrescription(contract *client.Contract, update *Prescription) error {
	tag := genPrescriptionTag(fmt.Sprintf("%v", update.Id), getCurrentUser())
	pubkey, err := localPubkey(getCurrentUser())
	if err != nil {
		return err
	}
	b64encrypted, err := packagePrescription(pubkey, update)
	if err != nil {
		return err
	}

	_, err = contract.SubmitTransaction("SavePrescription", tag, b64encrypted)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainReadPrescription(contract *client.Contract, pid string) (*Prescription, error) {
	// Get Tag of prescription to read
	tag := genPrescriptionTag(pid, getCurrentUser())

	// Retrieve from smart contract
	pdata, err := contract.EvaluateTransaction("ReadPrescription", tag)
	if err != nil {
		return nil, ChaincodeParseError(err)
	}
	// Unpackage and return the prescription
	return unpackagePrescription(string(pdata))
}

func ChainSharePrescription(contract *client.Contract, pid string, username string) error {
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
	tag := genPrescriptionTag(pid, username)

	//Save prescription with tag
	_, err = contract.SubmitTransaction("SavePrescription", tag, b64encrypted)
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

	tag := genPrescriptionTag(pid, getCurrentUser())

	pubkey, err := localPubkey(getCurrentUser())
	if err != nil {
		return err
	}

	b64encrypted, err := packagePrescription(pubkey, prescription)
	if err != nil {
		return err
	}

	_, err = contract.SubmitTransaction("SavePrescription", tag, b64encrypted)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
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
	prescription, err := ChainReadPrescription(contract, pid)
	if err != nil {
		return fmt.Errorf("Failed to read prescription %v : %v", pid, err)
	}
	readers, err := ChainReportGetReaders(contract)
	if err != nil {
		return err
	}

	for _, reader := range readers {
		username := string(reader)
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
		tag := genPrescriptionTag(pid, username)

		//Save prescription with tag
		_, err = contract.SubmitTransaction("SavePrescription", tag, b64encrypted)
		if err != nil {
			return ChaincodeParseError(err)
		}
	}
	return nil
}

/*
func ChainReportView(contract *client.Contract) error {
	pdatas, err := contract.EvaluateTransaction("GetAllPrescriptions")
	if err != nil {
		return err
	}

	for _, pdata := range pdatas {
		prescription, err := unpackagePrescription(pdata)
		if err != nil {
			return err
		}
	}
}
*/

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
