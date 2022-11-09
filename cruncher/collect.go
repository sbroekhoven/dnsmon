package cruncher

import (
	"dnsmon/checks"
	"dnsmon/config"
	"log"
)

func Collect(domain config.Domain, nameserver string) (*Domain, error) {
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
	if len(domainNameservers) != 0 {
		data.Nameservers = domainNameservers
	}

	// Get the mailservers for the domain
	domainMailservers, err := checks.GetMailservers(domain.Name, nameserver)
	if err != nil {
		return data, err
	}
	if len(domainMailservers) != 0 {
		data.Mailservers = domainMailservers
	}

	for _, r := range domain.Records {
		log.Println(r)
		record, err := GetRecord(r, nameserver)
		data.Records = append(data.Records, *record)
		if err != nil {
			// return data, err
			log.Println(err.Error())
		}
	}

	return data, nil
}

// getHosts function
func GetRecord(record string, nameserver string) (*Record, error) {
	r := new(Record)

	r.Hostname = record

	cname, err := checks.GetCNAME(r.Hostname, nameserver)
	if err != nil {
		return r, err
	}

	if len(cname) > 0 {
		r.CNAME = cname
		return r, nil
	}

	ar, err := checks.GetA(r.Hostname, nameserver)
	if err != nil {
		return r, err
	}
	r.IPv4 = ar

	aaaar, err := checks.GetAAAA(r.Hostname, nameserver)
	if err != nil {
		return r, err
	}
	r.IPv6 = aaaar

	return r, nil

}
