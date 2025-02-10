package cruncher

import (
	"dnsmon/alerting"
	"dnsmon/config"
	"fmt"
	"log"
	"reflect"
)

func Compare(alert config.Alerting, old Domain, new Domain) (bool, error) {
	eq := true

	// Check if new.Serial is set (assuming 0 is not a valid serial number)
	if new.Serial == 0 {
		message := fmt.Sprintf("⚠️ The serial in SOA record of domain %s is not found.", new.Domainname)
		log.Println(message)
		alerting.SendAlerts(alert, message)
	}

	// Compare domain zone serial
	if old.Serial != new.Serial {
		eq = false
		message := fmt.Sprintf("⚠️ The serial in SOA record of domain %s is changed.", new.Domainname)
		log.Println(message)
		alerting.SendAlerts(alert, message)
	}

	// Compare nameservers
	if !stringSlicesEqual(old.Nameservers, new.Nameservers) {
		eq = false
		message := fmt.Sprintf("⚠️ The nameservers of domain %s are changed.", new.Domainname)
		log.Println(message)
		alerting.SendAlerts(alert, message)
	}

	// Compare mailservers
	if !stringSlicesEqual(old.Mailservers, new.Mailservers) {
		eq = false
		message := fmt.Sprintf("⚠️ The MX records of domain %s are changed.", new.Domainname)
		log.Println(message)
		alerting.SendAlerts(alert, message)
	}

	// Compare SPF
	if old.SPFRecord != new.SPFRecord {
		eq = false
		message := fmt.Sprintf("⚠️ The SPF record of domain %s is changed.", new.Domainname)
		log.Println(message)
		alerting.SendAlerts(alert, message)
	}

	// Compare SPF
	if old.DMARCRecord != new.DMARCRecord {
		eq = false
		message := fmt.Sprintf("⚠️ The DMARC record of domain %s is changed.", new.Domainname)
		log.Println(message)
		alerting.SendAlerts(alert, message)
	}

	// Compare records
	// TODO: Needs some work
	if (len(old.Records) > 0) || (len(new.Records) > 0) {
		ref := reflect.DeepEqual(old.Records, new.Records)
		if !ref {
			eq = false
			message := fmt.Sprintf("⚠️ Some records of domain %s are changed.", new.Domainname)
			log.Println(message)
			alerting.SendAlerts(alert, message)
		}
	}
	// message := fmt.Sprintf("⚠️ Testing domain %s .", new.Domainname)
	// alerting.SendAlerts(alert, message)

	return eq, nil
}

// stringSlicesEqual function to check if slises are equal
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
