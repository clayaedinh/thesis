package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/clayaedinh/thesis/application/rsa/src"
)

/*
WARNING

It is best to use -benchtime=x instead of -benchtime=s
So that the number of iterations (and prescriptions made) is consistent
*/

/*
NOTICE

In practice, update and setfill time should add report generation time - since reports need to be generated every time
there is an update or setfills

*/

// ======================================================================//
// BENCHMARK STANDARD
// Benchmarks all functions, non-split
// ======================================================================//
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

	var keyname []string

	b.Run("GenerateKey", func(b *testing.B) {
		//Runtime Phase
		for i := 0; i < b.N; i++ {
			new_key := fmt.Sprintf("benchtest%v", i)
			keyname = append(keyname, new_key)
			src.GenerateUserKeyFiles(new_key)
		}
	})
	b.Run("SendKey", func(b *testing.B) {
		// Connection Phase
		src.SetConnectionVariables("org1", "Admin", "localhost:7051")
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
			keyNum := i % len(keyname)
			src.SendPubkey(contract, keyname[keyNum])
		}
	})

	b.Run("GetKey", func(b *testing.B) {
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
			keyNum := i % len(keyname)
			src.GetPubkey(contract, src.ObscureName(keyname[keyNum]))
		}
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

	b.Run("ReportAddReader", func(b *testing.B) {
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
			src.ChainReportAddReader(contract)
		}
	})

	b.Run("ReportUpdate", func(b *testing.B) {
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
			src.ReportUpdate(contract, pids[pidsNum])
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
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("recovered from error")
				}
			}()
			_ = src.DeletePrescription(contract, pids[i])
		}
	})
}

// ======================================================================//
// BENCHMARK SPLIT
// Benchmarks all functions, split
// Into two halves:
// Prepare (Application)
// Submit (Chaincode)
// OR
// Evaluate (Chaincode)
// Process (Application)
// ======================================================================//
func BenchmarkSplit(b *testing.B) {
	fmt.Println("BENCHMARK TEST -- SPLIT METHODS")
	var pids []string
	var keyname []string

	b.Run("GenerateKey", func(b *testing.B) {
		//Runtime Phase
		for i := 0; i < b.N; i++ {
			new_key := fmt.Sprintf("benchtest%v", i)
			keyname = append(keyname, new_key)
			src.GenerateUserKeyFiles(new_key)
		}
	})

	var keyobsname []string
	var keys []string

	b.Run("SendKeyPrepare", func(b *testing.B) {
		//Runtime Phase
		for i := 0; i < b.N; i++ {
			keyNum := i % len(keyname)
			newobs, newkey := src.PrepareSendPubkey(keyname[keyNum])
			keyobsname = append(keyobsname, newobs)
			keys = append(keys, newkey)
		}
	})

	b.Run("SendKeySubmit", func(b *testing.B) {
		// Connection Phase
		src.SetConnectionVariables("org1", "Admin", "localhost:7051")
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
			keyNum := i % len(keyname)
			src.SubmitSendPubkey(contract, keyobsname[keyNum], keys[keyNum])
		}
	})

	var getkeyout []string

	b.Run("GetKeyEvaluate", func(b *testing.B) {
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
			keyNum := i % len(keyname)
			getkeyout = append(getkeyout, src.EvaluateGetPubkey(contract, src.ObscureName(keyname[keyNum])))
		}
	})

	b.Run("GetKeyProcess", func(b *testing.B) {
		//Runtime Phase
		for i := 0; i < b.N; i++ {
			keyNum := i % len(keyname)
			src.ProcessGetPubkey(getkeyout[keyNum])
		}
	})

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

	b.Run("(SharetoDoctors)", func(b *testing.B) {
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

	var shareobs []string
	var shareenc []string

	b.Run("SharePrepare", func(b *testing.B) {
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
			newobs, newenc := src.PrepareSharePrescription(contract, pids[pidsNum], "user0003")
			shareobs = append(shareobs, newobs)
			shareenc = append(shareenc, newenc)
		}
	})

	b.Run("ShareSubmit", func(b *testing.B) {
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
			src.SubmitSharePrescription(contract, pids[pidsNum], shareobs[pidsNum], shareenc[pidsNum])
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
			updates = append(updates, src.PrepareUpdatePrescription(contract, pids[pidsNum], &prescription))
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

	b.Run("ReportAddReader", func(b *testing.B) {
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
			src.ChainReportAddReader(contract)
		}
	})

	var reportupdate []string

	b.Run("ReportUpdatePrepare", func(b *testing.B) {
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
			reportupdate = append(reportupdate, src.PrepareReportUpdate(contract, pids[pidsNum]))
		}
	})
	b.Run("ReportUpdateSubmit", func(b *testing.B) {
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
			src.SubmitReportUpdate(contract, pids[pidsNum], reportupdate[pidsNum])
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
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("recovered from error")
				}
			}()
			_ = src.DeletePrescription(contract, pids[i])
		}
	})
}

func BenchmarkPrescriptionAmountAndReportRead(b *testing.B) {
	var pids []string
	b.Run("(Create)", func(b *testing.B) {
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

	b.Run("(SharetoDoctors)", func(b *testing.B) {
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

	b.Run("(ReportAddReader)", func(b *testing.B) {
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
		src.ChainReportAddReader(contract)
	})

	b.Run("(ReportUpdate)", func(b *testing.B) {
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
			src.ReportUpdate(contract, pids[pidsNum])
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
