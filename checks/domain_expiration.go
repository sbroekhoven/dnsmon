package checks

import (
	"errors"
	"fmt"
	"time"

	"github.com/openrdap/rdap"
)

// GetSerial function to resolve SOA record from domain and return only the serial for now
func getExpirationDate(domain *rdap.Domain) (time.Time, error) {
	fmt.Printf("Domain Status: %v\n", domain.Status)

	for _, event := range domain.Events {
		if event.Action == "expiration" {

			t, err := time.Parse(time.RFC3339, event.Date)
			if err != nil {
				return time.Time{}, errors.New("error parsing time")
			} else {
				fmt.Printf("Parsed time: %v\n", t)
				fmt.Printf("Status: %v\n", domain.Status)
			}

			return t, nil
		}
	}
	return time.Time{}, errors.New("expiration date not found")
}
