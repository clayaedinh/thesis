package src

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

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

func encodePrescription(prescription *Prescription) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*prescription)
	if err != nil {
		return nil, fmt.Errorf("error encoding data %v: %v", prescription.Id, err)
	}
	return buf.Bytes(), nil
}

func decodePrescription(encoded []byte) (*Prescription, error) {
	pres := Prescription{}
	enc := gob.NewDecoder(bytes.NewReader(encoded))
	err := enc.Decode(&pres)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return &pres, nil
}

// ===============================================
// Generate Prescription Id
// a random uint64 used as an id
// ===============================================
func genPrescriptionId() uint64 {
	rand.Seed(time.Now().UnixNano())
	pid := rand.Uint64()
	log.Printf("Generated prescription id: %v", pid)
	return pid
}

// ===============================================
// Package Prescription
// gob-encodes, encrypts, and base-64s
// a prescription so that it's ready to be saved
// ===============================================

func packagePrescription(pubkey *rsa.PublicKey, prescription *Prescription) (string, error) {
	// Encode Prescription to Bytes
	encoded, err := encodePrescription(prescription)
	if err != nil {
		return "", fmt.Errorf("Failed to encode prescription: %v", err)
	}
	//Encrypt data with current user's public key
	encrypted, err := encryptBytes(encoded, pubkey)
	if err != nil {
		return "", fmt.Errorf("Failed to encrypt prescription: %v", err)
	}
	//Encode data as base64
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// ===============================================
// Package Prescription
// reverse of package prescription
// ===============================================
func unpackagePrescription(pdata string) (*Prescription, error) {
	decoded, err := base64.StdEncoding.DecodeString(pdata)
	if err != nil {
		return nil, fmt.Errorf("Base64 failed to decrypt prescription: %v", err)
	}

	// read user privkey
	privkey, err := localPrivkey(getCurrentUser())
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieved user private key: %v", err)
	}

	// Decrypt data with current user's public key
	decrypted, err := decryptBytes(decoded, privkey)
	if err != nil {
		return nil, fmt.Errorf("Failed to encrypt prescription: %v", err)
	}

	// Decode Prescription to Bytes
	prescription, err := decodePrescription(decrypted)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode prescription: %v", err)
	}
	return prescription, nil
}

func PrescriptionFromCmdArgs(pid string, brand string, dosage string, patientName string, patientAddress string,
	prescriberName string, prescriberNo string, piecesTotal string) *Prescription {

	uintpid, err := strconv.ParseUint(pid, 10, 64)

	prescriberNoConv, err := strconv.Atoi(prescriberNo)
	if err != nil {
		panic(fmt.Errorf("Failed to parse prescriber number into integer: %v", err))
	}
	piecesTotalConv, err := strconv.Atoi(piecesTotal)
	if err != nil {
		panic(fmt.Errorf("Failed to parse pieces total into integer: %v", err))
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
