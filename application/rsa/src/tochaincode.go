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

func ChainSendPubkey(contract *client.Contract, username string) error {
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

func ChainGetPubkey(contract *client.Contract, username string) (*rsa.PublicKey, error) {
	obscureName := obscureName(username)
	return chainGetPubkeyObscuredName(contract, obscureName)
}

func chainGetPubkeyObscuredName(contract *client.Contract, obscureName string) (*rsa.PublicKey, error) {
	evaluateResult, err := contract.EvaluateTransaction("RetrieveUserRSAPubkey", obscureName)
	if err != nil {
		return nil, ChaincodeParseError(err)
	}
	if evaluateResult == nil {
		return nil, fmt.Errorf("error: pubkey retrieved for user '%v' is nil", obscureName)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(evaluateResult))
	if err != nil {
		return nil, fmt.Errorf("base64 decoding failed on retrieved pubkey: %v", err)
	}
	return parsePubkey(decoded)
}

func ChainCreatePrescription(contract *client.Contract) (string, error) {
	var prescription Prescription
	pubkey, err := readLocalPubkey(currentUserObscure())
	if err != nil {
		return "", err
	}
	b64encrypted, err := packagePrescription(pubkey, &prescription)
	if err != nil {
		return "", err
	}
	pid, err := contract.SubmitTransaction("CreatePrescription", b64encrypted)
	if err != nil {
		return "", ChaincodeParseError(err)
	}
	return string(pid), nil
}

func ChainReadPrescription(contract *client.Contract, pid string) (*Prescription, error) {
	// Retrieve from smart contract
	pdata, err := contract.EvaluateTransaction("ReadPrescription", pid)
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
	otherPubkey, err := ChainGetPubkey(contract, username)
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
func ChainSharedToList(contract *client.Contract, pid string) (*[]string, error) {
	// Get list of all users that the prescription was shared to
	b64strings, err := contract.EvaluateTransaction("PrescriptionSharedTo", pid)
	if err != nil {
		return nil, ChaincodeParseError(err)
	}
	// base64 decode
	sharedto, err := unpackageStringSlice(string(b64strings))
	if err != nil {
		return nil, err
	}
	return sharedto, nil
}
func reencryptPrescriptionSet(contract *client.Contract, pid string, update *Prescription) (string, error) {
	usernames, err := ChainSharedToList(contract, pid)
	if err != nil {
		return "", err
	}
	pset := make(map[string]string)
	//Encrypt for each username
	for _, username := range *usernames {
		pubkey, err := readLocalPubkey(username)
		if err != nil {
			return "", err
		}
		b64encrypted, err := packagePrescription(pubkey, update)
		if err != nil {
			return "", err
		}
		pset[username] = b64encrypted
	}
	b64gob, err := packagePrescriptionSet(&pset)
	if err != nil {
		return "", err
	}
	return b64gob, nil
}

func ChainUpdatePrescription(contract *client.Contract, pid string, update *Prescription) error {
	b64gob, err := reencryptPrescriptionSet(contract, pid, update)
	if err != nil {
		return err
	}
	_, err = contract.SubmitTransaction("UpdatePrescription", pid, b64gob)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainSetfillPrescription(contract *client.Contract, pid string, newfill uint8) error {
	prescription, err := ChainReadPrescription(contract, pid)
	if err != nil {
		return fmt.Errorf("failed to read prescription %v : %v", pid, err)
	}
	prescription.PiecesFilled = newfill

	b64gob, err := reencryptPrescriptionSet(contract, pid, prescription)
	if err != nil {
		return err
	}
	_, err = contract.SubmitTransaction("SetfillPrescription", pid, b64gob)
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

func ChainReportAddReader(contract *client.Contract) error {
	_, err := contract.SubmitTransaction("RegisterMeAsReportReader")
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainReportGetReaders(contract *client.Contract) (*[]string, error) {
	b64readers, err := contract.EvaluateTransaction("GetAllReportReaders")
	if err != nil {
		return nil, err
	}
	strings, err := unpackageStringSlice(string(b64readers))
	if err != nil {
		return nil, err
	}
	return strings, nil
}

/*
Encrypts the given prescription at the given pid for all report readers
Probably has to be called every time prescription is updated
this is why role-based encryption is better
*/
func ChainReportUpdate(contract *client.Contract, pid string) error {
	readers, err := ChainReportGetReaders(contract)
	if err != nil {
		return err
	}
	prescription, err := ChainReadPrescription(contract, pid)
	if err != nil {
		return err
	}
	pset := make(map[string]string)
	for _, obscuredName := range *readers {
		pubkey, err := chainGetPubkeyObscuredName(contract, obscuredName)
		if err != nil {
			return err
		}
		b64encrypted, err := packagePrescription(pubkey, prescription)
		if err != nil {
			return err
		}
		pset[obscuredName] = b64encrypted
	}
	b64reports, err := packagePrescriptionSet(&pset)
	if err != nil {
		return err
	}
	_, err = contract.SubmitTransaction("UpdateReport", pid, b64reports)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainReportView(contract *client.Contract) (string, error) {
	b64all, err := contract.EvaluateTransaction("GetPrescriptionReport")
	if err != nil {
		return "", err
	}
	prescriptions, err := unpackagePrescriptionSet(string(b64all))
	if err != nil {
		return "", err
	}
	var output string
	for _, pdata := range *prescriptions {
		prescription, err := unpackagePrescription(pdata)
		if err != nil {
			return "", err
		}
		output += fmt.Sprintf("prescription: %v\n", prescription)
	}

	return output, nil
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
