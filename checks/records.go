package checks

import "github.com/miekg/dns"

// GetCNAME function
func GetCNAME(hostname string, nameserver string) (string, error) {
	var cname string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(hostname), dns.TypeCNAME)
	c := new(dns.Client)
	m.MsgHdr.RecursionDesired = true
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		return "none", err
	}
	for _, rin := range in.Answer {
		if r, ok := rin.(*dns.CNAME); ok {
			cname = r.Target
		}
	}
	return cname, nil
}

// GetA function
func GetA(hostname string, nameserver string) ([]string, error) {
	var record []string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(hostname), dns.TypeA)
	c := new(dns.Client)
	m.MsgHdr.RecursionDesired = true
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		return nil, err
	}
	for _, rin := range in.Answer {
		if r, ok := rin.(*dns.A); ok {
			record = append(record, r.A.String())
		}
	}

	return record, nil
}

// GetAAAA function
func GetAAAA(hostname string, nameserver string) ([]string, error) {
	var record []string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(hostname), dns.TypeAAAA)
	c := new(dns.Client)
	m.MsgHdr.RecursionDesired = true
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		return nil, err
	}
	for _, rin := range in.Answer {
		if r, ok := rin.(*dns.AAAA); ok {
			record = append(record, r.AAAA.String())
		}
	}

	return record, nil
}
