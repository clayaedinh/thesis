package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/clayaedinh/thesis/application/rsa/src"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

const (
	RED    = "\033[0;31m"
	YELLOW = "\033[1;33m"
	GREEN  = "\033[0;32m"
	CYAN   = "\033[0;36m"
	PURPLE = "\033[0;35m"
	NC     = "\033[0m"
)

func printHelp() {
	fmt.Println("")
	fmt.Printf("%vPrescription Blockchain Thesis Application, RSA version%v\n", YELLOW, NC)
	fmt.Println("This application enables users to call chaincode remotely.")
	fmt.Println("")
	fmt.Printf("%vUsage%v: ./rsa %v<Method> %v<Method Args>%v\n", GREEN, NC, CYAN, PURPLE, NC)
	fmt.Println("")
	fmt.Printf("%vAvailable Methods%v:\n", GREEN, NC)
	fmt.Printf("./rsa %vstorekey%v %v<username>%v\n", CYAN, NC, PURPLE, NC)
	fmt.Println("Saves the RSA public key of the given username on the ledger.")
	fmt.Println("")
	fmt.Printf("./rsa %vgetkey%v %v<username>%v\n", CYAN, NC, PURPLE, NC)
	fmt.Println("Retrieves the RSA public key of the given username from the ledger.")
	fmt.Println("")

}
func main() {
	if len(os.Args) == 1 || strings.ToLower(os.Args[1]) == "help" {
		printHelp()
		os.Exit(0)
	}

	//If application is not printing help, it will be interacting with chaincode
	//So start connection
	src.SetConnectionVariables("org1", "Admin", "localhost:7051")
	clientConnection := src.NewGrpcConnection()
	defer clientConnection.Close()
	gw, err := src.DefaultGateway(clientConnection)
	if err != nil {
		log.Panic(err)
	}
	defer gw.Close()
	contract := src.SmartContract(gw)

	//We now check which chaincode function is being called

	if os.Args[1] == "storekey" {
		storekey(contract, os.Args[2])
	}
	if os.Args[1] == "getkey" {
		getkey(contract, os.Args[2])
	}
}

func storekey(contract *client.Contract, username string) {
	pubkey, err := src.ReadUserPubkey(username)
	if err != nil {
		log.Panic(err)
	}
	err = src.ChainStoreUserPubkey(contract, "user0002", pubkey)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%vKey stored successfully for user %v%v\n", GREEN, username, NC)
}

func getkey(contract *client.Contract, username string) {
	out, err := src.ChainRetrieveUserPubkey(contract, username)
	if err != nil {
		log.Panic(err)
	}
	fmt.Print(out)
	fmt.Printf("\n%vKey retrieved successfully for user %v%v\n", GREEN, username, NC)
}
