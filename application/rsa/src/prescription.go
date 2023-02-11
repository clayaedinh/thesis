package src

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/rand"
)

type Prescription struct {
	Id             string `json:"Id"`
	DrugBrand      string `json:"DrugBrand"`
	DrugDoseSched  string `json:"DrugDoseSched"`
	DrugPrice      string `json:"DrugPrice"`
	PatientName    string `json:"PatientName"`
	PatientAddress string `json:"PatientAddress"`
	PrescriberName string `json:"PrescriberName"`
	PrescriberNo   string `json:"PrescriberNo"`
	FilledAmount   string `json:"FilledAmount"`
}

type PrescriptionCmdInput struct {
	DrugBrand      string `json:"DrugBrand"`
	DrugDoseSched  string `json:"DrugDoseSched"`
	DrugPrice      string `json:"DrugPrice"`
	PatientName    string `json:"PatientName"`
	PatientAddress string `json:"PatientAddress"`
	PrescriberName string `json:"PrescriberName"`
	PrescriberNo   string `json:"PrescriberNo"`
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

func genPrescriptionId() string {
	pid := fmt.Sprint(rand.Intn(math.MaxInt32))
	log.Printf("Generated prescription id: %v", pid)
	return pid
}

func PrescriptionFromCmdArgs(drugBrand string, drugSched string, drugPrice string,
	patientName string, patientAddress string, prescriberName string, prescriberNo string) *Prescription {

	prescription := Prescription{
		Id:             genPrescriptionId(),
		DrugBrand:      drugBrand,
		DrugDoseSched:  drugSched,
		DrugPrice:      drugPrice,
		PatientName:    patientName,
		PatientAddress: patientAddress,
		PrescriberName: prescriberName,
		PrescriberNo:   prescriberNo,
	}
	return &prescription

}
