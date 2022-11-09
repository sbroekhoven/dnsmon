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

	// Compare domain zone serial
	if old.Serial != new.Serial {
		eq = false
		message := fmt.Sprintf("The serial of domain %s is changed.", new.Domainname)
		log.Println(message)
		err := alerting.Discord(alert.DiscordWebhookURL, alert.DiscordUsername, message)
		if err != nil {
			log.Println(err.Error())
		}
	}

	// Compare nameservers
	if !stringSlicesEqual(old.Nameservers, new.Nameservers) {
		eq = false
		message := fmt.Sprintf("The nameservers of domain %s are changed.", new.Domainname)
		log.Println(message)
		err := alerting.Discord(alert.DiscordWebhookURL, alert.DiscordUsername, message)
		if err != nil {
			log.Println(err.Error())
		}
	}

	// Compare mailservers
	if !stringSlicesEqual(old.Mailservers, new.Mailservers) {
		eq = false
		message := fmt.Sprintf("The MX records of domain %s are changed.", new.Domainname)
		log.Println(message)
		err := alerting.Discord(alert.DiscordWebhookURL, alert.DiscordUsername, message)
		if err != nil {
			log.Println(err.Error())
		}
	}

	// Compare records
	// TODO: Needs some work
	if (len(old.Records) > 0) || (len(new.Records) > 0) {
		ref := reflect.DeepEqual(old.Records, new.Records)
		if !ref {
			eq = false
			message := fmt.Sprintf("Some records of domain %s are changed.", new.Domainname)
			log.Println(message)
			err := alerting.Discord(alert.DiscordWebhookURL, alert.DiscordUsername, message)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

	// done
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
