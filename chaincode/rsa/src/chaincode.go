package src

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

const (
	collectionPrescription  = "collectionPrescription"
	collectionPubkeyRSA     = "collectionPubkeyRSA"
	collectionReportReaders = "collectionReportReaders"
	collectionReports       = "collectionReports"
)

type SmartContract struct {
	contractapi.Contract
}
