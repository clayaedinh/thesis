package main

import (
	"flag"
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
	BLUE   = "\033[0;34m"
	PURPLE = "\033[0;35m"
	GRAY   = "\033[1;30m"
	NC     = "\033[0m"
)

const (
	FLAG_H_ORG  = "Specifies the org that the current user belongs to."
	FLAG_H_USER = "Specifies the user that connects to the network."
	FLAG_H_PORT = "Specifies the port which the organization peer belongs to."
)

func printHelp() {
	fmt.Println("")
	fmt.Printf("%vPrescription Blockchain Thesis Application, RSA version%v\n", YELLOW, NC)
	fmt.Println("This application enables users to call chaincode remotely.")
	fmt.Println("")
	fmt.Printf("%vUsage%v: ./rsa %v[-options] %v<Method> %v<Method Args>\n", GREEN, NC, PURPLE, CYAN, NC)
	fmt.Println("")
	fmt.Printf("%vAvailable Options%v:\n", GREEN, NC)
	fmt.Println("")
	fmt.Printf("./rsa %v-user=%vstring\n", PURPLE, NC)
	fmt.Println(FLAG_H_USER)
	fmt.Println("")
	fmt.Printf("./rsa %v-org=%vstring\n", PURPLE, NC)
	fmt.Println(FLAG_H_ORG)
	fmt.Println("")
	fmt.Printf("./rsa %v-port=%vlocalhost:port\n", PURPLE, NC)
	fmt.Println(FLAG_H_PORT)
	fmt.Println("")
	fmt.Printf("%vAvailable Methods (must be AFTER options)%v:\n", GREEN, NC)
	fmt.Printf("./rsa %vstorekey%v <username>\n", CYAN, NC)
	fmt.Println("Saves the RSA public key of the given username on the ledger.")
	fmt.Println("")
	fmt.Printf("./rsa %vgetkey%v <username>\n", CYAN, NC)
	fmt.Println("Retrieves the RSA public key of the given username from the ledger.")
	fmt.Println("")
	fmt.Printf("./rsa %vcreatep%v <json>\n", CYAN, NC)
	fmt.Println("Creates a prescription with json data. JSON follows this format:")
	fmt.Println(src.PrintEmptyPrescriptionJSON())
	fmt.Println("")

}
func main() {
	//Help Menu
	if len(os.Args) == 1 || strings.ToLower(os.Args[1]) == "help" {
		printHelp()
		os.Exit(0)
	}

	//Flags
	flagOrg := flag.String("org", "org1", FLAG_H_ORG)
	flagUser := flag.String("user", "Admin", FLAG_H_USER)
	flagPort := flag.String("port", "localhost:7051", FLAG_H_PORT)

	flag.Parse()

	//If application is not printing help, it will be interacting with chaincode
	//So start connection
	src.SetConnectionVariables(*flagOrg, *flagUser, *flagPort)
	clientConnection, err := src.NewGrpcConnection()
	if err != nil {
		log.Panic(err)
	}
	defer clientConnection.Close()

	gw, err := src.DefaultGateway(clientConnection)
	if err != nil {
		log.Panic(err)
	}
	defer gw.Close()
	contract := src.SmartContract(gw)

	//We now check which chaincode function is being called
	if flag.Arg(0) == "storekey" {
		storekey(contract, flag.Arg(1))
	} else if flag.Arg(0) == "getkey" {
		getkey(contract, flag.Arg(1))
	} else if flag.Arg(0) == "testprint" {
		src.TestPrintPrescriptionJSON()
	} else {
		fmt.Printf("%vInvalid method '%v'. Do './rsa help' for method options.\n", RED, flag.Arg(0))
	}
}

func storekey(contract *client.Contract, username string) {
	pubkey, err := src.ReadUserPubkey(username)
	if err != nil {
		log.Panic(err)
	}
	err = src.ChainStoreUserPubkey(contract, username, pubkey)
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

func createp(contract *client.Contract, jsonstring string) {
	prescription, err := src.PrescriptionFromCmdInput(jsonstring)
	if err != nil {
		log.Panic(err)
	}
	err = src.ChainCreatePrescriptionSimple(contract, prescription)
	if err != nil {
		log.Panic(err)
	}
}
