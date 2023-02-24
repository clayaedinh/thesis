package src

import (
	"context"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc/status"
)

/*
Each function is awkwardly split into two, so that each half may be benchmarked
*/

func CreatePrescription(contract *client.Contract) string {
	return SubmitCreatePrescription(contract, PrepareCreatePrescription())
}
func PrepareCreatePrescription() string {
	var prescription Prescription
	b64encrypted, err := packagePrescription(&prescription)
	if err != nil {
		panic(err)
	}
	return b64encrypted
}
func SubmitCreatePrescription(contract *client.Contract, b64encrypted string) string {
	pid, err := contract.SubmitTransaction("CreatePrescription", b64encrypted)
	if err != nil {
		panic(ChaincodeParseError(err))
	}
	return string(pid)
}

func ReadPrescription(contract *client.Contract, pid string) *Prescription {
	return ProcessReadPrescription(EvaluateReadPrescription(contract, pid))
}
func EvaluateReadPrescription(contract *client.Contract, pid string) string {
	pdata, err := contract.EvaluateTransaction("ReadPrescription", pid, currentUserObscure())
	if err != nil {
		panic(ChaincodeParseError(err))
	}
	return string(pdata)
}
func ProcessReadPrescription(pdata string) *Prescription {
	prescription, err := unpackagePrescription(pdata)
	if err != nil {
		panic(err)
	}
	return prescription
}

func ChainSharePrescription(contract *client.Contract, pid string, username string) error {
	obscureName := obscureName(username)
	// Invoke Share Prescription on chaincode
	_, err := contract.SubmitTransaction("SharePrescription", pid, obscureName)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainUpdatePrescription(contract *client.Contract, pid string, update *Prescription) error {
	b64prescription, err := packagePrescription(update)
	if err != nil {
		return err
	}
	_, err = contract.SubmitTransaction("UpdatePrescription", pid, b64prescription)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

func ChainSetfillPrescription(contract *client.Contract, pid string, newfill uint8) error {
	prescription := ReadPrescription(contract, pid)

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

func ChainReportView(contract *client.Contract) (*[]string, error) {
	b64report, err := contract.SubmitTransaction("GetPrescriptionReport")
	if err != nil {
		return nil, ChaincodeParseError(err)
	}
	report, err := unpackageStringSlice(string(b64report))
	if err != nil {
		return nil, err
	}
	return report, nil
}

// From Fabric Samples: Basic Asset Transfer (Gateway), Golang
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
