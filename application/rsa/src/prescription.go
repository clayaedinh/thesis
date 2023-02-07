package src

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
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

type PrescriptionCmdInput struct {
	DrugBrand      string  `json:"DrugBrand"`
	DrugDoseSched  string  `json:"DrugDoseSched"`
	DrugPrice      float64 `json:"DrugPrice"`
	PatientName    string  `json:"PatientName"`
	PatientAddress string  `json:"PatientAddress"`
	PrescriberName string  `json:"PrescriberName"`
	PrescriberNo   string  `json:"PrescriberNo"`
	Notes          string  `json:"Notes"`
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

func PrescriptionFromCmdInput(jsonstring string) (*Prescription, error) {
	var cmdInput *PrescriptionCmdInput
	err := json.Unmarshal([]byte(jsonstring), cmdInput)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse json string into prescription: %v", err)
	}

	prescription_id := "prescription_" + string(rand.Intn(math.MaxInt32))
	prescription := Prescription{
		Id:             prescription_id,
		DrugBrand:      cmdInput.DrugBrand,
		DrugDoseSched:  cmdInput.DrugDoseSched,
		DrugPrice:      cmdInput.DrugPrice,
		PatientName:    cmdInput.PatientName,
		PatientAddress: cmdInput.PatientName,
		PrescriberName: cmdInput.PrescriberName,
		PrescriberNo:   cmdInput.PrescriberNo,
		Notes:          cmdInput.Notes,
		FilledAmount:   "none",
	}
	return &prescription, nil

}

func TestPrintPrescriptionJSON() {
	sample := PrescriptionCmdInput{
		DrugBrand:      "Paracetamol",
		DrugDoseSched:  "1 tablet / day",
		DrugPrice:      20.00,
		PatientName:    "Juan de la Cruz",
		PatientAddress: "Katipunan Ave, QC",
		PrescriberName: "Doctor Doctor",
		PrescriberNo:   "12345678",
		Notes:          ""}
	fmt.Println("RAW STRUCT")
	fmt.Printf("sample: %v\n", sample)
	fmt.Println("JSON")

	pJSON, err := json.Marshal(sample)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to convert prescription to JSON: %v", err))
	}
	fmt.Printf("pJSON: %v\n", string(pJSON))

}

func PrintEmptyPrescriptionJSON() string {
	sample := PrescriptionCmdInput{
		DrugBrand:      "",
		DrugDoseSched:  "",
		DrugPrice:      0,
		PatientName:    "",
		PatientAddress: "",
		PrescriberName: "",
		PrescriberNo:   "",
		Notes:          ""}

	pJSON, err := json.Marshal(sample)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to convert prescription to JSON: %v", err))
	}
	return string(pJSON)
}
