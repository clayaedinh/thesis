package chaincode

import "github.com/hyperledger/fabric-contract-api-go/contractapi"

const (
	collectionPrescription = "collectionPrescription"
	collectionPubkeyRSA    = "collectionPubkeyRSA"
)

type SmartContract struct {
	contractapi.Contract
}
