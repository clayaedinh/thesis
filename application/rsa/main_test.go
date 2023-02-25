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

var keyname []string

func BenchmarkGenerateKey(b *testing.B) {
	//Runtime Phase
	for i := 0; i < b.N; i++ {
		new_key := fmt.Sprintf("benchtest%v", i)
		keyname = append(keyname, new_key)
		src.GenerateUserKeyFiles(new_key)
	}

}

// ======================================================================//
// Send Pubkey
// ======================================================================//
func BenchmarkSendKey(b *testing.B) {
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
}

var keyobsname []string
var keys []string

func BenchmarkSendKeyPrepare(b *testing.B) {
	//Runtime Phase
	for i := 0; i < b.N; i++ {
		keyNum := i % len(keyname)
		newobs, newkey := src.PrepareSendPubkey(keyname[keyNum])
		keyobsname = append(keyobsname, newobs)
		keys = append(keys, newkey)
	}
}
func BenchmarkSendKeySubmit(b *testing.B) {
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
}

// ======================================================================//
// Get Pubkey
// ======================================================================//
func BenchmarkGetKey(b *testing.B) {
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
}

var getkeyout []string

func BenchmarkGetKeyEvaluate(b *testing.B) {
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
}
func BenchmarkGetKeyProcess(b *testing.B) {
	//Runtime Phase
	for i := 0; i < b.N; i++ {
		keyNum := i % len(keyname)
		src.ProcessGetPubkey(getkeyout[keyNum])
	}
}

var pids []string

// ======================================================================//
// Create Prescription
// ======================================================================//
func BenchmarkCreate(b *testing.B) {
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
}

var created []string

func BenchmarkCreatePrepare(b *testing.B) {
	//Runtime Phase
	for i := 0; i < b.N; i++ {
		created = append(created, src.PrepareCreatePrescription())
	}
}

func BenchmarkCreateSubmit(b *testing.B) {
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
		pids = append(pids, src.SubmitCreatePrescription(contract, created[i]))
	}
}

// ======================================================================//
// Read Prescription
// ======================================================================//
func BenchmarkRead(b *testing.B) {
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
		src.ReadPrescription(contract, pids[randPIDNum])
	}
}

var readout []string

func BenchmarkReadEvaluate(b *testing.B) {
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
		readout = append(readout, src.EvaluateReadPrescription(contract, pids[randPIDNum]))
	}
}
func BenchmarkReadProcess(b *testing.B) {
	//Runtime Phase
	for i := 0; i < b.N; i++ {
		src.ProcessReadPrescription(readout[i])
	}
}

// ======================================================================//
// Share Prescription
// Non-split shares to doctors
// Split sharess to pharmacists. DO NOT REMOVE EITHER
// ======================================================================//
func BenchmarkShare(b *testing.B) {
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
}

var shareobs []string
var shareenc []string

func BenchmarkSharePrepare(b *testing.B) {
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
}

func BenchmarkShareSubmit(b *testing.B) {
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
}

// ======================================================================//
// Update Prescription
// ======================================================================//
func BenchmarkUpdate(b *testing.B) {
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
		src.UpdatePrescription(contract, pids[pidsNum], &prescription)
	}
}

var updates []string

func BenchmarkUpdatePrepare(b *testing.B) {
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
		updates = append(updates, src.PrepareUpdatePrescription(contract, pids[pidsNum], &prescription))
	}
}

func BenchmarkUpdateSubmit(b *testing.B) {
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
}

// ======================================================================//
// Setfill Prescription
// ======================================================================//
func BenchmarkSetfill(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		pidsNum := i % len(pids)
		src.SetfillPrescription(contract, pids[pidsNum], uint8(rand.Intn(100)))
	}
}

var setfills []string

func BenchmarkSetfillPrepare(b *testing.B) {
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

	for i := 0; i < b.N; i++ {
		pidsNum := i % len(pids)
		setfills = append(setfills, src.PrepareSetfillPrescription(contract, pids[pidsNum], uint8(rand.Intn(100))))
	}
}

func BenchmarkSetfillSubmit(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		pidsNum := i % len(pids)
		src.SubmitSetfillPrescription(contract, pids[pidsNum], setfills[pidsNum])
	}
}

// ======================================================================//
// Report Register
// Has no application stage, so only one benchmark
// ======================================================================//
func BenchmarkReportRegister(b *testing.B) {
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
}

// ======================================================================//
// Report Update
// ======================================================================//
func BenchmarkReportUpdate(b *testing.B) {
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
		src.ReportUpdate(contract, pids[pidsNum])
	}
}

var reportupdate []string

func BenchmarkReportUpdatePrepare(b *testing.B) {
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
		reportupdate = append(reportupdate, src.PrepareReportUpdate(contract, pids[pidsNum]))
	}
}
func BenchmarkReportUpdateSubmit(b *testing.B) {
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
		src.SubmitReportUpdate(contract, pids[pidsNum], reportupdate[pidsNum])
	}
}

// ======================================================================//
// Report Read
// ======================================================================//
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
		src.ReportView(contract)

	}
}

var reportview []string

func BenchmarkReportReadEvaluate(b *testing.B) {
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
		reportview = append(reportview, src.EvaluateReportView(contract))
	}
}
func BenchmarkReportReadProcess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		src.ProcessReportView(reportview[i])

	}
}

// ======================================================================//
// Delete Prescription
// Has no application stage, so only one benchmark
// ======================================================================//
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
		src.DeletePrescription(contract, pids[i])
	}
}
