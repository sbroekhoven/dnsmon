package cruncher

import (
	log "github.com/sirupsen/logrus"
)

func Compare(old Domain, new Domain) (bool, error) {
	eq := true
	compareLogger := log.WithFields(log.Fields{
		"domain":    old.Domainname,
		"function:": "compare",
	})
	compareLogger.Info("start comparing")

	// Compare domain zone serial
	if old.Serial != new.Serial {
		eq = false
		compareLogger.Warn("serial changed")
	}

	// Compare nameservers
	if !stringSlicesEqual(old.Nameservers, new.Nameservers) {
		eq = false
		compareLogger.Warn("nameservers changed")
	}

	// Compare mailservers
	if !stringSlicesEqual(old.Mailservers, new.Mailservers) {
		eq = false
		compareLogger.Warn("mailservers changed")
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
