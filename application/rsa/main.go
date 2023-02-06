package main

import (
	"fmt"
	"log"

	"github.com/clayaedinh/thesis/application/rsa/src"
)

func main() {
	//get public key

	pubkey, err := src.ReadUserPubkey("user0001")
	if err != nil {
		log.Panic(err)
	}
	//connection

	src.SetConnectionVariables("org1", "Admin", "localhost:7051")
	clientConnection := src.NewGrpcConnection()
	defer clientConnection.Close()

	gw, err := src.DefaultGateway(clientConnection)
	if err != nil {
		log.Panic(err)
	}
	defer gw.Close()

	contract := src.SmartContract(gw)

	err = src.ChainStoreUserPubkey(contract, "user0002", pubkey)
	if err != nil {
		log.Panic(err)
	}
	out, err := src.ChainRetrieveUserPubkey(contract, "user0002")
	if err != nil {
		log.Panic(err)
	}
	fmt.Print(out)

}
