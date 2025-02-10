package main

import (
	"encoding/json"
	"fmt"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func main() {
	raw, err := whois.Whois("degiro.pt")
	if err == nil {
		fmt.Println(raw)
	}

	result, err := whoisparser.Parse(raw)
	if err == nil {

		// Convert domain to JSON
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Printf("Error converting to JSON: %v\n", err)
		}

		// Print JSON
		fmt.Println(string(jsonData))

		// Print the domain status
		fmt.Println(result.Domain.Status)

		// Print the domain expiration date
		fmt.Println(result.Domain.ExpirationDate)

		// Print the registrar name
		// fmt.Println(result.Registrar.Name)

	}
}
