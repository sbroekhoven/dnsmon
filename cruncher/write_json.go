package cruncher

import (
	"encoding/json"
	"os"
)

func WriteJSON(file string, data Domain) (int, error) {
	json, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return 0, err
	}

	f, err := os.OpenFile(file, os.O_CREATE, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	n, err := f.Write(json)
	if err != nil {
		return 0, err
	}
	return n, nil
}
