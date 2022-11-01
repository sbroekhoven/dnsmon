package cruncher

import (
	"encoding/json"
	"os"
)

func ReadJSON(file string) (Domain, error) {
	var domaindata Domain
	domainDataFile, err := os.Open(file)
	if err != nil {
		return domaindata, err
	}
	defer domainDataFile.Close()
	jsonParser := json.NewDecoder(domainDataFile)
	jsonParser.Decode(&domaindata)
	return domaindata, err
}
