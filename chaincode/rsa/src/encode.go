package src

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func encodePrescriptionSet(pset *map[string]string) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(pset)
	if err != nil {
		return nil, fmt.Errorf("Failed to gob the prescription set: %v", err)
	}
	return buf.Bytes(), nil
}

func decodePrescriptionSet(rawgob []byte) (map[string]string, error) {
	pset := make(map[string]string)
	enc := gob.NewDecoder(bytes.NewReader(rawgob))
	err := enc.Decode(&pset)
	if err != nil {
		return nil, fmt.Errorf("error decoding data : %v", err)
	}
	return pset, nil
}

func encodeStringSlice(strings []string) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(strings)
	if err != nil {
		return nil, fmt.Errorf("Failed to gob the string slice: %v", err)
	}
	return buf.Bytes(), nil
}
