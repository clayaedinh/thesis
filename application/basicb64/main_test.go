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

	b.ResetTimer()
	//Runtime Phase
	for i := 0; i < b.N; i++ {
		var b64prescription string
		var new_pid string
		b.Run("AppPrepare", func(b *testing.B) {
			b64prescription = src.PrepareCreatePrescription()
		})
		b.Run("ChainSubmit", func(b *testing.B) {
			new_pid = src.SubmitCreatePrescription(contract, b64prescription)
		})
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

	b.ResetTimer()
	//Runtime Phase
	for i := 0; i < b.N; i++ {
		rand.Seed(time.Now().UnixNano())
		randPIDNum := rand.Intn(len(pids) - 1)
		var pdata string
		b.Run("ChainEvaluate", func(b *testing.B) {
			pdata = src.EvaluateReadPrescription(contract, pids[randPIDNum])
		})
		b.Run("AppProcess", func(b *testing.B) {
			src.ProcessReadPrescription(pdata)
		})
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pidsNum := i % len(pids)
		b.Run("ToDoctors", func(b *testing.B) {
			src.SharePrescription(contract, pids[pidsNum], "user0001")
		})
		b.Run("ToPharmas", func(b *testing.B) {
			src.SharePrescription(contract, pids[pidsNum], "user0003")
		})
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pidsNum := i % len(pids)

		var b64prescription string
		b.Run("AppPrepare", func(b *testing.B) {
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
			b64prescription = src.PrepareUpdatePrescription(&prescription)
		})
		b.Run("ChainSubmit", func(b *testing.B) {
			src.SubmitUpdatePrescription(contract, pids[pidsNum], b64prescription)
		})
	}
}

func BenchmarkSetfillPrescription(b *testing.B) {
	// Connection Phase
	src.SetConnectionVariables("org2", "user0003", "localhost:9051")
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

	b.ResetTimer()
	// RANDOM PRESCRIPTION VARIABLES
	for i := 0; i < b.N; i++ {
		pidsNum := i % len(pids)
		var b64prescription string
		b.Run("AppPrepare", func(b *testing.B) {
			b64prescription = src.PrepareSetfillPrescription(contract, pids[pidsNum], uint8(rand.Intn(100)))
		})
		b.Run("ChainSubmit", func(b *testing.B) {
			src.SubmitSetfillPrescription(contract, pids[pidsNum], b64prescription)
		})
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i >= len(pids) {
			break
		}
		src.ChainDeletePrescription(contract, pids[i])
	}
}

func BenchmarkReportRead(b *testing.B) {
	// Connection Phase
	src.SetConnectionVariables("org1", "user0004", "localhost:7051")
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var b64report string
		b.Run("ChainEvaluate", func(b *testing.B) {
			b64report = src.EvaluateReportView(contract)
		})
		b.Run("AppProcess", func(b *testing.B) {
			src.ProcessReportView(b64report)
		})
	}
}
