package src

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
)

func packageStringSlice(strings *[]string) (string, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*strings)
	if err != nil {
		return "", fmt.Errorf("Failed to gob the string slice: %v", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func packagePrescriptionSet(pset *map[string]string) (string, error) {

	// STEP 1: Gob-Encode
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(*pset)
	if err != nil {
		return "", fmt.Errorf("Failed to gob the prescription set: %v", err)
	}

	// STEP 2: Base-64 it
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func unpackagePrescriptionSet(packaged string) (*map[string]string, error) {
	rawgob, err := base64.StdEncoding.DecodeString(packaged)
	if err != nil {
		return nil, err
	}
	pset := make(map[string]string)
	enc := gob.NewDecoder(bytes.NewReader(rawgob))
	err = enc.Decode(&pset)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return &pset, nil
}
