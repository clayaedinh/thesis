package src

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"strconv"
)

// ============================================================ //
// Prescription
// ============================================================ //

type Prescription struct {
	Brand          string `json:"Brand"`
	Dosage         string `json:"Dosage"`
	PatientName    string `json:"PatientName"`
	PatientAddress string `json:"PatientAddress"`
	PrescriberName string `json:"PrescriberName"`
	PrescriberNo   uint32 `json:"PrescriberNo"`
	PiecesTotal    uint8  `json:"AmountTotal"`
	PiecesFilled   uint8  `json:"AmountFilled"`
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
		return "", fmt.Errorf("error encoding prescription data: %v", err)
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

func unpackageStringSlice(b64slice string) (*[]string, error) {
	gobslice, err := base64.StdEncoding.DecodeString(b64slice)
	if err != nil {
		return nil, err
	}
	var strings []string
	enc := gob.NewDecoder(bytes.NewReader(gobslice))
	err = enc.Decode(&strings)
	if err != nil {
		return nil, fmt.Errorf("error decoding string slice : %v", err)
	}
	return &strings, nil
}

// ============================================================ //
// Prescription from Command Arguments
// ============================================================ //

func PrescriptionFromCmdArgs(brand string, dosage string, patientName string, patientAddress string,
	prescriberName string, prescriberNo string, piecesTotal string) *Prescription {

	prescriberNoConv, err := strconv.Atoi(prescriberNo)
	if err != nil {
		panic(fmt.Errorf("failed to parse prescriber number into integer: %v", err))
	}
	piecesTotalConv, err := strconv.Atoi(piecesTotal)
	if err != nil {
		panic(fmt.Errorf("failed to parse pieces total into integer: %v", err))
	}

	prescription := Prescription{
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
