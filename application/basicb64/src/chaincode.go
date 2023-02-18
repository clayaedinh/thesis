package src

import (
	"context"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc/status"
)

func ChainCreatePrescription(contract *client.Contract) error {
	// Create a new, blank prescription
	var prescription Prescription

	// Assign an id to the prescription
	prescription.Id = genPrescriptionId()

	obscureName := currentUserObscure()

	// package the prescription
	b64prescription, err := packagePrescription(&prescription)
	if err != nil {
		return err
	}
	_, err = contract.SubmitTransaction("CreatePrescription", fmt.Sprintf("%v", prescription.Id), obscureName, b64prescription)
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

	//Save prescription with tag
	_, err := contract.SubmitTransaction("SharePrescription", pid, obscureName)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainUpdatePrescription(contract *client.Contract, update *Prescription) error {
	b64prescription, err := packagePrescription(update)
	if err != nil {
		return err
	}
	_, err = contract.SubmitTransaction("UpdatePrescription", fmt.Sprintf("%v", update.Id), b64prescription)
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

	b64prescription, err := packagePrescription(prescription)
	if err != nil {
		return err
	}
	_, err = contract.SubmitTransaction("SetfillPrescription", pid, b64prescription)
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
