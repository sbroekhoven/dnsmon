package cruncher

import (
	"dnsmon/checks"
)

func Collect(domain string, nameserver string) (*Domain, error) {
	// Collect domain information
	data := new(Domain)
	data.Domainname = domain

	// Get the serial number of the zone file
	domainSerial, err := checks.GetSerial(domain, nameserver)
	if err != nil {
		return data, err
	}
	data.Serial = domainSerial

	// Get the nameservers for the domains
	domainNameservers, err := checks.GetNameservers(domain, nameserver)
	if err != nil {
		return data, err
	}
	if len(domainNameservers) != 0 {
		data.Nameservers = domainNameservers
	}

	// Get the mailservers for the domain
	domainMailservers, err := checks.GetMailservers(domain, nameserver)
	if err != nil {
		return data, err
	}
	if len(domainMailservers) != 0 {
		data.Mailservers = domainMailservers
	}

	// Get SPF record for the domain
	spfRecord, err := checks.GetSPFRecord(domain, nameserver)
	if err != nil {
		return data, err
	}
	data.SPFRecord = spfRecord

	// Get DMARC record for the domain
	dmarcRecord, err := checks.GetDMARCRecord(domain, nameserver)
	if err != nil {
		return data, err
	}
	data.DMARCRecord = dmarcRecord

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
