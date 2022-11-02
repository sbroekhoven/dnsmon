package cruncher

import (
	"dnsmon/checks"
	"dnsmon/config"
)

func Collect(domain config.ConfigDomain, nameserver string) (*Domain, error) {
	// Collect domain information
	data := new(Domain)
	data.Domainname = domain.Name

	// Get the serial number of the zone file
	domainSerial, err := checks.GetSerial(domain.Name, nameserver)
	if err != nil {
		return data, err
	}
	data.Serial = domainSerial

	// Get the nameservers for the domains
	domainNameservers, err := checks.GetNameservers(domain.Name, nameserver)
	if err != nil {
		return data, err
	}
	if len(domainNameservers) > 0 {
		data.Nameservers = domainNameservers
	}

	// Get the mailservers for the domain
	domainMailservers, err := checks.GetMailservers(domain.Name, nameserver)
	if err != nil {
		return data, err
	}
	if len(domainMailservers) > 0 {
		data.Mailservers = domainMailservers
	}

	return data, nil
}
