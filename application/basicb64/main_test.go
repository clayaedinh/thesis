package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/clayaedinh/thesis/application/basicb64/src"
)

/*
	func BenchmarkObscureName(b *testing.B) {
		for i := 0; i < b.N; i++ {
			namei := fmt.Sprintf("user%v", i)
			src.ObscureName(namei)
		}
	}
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
