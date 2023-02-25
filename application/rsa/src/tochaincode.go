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

// ====================================================================//
// Send Pubkey
// ====================================================================//
func SendPubkey(contract *client.Contract, username string) {
	obscureName, b64pubkey := PrepareSendPubkey(username)
	SubmitSendPubkey(contract, obscureName, b64pubkey)
}
func PrepareSendPubkey(username string) (string, string) {
	obscureName := obscureName(username)
	pubkey, err := readLocalKey(obscureName, pubFilename)
	if err != nil {
		panic(err)
	}
	return obscureName, base64.StdEncoding.EncodeToString(pubkey)

}
func SubmitSendPubkey(contract *client.Contract, obscureName string, b64pubkey string) {
	_, err := contract.SubmitTransaction("StoreUserRSAPubkey", obscureName, b64pubkey)
	if err != nil {
		panic(ChaincodeParseError(err))
	}
}

// ====================================================================//
// Get Pubkey
// ====================================================================//
func GetPubkey(contract *client.Contract, obscureName string) *rsa.PublicKey {
	return ProcessGetPubkey(EvaluateGetPubkey(contract, obscureName))
}
func EvaluateGetPubkey(contract *client.Contract, obscureName string) string {
	evaluateResult, err := contract.EvaluateTransaction("RetrieveUserRSAPubkey", obscureName)
	if err != nil {
		panic(ChaincodeParseError(err))
	}
	if evaluateResult == nil {
		panic(fmt.Errorf("error: pubkey retrieved for user '%v' is nil", obscureName))
	}
	return string(evaluateResult)
}
func ProcessGetPubkey(evaluateResult string) *rsa.PublicKey {
	decoded, err := base64.StdEncoding.DecodeString(string(evaluateResult))
	if err != nil {
		panic(fmt.Errorf("base64 decoding failed on retrieved pubkey: %v", err))
	}
	pubkey, err := parsePubkey(decoded)
	if err != nil {
		panic(err)
	}
	return pubkey
}

// ====================================================================//
// Create Prescription
// ====================================================================//
func CreatePrescription(contract *client.Contract) string {
	return SubmitCreatePrescription(contract, PrepareCreatePrescription())
}
func PrepareCreatePrescription() string {
	prescription := Prescription{
		Brand:          "NULL",
		Dosage:         "NULL",
		PatientName:    "NULL",
		PatientAddress: "NULL",
		PrescriberName: "NULL",
		PrescriberNo:   0,
		PiecesFilled:   0,
		PiecesTotal:    0,
	}
	pubkey, err := readLocalPubkey(currentUserObscure())
	if err != nil {
		panic(err)
	}
	b64encrypted, err := packagePrescription(pubkey, &prescription)
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

// ====================================================================//
// Read Prescription
// ====================================================================//
func ReadPrescription(contract *client.Contract, pid string) *Prescription {
	return ProcessReadPrescription(EvaluateReadPrescription(contract, pid))
}
func EvaluateReadPrescription(contract *client.Contract, pid string) string {
	// Retrieve from smart contract
	pdata, err := contract.EvaluateTransaction("ReadPrescription", pid)
	if err != nil {
		panic(ChaincodeParseError(err))
	}
	return string(pdata)
}
func ProcessReadPrescription(pdata string) *Prescription {
	// Unpackage and return the prescription
	out, err := unpackagePrescription(string(pdata))
	if err != nil {
		panic(err)
	}
	return out
}

// ====================================================================//
// Share Prescription
// ====================================================================//
func SharePrescription(contract *client.Contract, pid string, username string) {
	obscureName, b64encrypted := PrepareSharePrescription(contract, pid, username)
	SubmitSharePrescription(contract, pid, obscureName, b64encrypted)
}

func PrepareSharePrescription(contract *client.Contract, pid string, username string) (string, string) {
	obscureName := obscureName(username)
	//Retrieve prescription with current user credentials
	prescription := ReadPrescription(contract, pid)
	//Request pubkey from username to share to
	otherPubkey := GetPubkey(contract, obscureName)

	//Re-encrypt the prescription with the new user credentials
	b64encrypted, err := packagePrescription(otherPubkey, prescription)
	if err != nil {
		panic(err)
	}
	return obscureName, b64encrypted
}
func SubmitSharePrescription(contract *client.Contract, pid string, obscureName string, b64encrypted string) {
	//Save prescription with tag
	_, err := contract.SubmitTransaction("SharePrescription", pid, obscureName, b64encrypted)
	if err != nil {
		panic(ChaincodeParseError(err))
	}
}

func SharedToList(contract *client.Contract, pid string) *[]string {
	// Get list of all users that the prescription was shared to
	b64strings, err := contract.EvaluateTransaction("PrescriptionSharedTo", pid)
	if err != nil {
		panic(ChaincodeParseError(err))
	}
	// base64 decode
	sharedto, err := unpackageStringSlice(string(b64strings))
	if err != nil {
		panic(err)
	}
	return sharedto
}

// ====================================================================//
// Re-encrypt Prescription Set
// ====================================================================//
func reencryptPrescriptionSet(contract *client.Contract, pid string, update *Prescription) (string, error) {
	usernames := SharedToList(contract, pid)
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

// ====================================================================//
// Update Prescription
// ====================================================================//
func UpdatePrescription(contract *client.Contract, pid string, update *Prescription) {
	SubmitUpdatePrescription(contract, pid, PrepareUpdatePrescription(contract, pid, update))
}
func PrepareUpdatePrescription(contract *client.Contract, pid string, update *Prescription) string {
	b64gob, err := reencryptPrescriptionSet(contract, pid, update)
	if err != nil {
		panic(err)
	}
	return b64gob
}
func SubmitUpdatePrescription(contract *client.Contract, pid string, b64gob string) {
	_, err := contract.SubmitTransaction("UpdatePrescription", pid, b64gob)
	if err != nil {
		panic(ChaincodeParseError(err))
	}
}

// ====================================================================//
// Setfill Prescription
// ====================================================================//
func SetfillPrescription(contract *client.Contract, pid string, newfill uint8) {
	SubmitSetfillPrescription(contract, pid, PrepareSetfillPrescription(contract, pid, newfill))
}
func PrepareSetfillPrescription(contract *client.Contract, pid string, newfill uint8) string {
	prescription := ReadPrescription(contract, pid)

	prescription.PiecesFilled = newfill

	b64gob, err := reencryptPrescriptionSet(contract, pid, prescription)
	if err != nil {
		panic(err)
	}
	return b64gob
}
func SubmitSetfillPrescription(contract *client.Contract, pid string, b64gob string) {
	_, err := contract.SubmitTransaction("SetfillPrescription", pid, b64gob)
	if err != nil {
		panic(ChaincodeParseError(err))
	}
}

// ====================================================================//
// Delete Prescription
// ====================================================================//
func DeletePrescription(contract *client.Contract, pid string) error {
	_, err := contract.SubmitTransaction("DeletePrescription", pid)
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

// ====================================================================//
// Report Register
// ====================================================================//
func ChainReportAddReader(contract *client.Contract) error {
	_, err := contract.SubmitTransaction("RegisterMeAsReportReader")
	if err != nil {
		return ChaincodeParseError(err)
	}
	return nil
}

// ====================================================================//
// Report Get Readers
// ====================================================================//
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

// ====================================================================//
// Report Update
// ====================================================================//
func ReportUpdate(contract *client.Contract, pid string) {
	SubmitReportUpdate(contract, pid, PrepareReportUpdate(contract, pid))
}

func PrepareReportUpdate(contract *client.Contract, pid string) string {
	readers, err := ChainReportGetReaders(contract)
	if err != nil {
		panic(err)
	}
	prescription := ReadPrescription(contract, pid)

	pset := make(map[string]string)
	for _, obscuredName := range *readers {
		pubkey := GetPubkey(contract, obscuredName)
		b64encrypted, err := packagePrescription(pubkey, prescription)
		if err != nil {
			panic(err)
		}
		pset[obscuredName] = b64encrypted
	}
	b64reports, err := packagePrescriptionSet(&pset)
	if err != nil {
		panic(err)
	}
	return b64reports
}

func SubmitReportUpdate(contract *client.Contract, pid string, b64reports string) {
	_, err := contract.SubmitTransaction("UpdateReport", pid, b64reports)
	if err != nil {
		panic(ChaincodeParseError(err))
	}
}

// ====================================================================//
// Report View
// ====================================================================//
func ReportView(contract *client.Contract) string {
	return ProcessReportView(EvaluateReportView(contract))
}
func EvaluateReportView(contract *client.Contract) string {
	b64all, err := contract.EvaluateTransaction("GetPrescriptionReport")
	if err != nil {
		panic(err)
	}
	return string(b64all)
}
func ProcessReportView(b64all string) string {
	prescriptions, err := unpackagePrescriptionSet(string(b64all))
	if err != nil {
		panic(err)
	}
	var output string
	for _, pdata := range *prescriptions {
		if pdata != "" {
			prescription, err := unpackagePrescription(pdata)
			if err != nil {
				panic(err)
			}
			output += fmt.Sprintf("prescription: %v\n", prescription)
		}
	}

	return output
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
