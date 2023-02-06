package src

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
)

type Prescription struct {
	Id             string  `json:"Id"`
	DrugBrand      string  `json:"DrugBrand"`
	DrugDoseSched  string  `json:"DrugDoseSched"`
	DrugPrice      float64 `json:"DrugPrice"`
	PatientName    string  `json:"PatientName"`
	PatientAddress string  `json:"PatientAddress"`
	PrescriberName string  `json:"PrescriberName"`
	PrescriberNo   string  `json:"PrescriberNo"`
	Notes          string  `json:"Notes"`
	FilledAmount   string  `json:"FilledAmount"`
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

func printPrescriptionAsJSON(prescription *Prescription) {
	pJSON, err := json.Marshal(prescription)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to convert prescription to JSON: %v", err))
	}
	fmt.Printf("pJSON: %v\n", pJSON)
}

/*
func TestPrintPrescriptionJSON(){
	sample := Prescription{
	Id:             "prescription_1",
	DrugBrand: "Paracetamol",
	DrugDoseSched:  "1 tablet / day",
	DrugPrice:      20.00,
	Notes:          "",
	PatientName:    "Juan de la Cruz",
	PatientAddress: "Katipunan Ave, QC",
	PrescriberName: "Doctor Doctor",
	PrescriberNo:   "12345678"},
}
*/
