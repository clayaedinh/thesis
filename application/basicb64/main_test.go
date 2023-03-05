package main

import (
	"fmt"
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

func BenchmarkStandard(b *testing.B) {
	fmt.Println("BENCHMARK TEST -- STANDARD METHODS")
	b.Run("Connect", func(b *testing.B) {
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
		src.SmartContract(gw)
	})

	var pids []string

	b.Run("Create", func(b *testing.B) {
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
			pids = append(pids, src.CreatePrescription(contract))
		}
	})

	b.Run("Read", func(b *testing.B) {
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
			src.ReadPrescription(contract, pids[randPIDNum])
		}
	})

	b.Run("SharetoDoctors", func(b *testing.B) {
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
			src.SharePrescription(contract, pids[pidsNum], "user0001")
		}
	})

	b.Run("SharetoPharmas", func(b *testing.B) {
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
			src.SharePrescription(contract, pids[pidsNum], "user0003")
		}
	})

	b.Run("Update", func(b *testing.B) {
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

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pidsNum := i % len(pids)
			src.UpdatePrescription(contract, pids[pidsNum], &prescription)
		}
	})

	b.Run("Setfill", func(b *testing.B) {
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
			src.SetfillPrescription(contract, pids[pidsNum], uint8(rand.Intn(100)))
		}
	})

	b.Run("ReportRead", func(b *testing.B) {
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
			src.ReportView(contract)
		}
	})

	b.Run("Delete", func(b *testing.B) {
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
	})
}

func BenchmarkSplit(b *testing.B) {
	fmt.Println("BENCHMARK TEST -- SPLIT METHODS")
	var pids []string
	var b64prescriptions []string
	b.Run("CreatePrepare", func(b *testing.B) {
		//Runtime Phase
		for i := 0; i < b.N; i++ {
			b64prescriptions = append(b64prescriptions, src.PrepareCreatePrescription())
		}
	})
	b.Run("CreateSubmit", func(b *testing.B) {
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
			pids = append(pids, src.SubmitCreatePrescription(contract, b64prescriptions[i]))
		}
	})
	var pdata []string

	b.Run("ReadEvaluate", func(b *testing.B) {
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
			pdata = append(pdata, src.EvaluateReadPrescription(contract, pids[randPIDNum]))
		}
	})

	b.Run("ReadProcess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			src.ProcessReadPrescription(pdata[i])
		}
	})
	b.Run("SharetoDoctors", func(b *testing.B) {
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
			src.SharePrescription(contract, pids[pidsNum], "user0001")
		}
	})

	b.Run("SharetoPharmas", func(b *testing.B) {
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
			src.SharePrescription(contract, pids[pidsNum], "user0003")
		}
	})
	var updates []string
	b.Run("UpdatePrepare", func(b *testing.B) {
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
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			updates = append(updates, src.PrepareUpdatePrescription(&prescription))
		}
	})

	b.Run("UpdateSubmit", func(b *testing.B) {
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
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pidsNum := i % len(pids)
			src.SubmitUpdatePrescription(contract, pids[pidsNum], updates[pidsNum])
		}
	})

	var setfills []string
	b.Run("SetfillPrepare", func(b *testing.B) {
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
			setfills = append(setfills, src.PrepareSetfillPrescription(contract, pids[pidsNum], uint8(rand.Intn(100))))
		}
	})
	b.Run("SetfillSubmit", func(b *testing.B) {
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
			src.SubmitSetfillPrescription(contract, pids[pidsNum], setfills[pidsNum])
		}
	})

	var reports []string
	b.Run("ReportReadEvaluate", func(b *testing.B) {
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
			reports = append(reports, src.EvaluateReportView(contract))
		}
	})

	b.Run("ReportReadProcess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			src.ProcessReportView(reports[i])
		}
	})
	b.Run("Delete", func(b *testing.B) {
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
	})
}

func BenchmarkPrescriptionAmountAndReportRead(b *testing.B) {
	var pids []string
	b.Run("Create", func(b *testing.B) {
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
			pids = append(pids, src.CreatePrescription(contract))
		}
	})
	b.Run("ReportRead", func(b *testing.B) {
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
			src.ReportView(contract)
		}
	})
}
