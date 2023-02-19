package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/clayaedinh/thesis/application/basicb64/src"
)

/*
WARNING

It is best to use -benchtime=x instead of -benchtime=s
So that the number of iterations (and prescriptions made) is consistent

*/

var pids []string

func BenchmarkCreatePrescription(b *testing.B) {

	// Connection Phase
	src.SetConnectionVariables("org1", "user0002", "localhost:7051")
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

	//Runtime Phase
	for i := 0; i < b.N; i++ {
		new_pid, _ := src.ChainCreatePrescription(contract)
		pids = append(pids, new_pid)
	}
}

// please run BenchmarkCreatePrescription prior to this.
func BenchmarkReadPrescription(b *testing.B) {
	// Connection Phase
	src.SetConnectionVariables("org1", "user0002", "localhost:7051")
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

	//Runtime Phase
	for i := 0; i < b.N; i++ {
		rand.Seed(time.Now().UnixNano())
		randPIDNum := rand.Intn(len(pids) - 1)
		src.ChainReadPrescription(contract, pids[randPIDNum])
	}
}

func BenchmarkSharePrescription(b *testing.B) {
	// Connection Phase
	src.SetConnectionVariables("org1", "user0002", "localhost:7051")
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

	for i := 0; i < b.N; i++ {
		pidsNum := i % len(pids)
		src.ChainSharePrescription(contract, pids[pidsNum], "user0001")
	}
}

func BenchmarkUpdatePrescription(b *testing.B) {
	// Connection Phase
	src.SetConnectionVariables("org1", "user0001", "localhost:7051")
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

	// RANDOM PRESCRIPTION VARIABLES
	brands := []string{"amoxicillin", "azithromycin", "penicillin", "epinephrine", "aspirin", "insulin", "vitamin D", "paracetamol", "oxytocin", "aluminum hyrdoxide"}
	doses := []string{"once daily", "twice daily", "thrice daily", "4x daily", "weekly", "once only", "twice only", "MWF", "monthly", "when needed"}
	patients := []string{"Aedin", "Sean", "Lance", "Paolo", "Felizia", "Roswold", "Martina", "Christine", "Carmina", "Fide", "Migo", "Leanne", "Haze", "Kyle"}
	doctors := []string{"Dr. Pulmano", "Dr. Tamayo", "Dr. Pangan", "Dr. Sugay", "Dr. Rodrigo", "Dr. Diy", "Dr. Abu", "Dr. Casano", "Dr. Estuar", "Dr. Montalan", "Dr. Jongko"}
	addrs := []string{"Las Pinas", "Quezon City", "Pasay", "Naga", "Paranaque", "Pasig", "Taguig", "Muntinlupa", "Marikina", "Dasmarinas", "Santa Rosa"}

	for i := 0; i < b.N; i++ {
		pidsNum := i % len(pids)
		prescription := src.Prescription{
			Brand:          brands[i%len(brands)],
			Dosage:         doses[i%len(doses)],
			PatientName:    patients[i%len(patients)],
			PatientAddress: addrs[i%len(addrs)],
			PrescriberName: doctors[i%len(doctors)],
			PrescriberNo:   rand.Uint32(),
			PiecesTotal:    uint8(rand.Intn(100)),
			PiecesFilled:   0,
		}
		src.ChainUpdatePrescription(contract, pids[pidsNum], &prescription)
	}
}

func BenchmarkDelete(b *testing.B) {
	// Connection Phase
	src.SetConnectionVariables("org1", "user0002", "localhost:7051")
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
	for i := 0; i < b.N; i++ {
		if i >= len(pids) {
			break
		}
		src.ChainDeletePrescription(contract, pids[i])
	}
}
