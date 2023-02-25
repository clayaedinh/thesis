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

func Benchmark_Connect(b *testing.B) {
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
}

func Benchmark_Create(b *testing.B) {
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

var b64prescriptions []string

func Benchmark_Create_Prepare(b *testing.B) {
	//Runtime Phase
	for i := 0; i < b.N; i++ {
		b64prescriptions = append(b64prescriptions, src.PrepareCreatePrescription())
	}
}

func Benchmark_Create_Submit(b *testing.B) {
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
}

func Benchmark_Read(b *testing.B) {
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
}

var pdata []string

func Benchmark_Read_Evaluate(b *testing.B) {
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
}

func Benchmark_Read_Process(b *testing.B) {
	for i := 0; i < b.N; i++ {
		src.ProcessReadPrescription(pdata[i])
	}
}

func Benchmark_Share_ToDoctors(b *testing.B) {
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
func Benchmark_Share_ToPharmas(b *testing.B) {
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
}

func Benchmark_Update(b *testing.B) {
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

func Benchmark_Update_Prepare(b *testing.B) {
	// RANDOM PRESCRIPTION VARIABLES
	brands := []string{"amoxicillin", "azithromycin", "penicillin", "epinephrine", "aspirin", "insulin", "vitamin D", "paracetamol", "oxytocin", "aluminum hyrdoxide"}
	doses := []string{"once daily", "twice daily", "thrice daily", "4x daily", "weekly", "once only", "twice only", "MWF", "monthly", "when needed"}
	patients := []string{"Aedin", "Sean", "Lance", "Paolo", "Felizia", "Roswold", "Martina", "Christine", "Carmina", "Fide", "Migo", "Leanne", "Haze", "Kyle"}
	doctors := []string{"Dr. Pulmano", "Dr. Tamayo", "Dr. Pangan", "Dr. Sugay", "Dr. Rodrigo", "Dr. Diy", "Dr. Abu", "Dr. Casano", "Dr. Estuar", "Dr. Montalan", "Dr. Jongko"}
	addrs := []string{"Las Pinas", "Quezon City", "Pasay", "Naga", "Paranaque", "Pasig", "Taguig", "Muntinlupa", "Marikina", "Dasmarinas", "Santa Rosa"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
		updates = append(updates, src.PrepareUpdatePrescription(&prescription))
	}
}

func Benchmark_Update_Submit(b *testing.B) {
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

func Benchmark_Setfill(b *testing.B) {
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
}

var setfills []string

func Benchmark_Setfill_Prepare(b *testing.B) {
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
}
func Benchmark_Setfill_Submit(b *testing.B) {
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
}

func Benchmark_Delete(b *testing.B) {
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

func Benchmark_ReportRead(b *testing.B) {
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

var reports []string

func Benchmark_ReportRead_Evaluate(b *testing.B) {
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
}

func Benchmark_ReportRead_Process(b *testing.B) {
	for i := 0; i < b.N; i++ {
		src.ProcessReportView(reports[i])
	}
}
