package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/clayaedinh/thesis/application/basicb64/src"
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
	fmt.Printf("%vUsage%v: ./basicb64 %v[-options] %v<Method> %v<Method Args>\n", GREEN, NC, PURPLE, CYAN, NC)
	fmt.Println("")
	fmt.Printf("%vAvailable Options%v:\n", GREEN, NC)
	fmt.Println("")
	fmt.Printf("./basicb64 %v-user=%vstring\n", PURPLE, NC)
	fmt.Println(FLAG_H_USER)
	fmt.Println("")
	fmt.Printf("./basicb64 %v-org=%vstring\n", PURPLE, NC)
	fmt.Println(FLAG_H_ORG)
	fmt.Println("")
	fmt.Printf("./basicb64 %v-port=%vlocalhost:port\n", PURPLE, NC)
	fmt.Println(FLAG_H_PORT)
	fmt.Println("")
	fmt.Printf("%vAvailable Methods (must be AFTER options)%v:\n", GREEN, NC)
	fmt.Printf("./basicb64 %vcreatep%v\n", CYAN, NC)
	fmt.Printf("./basicb64 %vsharep%v <pid> <username>\n", CYAN, NC)
	fmt.Printf("./basicb64 %vreadp%v <id>\n", CYAN, NC)
	fmt.Printf("./basicb64 %vupdatep%v <brand> <dosage> <patient_name> <patient_address> <doctor_name> <doctor_prc> <pieces_total>\n", CYAN, NC)
	fmt.Printf("./basicb64 %vsetfillp%v <pid> <newfill>\n", CYAN, NC)
	fmt.Printf("./basicb64 %vdeletep%v <pid>\n", CYAN, NC)
	fmt.Printf("./basicb64 %vreportread%v <username>\n", CYAN, NC)
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
	//Connect to chaincode:
	src.SetConnectionVariables(*flagOrg, *flagUser, *flagPort)
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
	if flag.Arg(0) == "createp" {
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
	} else if flag.Arg(0) == "deletep" {
		checkEnoughArgs(2)
		deletep(contract, flag.Arg(1))
	} else if flag.Arg(0) == "reportread" {
		reportread(contract)
	} else if flag.Arg(0) == "getCreate" {
		getCreateOutput(contract)
	} else if flag.Arg(0) == "getUpdate" {
		getUpdateOutput(contract)
	} else if flag.Arg(0) == "getSetfill" {
		checkEnoughArgs(2)
		getSetfillOutput(contract, flag.Arg(2))
	} else {
		fmt.Printf("%vInvalid method '%v'. Do './rsa help' for method options.\n", RED, flag.Arg(0))
	}
}

func createp(contract *client.Contract) {
	pid := src.CreatePrescription(contract)
	fmt.Printf("%vCreate Prescription Successful. PID: %v.%v\n", GREEN, pid, NC)
}
func readp(contract *client.Contract, pid string) {
	prescription := src.ReadPrescription(contract, pid)
	fmt.Printf("Prescription: %v\n", prescription)
}

func sharep(contract *client.Contract, pid string, username string) {
	src.SharePrescription(contract, pid, username)
	fmt.Printf("%vShare Prescription Successful%v\n", GREEN, NC)
}

func updatep(contract *client.Contract, args []string) {
	cmdInput := src.PrescriptionFromCmdArgs(args[2], args[3], args[4], args[5], args[6], args[7], args[8])
	src.UpdatePrescription(contract, args[1], cmdInput)
	fmt.Printf("%vUpdate Prescription Successful%v\n", GREEN, NC)
}

func setfillp(contract *client.Contract, pid string, newfill string) {
	newfillInt, err := strconv.Atoi(newfill)
	if err != nil {
		panic(fmt.Errorf("failed to parse newfill into integer: %v", err))
	}
	src.SetfillPrescription(contract, pid, uint8(newfillInt))
	fmt.Printf("%vSetfill Prescription Successful%v\n", GREEN, NC)
}

func deletep(contract *client.Contract, pid string) {
	err := src.ChainDeletePrescription(contract, pid)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("%vDelete Prescription Successful%v\n", GREEN, NC)
	}
}

func reportread(contract *client.Contract) {
	reports := src.ReportView(contract)
	for _, x := range *reports {
		fmt.Printf("x: %v\n", x)
	}
}

func getCreateOutput(contract *client.Contract) {
	createb64 := src.PrepareCreatePrescription()
	testpid := src.SubmitCreatePrescription(contract, createb64)
	fmt.Printf("createb64: %v\n", createb64)
	fmt.Printf("testpid: %v\n", testpid)
}

func getUpdateOutput(contract *client.Contract) {
	prescription := src.Prescription{
		Brand:          "DRUG BRAND",
		Dosage:         "DRUG DOSAGE",
		PatientName:    "PATIENT NAME",
		PatientAddress: "PATIENT ADDR",
		PrescriberName: "PRESC NAME",
		PrescriberNo:   1234567,
		PiecesTotal:    100,
		PiecesFilled:   0,
	}
	updateb64 := src.PrepareUpdatePrescription(&prescription)
	fmt.Printf("updateb64: %v\n", updateb64)
}
func getSetfillOutput(contract *client.Contract, testpid string) {
	setfillb64 := src.PrepareSetfillPrescription(contract, testpid, 100)
	fmt.Printf("setfillb64: %v\n", setfillb64)
}

func checkEnoughArgs(expected int) {
	if len(flag.Args()) < expected {
		panic(fmt.Errorf("%vmethod '%v' expected %v arguments, but was only given %v. Do './rsa help' for method options", RED, flag.Arg(0), expected-1, len(flag.Args())-1))
	}
}

/*
	generates 10 random of the following:

- pid
- blank prescription with pid
- update prescription

also generates:
- obscured users (user0001 to user0003)
*/
/*
func inputgen() {
	fmt.Printf("%v Prescription Ids %v\n", YELLOW, NC)
	var pids [10]uint64
	for i := 0; i < 10; i++ {
		pids[i] = src.GenPrescriptionId()
		fmt.Printf("pids[i]: %v\n", pids[i])
	}
	fmt.Printf("%v Blank Prescriptions %v\n", YELLOW, NC)

}
*/
