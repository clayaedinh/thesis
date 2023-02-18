package src

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// ============================================================ //
// Prescription
// ============================================================ //

type Prescription struct {
	Id             uint64 `json:"Id"`
	Brand          string `json:"Brand"`
	Dosage         string `json:"Dosage"`
	PatientName    string `json:"PatientName"`
	PatientAddress string `json:"PatientAddress"`
	PrescriberName string `json:"PrescriberName"`
	PrescriberNo   uint32 `json:"PrescriberNo"`
	PiecesTotal    uint8  `json:"AmountTotal"`
	PiecesFilled   uint8  `json:"AmountFilled"` // in terms of percentage
}

func genPrescriptionId() uint64 {
	rand.Seed(time.Now().UnixNano())
	pid := rand.Uint64()
	log.Printf("Generated prescription id: %v", pid)
	return pid
}

// ============================================================ //
// ENCODING
// ============================================================ //

func obscureName(username string) string {
	raw := sha256.Sum256([]byte(username))
	return hex.EncodeToString(raw[:])
}

func packagePrescription(prescription *Prescription) (string, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*prescription)
	if err != nil {
		return "", fmt.Errorf("error encoding data %v: %v", prescription.Id, err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func unpackagePrescription(b64prescription string) (*Prescription, error) {
	encoded, err := base64.StdEncoding.DecodeString(b64prescription)
	if err != nil {
		return nil, err
	}
	prescription := Prescription{}
	enc := gob.NewDecoder(bytes.NewReader(encoded))
	err = enc.Decode(&prescription)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return &prescription, nil
}

// ============================================================ //
// Prescription from Command Arguments
// ============================================================ //

func PrescriptionFromCmdArgs(pid string, brand string, dosage string, patientName string, patientAddress string,
	prescriberName string, prescriberNo string, piecesTotal string) *Prescription {

	uintpid, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		panic(fmt.Errorf("failed to parse patient ID into integer: %v", err))
	}

	prescriberNoConv, err := strconv.Atoi(prescriberNo)
	if err != nil {
		panic(fmt.Errorf("failed to parse prescriber number into integer: %v", err))
	}
	piecesTotalConv, err := strconv.Atoi(piecesTotal)
	if err != nil {
		panic(fmt.Errorf("failed to parse pieces total into integer: %v", err))
	}

	prescription := Prescription{
		Id:             uintpid,
		Brand:          brand,
		Dosage:         dosage,
		PatientName:    patientName,
		PatientAddress: patientAddress,
		PrescriberName: prescriberName,
		PrescriberNo:   uint32(prescriberNoConv),
		PiecesTotal:    uint8(piecesTotalConv),
		PiecesFilled:   0,
	}
	return &prescription
}
