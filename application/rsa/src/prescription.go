package src

import (
	"bytes"
	"encoding/gob"
	"fmt"
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

func EncodePrescription(prescription *Prescription) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*prescription)
	if err != nil {
		return nil, fmt.Errorf("error encoding data %v: %v", prescription.Id, err)
	}
	return buf.Bytes(), nil
}

func DecodePrescription(encoded []byte) (*Prescription, error) {
	pres := Prescription{}
	enc := gob.NewDecoder(bytes.NewReader(encoded))
	err := enc.Decode(&pres)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return &pres, nil
}
