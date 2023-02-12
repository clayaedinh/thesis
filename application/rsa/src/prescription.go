package src

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"strconv"
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

func genPrescriptionId() uint64 {
	pid := rand.Uint64()
	log.Printf("Generated prescription id: %v", pid)
	return pid
}

type PrescriptionCmdInput struct {
	DrugBrand      string `json:"DrugBrand"`
	Dosage         string `json:"Dosage"`
	PatientName    string `json:"PatientName"`
	PatientAddress string `json:"PatientAddress"`
	PrescriberName string `json:"PrescriberName"`
	PrescriberNo   uint32 `json:"PrescriberNo"`
	PiecesTotal    uint8  `json:"AmountTotal"`
}

func PrescriptionFromCmdArgs(brand string, dosage string, patientName string, patientAddress string,
	prescriberName string, prescriberNo string, piecesTotal string) *Prescription {
	prescriberNoConv, err := strconv.Atoi(prescriberNo)
	if err != nil {
		panic(fmt.Errorf("Failed to parse prescriber number into integer: %v", err))
	}
	piecesTotalConv, err := strconv.Atoi(piecesTotal)
	if err != nil {
		panic(fmt.Errorf("Failed to parse pieces total into integer: %v", err))
	}

	prescription := Prescription{
		Id:             genPrescriptionId(),
		Brand:          brand,
		Dosage:         dosage,
		PatientName:    patientName,
		PatientAddress: patientAddress,
		PrescriberName: prescriberName,
		PrescriberNo:   uint32(prescriberNoConv),
		PiecesTotal:    uint8(piecesTotalConv),
	}
	return &prescription

}
