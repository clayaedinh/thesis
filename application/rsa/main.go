package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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
	fmt.Printf("./rsa %vgetkey%v <username>\n", CYAN, NC)
	fmt.Printf("./rsa %vcreatep%v\n", CYAN, NC)
	fmt.Printf("./rsa %vsharep%v <pid> <username>\n", CYAN, NC)
	fmt.Printf("./rsa %vreadp%v <id>\n", CYAN, NC)
	fmt.Printf("./rsa %vupdatep%v <brand> <dosage> <patient_name> <patient_address> <doctor_name> <doctor_prc> <pieces_total>\n", CYAN, NC)
	fmt.Printf("./rsa %vsetfillp%v <pid> <newfill>\n", CYAN, NC)
	fmt.Printf("./rsa %vdeletep%v <pid>\n", CYAN, NC)
	fmt.Printf("./rsa %vreaderadd%v <username>\n", CYAN, NC)
	fmt.Printf("./rsa %vreaderall%v\n", CYAN, NC)
	fmt.Printf("./rsa %vreportgen%v <pid>\n", CYAN, NC)
	fmt.Printf("./rsa %vreportread%v <username>\n", CYAN, NC)
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

	// Methods which do not require a connection to the chaincode

	if flag.Arg(0) == "genkey" {
		checkEnoughArgs(2)
		genkey(flag.Arg(1))
		os.Exit(0)
	}

	//If application is not printing help, it will be interacting with chaincode
	//So start connection
	src.SetConnectionVariables(*flagOrg, *flagUser, *flagPort)
	//src.PrintConnectionVariables()
	clientConnection, err := src.NewGrpcConnection()
	if err != nil {
		panic(err)
	}
	defer clientConnection.Close()

	gw, err := src.DefaultGateway(clientConnection)
	if err != nil {
		panic(err)
	}
	defer gw.Close()
	contract := src.SmartContract(gw)

	//We now check which chaincode function is being called
	if flag.Arg(0) == "storekey" {
		checkEnoughArgs(2)
		storekey(contract, flag.Arg(1))
	} else if flag.Arg(0) == "getkey" {
		checkEnoughArgs(2)
		getkey(contract, flag.Arg(1))
	} else if flag.Arg(0) == "createp" {
		createp(contract)
	} else if flag.Arg(0) == "updatep" {
		checkEnoughArgs(9)
		updatep(contract, flag.Args())
	} else if flag.Arg(0) == "setfillp" {
		checkEnoughArgs(3)
		setfillp(contract, flag.Arg(1), flag.Arg(2))
	} else if flag.Arg(0) == "readp" {
		checkEnoughArgs(2)
		readp(contract, flag.Arg(1))
	} else if flag.Arg(0) == "sharep" {
		checkEnoughArgs(3)
		sharep(contract, flag.Arg(1), flag.Arg(2))
	} else if flag.Arg(0) == "sharedto" {
		checkEnoughArgs(2)
		sharedto(contract, flag.Arg(1))
	} else if flag.Arg(0) == "deletep" {
		checkEnoughArgs(2)
		deletep(contract, flag.Arg(1))
	} else if flag.Arg(0) == "readeradd" {
		checkEnoughArgs(2)
		readeradd(contract, flag.Arg(1))
	} else if flag.Arg(0) == "readerall" {
		checkEnoughArgs(1)
		readerall(contract)
	} else if flag.Arg(0) == "reportgen" {
		checkEnoughArgs(2)
		reportgen(contract, flag.Arg(1))
	} else if flag.Arg(0) == "reportread" {
		checkEnoughArgs(2)
		reportread(contract, flag.Arg(1))
	} else if flag.Arg(0) == "test" {
		src.ChainTestMethod(contract)
	} else {
		fmt.Printf("%vInvalid method '%v'. Do './rsa help' for method options.\n", RED, flag.Arg(0))
	}
}

func storekey(contract *client.Contract, username string) {
	err := src.ChainStoreLocalPubkey(contract, username)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%vKey stored successfully for user %v%v\n", GREEN, username, NC)
}

func getkey(contract *client.Contract, username string) {
	out, err := src.ChainRetrievePubkey(contract, username)
	if err != nil {
		panic(err)
	}
	fmt.Print(out)
	fmt.Printf("\n%vKey retrieved successfully for user %v%v\n", GREEN, username, NC)
}

func genkey(username string) {
	src.GenerateUserKeyFiles(username)
	fmt.Printf("%vKey generated successfully for user %v%v\n", GREEN, username, NC)
}

func createp(contract *client.Contract) {
	err := src.ChainCreatePrescription(contract)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("%vCreate Prescription Successful%v\n", GREEN, NC)
	}
}

func updatep(contract *client.Contract, args []string) {
	cmdInput := src.PrescriptionFromCmdArgs(args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8])
	err := src.ChainUpdatePrescription(contract, cmdInput)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("%vUpdate Prescription Successful%v\n", GREEN, NC)
	}
}

func deletep(contract *client.Contract, pid string) {
	err := src.ChainDeletePrescription(contract, pid)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("%vDelete Prescription Successful%v\n", GREEN, NC)
	}
}

func setfillp(contract *client.Contract, pid string, newfill string) {
	newfillInt, err := strconv.Atoi(newfill)
	if err != nil {
		panic(fmt.Errorf("Failed to parse newfill into integer: %v", err))
	}
	err = src.ChainSetfillPrescription(contract, pid, uint8(newfillInt))
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("%vSetfill Prescription Successful%v\n", GREEN, NC)
	}

}

func readp(contract *client.Contract, pid string) {
	prescription, err := src.ChainReadPrescription(contract, pid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("prescription: %v\n", prescription)
}

func sharep(contract *client.Contract, pid string, username string) {
	err := src.ChainSharePrescription(contract, pid, username)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%vShare Prescription Successful%v\n", GREEN, NC)
}

func sharedto(contract *client.Contract, pid string) {
	list, err := src.ChainSharedToList(contract, pid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("list: %v\n", list)
}

func readeradd(contract *client.Contract, username string) {
	err := src.ChainReportAddReader(contract, username)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%vUser %v is now a report reader%v\n", GREEN, username, NC)
}

func reportgen(contract *client.Contract, pid string) {
	err := src.ChainReportEncrypt(contract, pid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%vReports generated successfully%v\n", GREEN, NC)
}

func reportread(contract *client.Contract, username string) {
	err := src.ChainReportView(contract, username)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%vReports displayed successfully%v\n", GREEN, NC)
}

func readerall(contract *client.Contract) {
	them, err := src.ChainReportGetReaders(contract)
	if err != nil {
		panic(err)
	}
	fmt.Printf("them: %v\n", *them)
}

func checkEnoughArgs(expected int) {
	if len(flag.Args()) < expected {
		panic(fmt.Errorf("%vMethod '%v' expected %v arguments, but was only given %v. Do './rsa help' for method options.\n", RED, flag.Arg(0), expected-1, len(flag.Args())-1))
	}
}
