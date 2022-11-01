package checks

import "github.com/miekg/dns"

// getSerial function to resolve SOA record from domain
func GetSerial(domain string, nameserver string) (uint32, error) {
	var answer uint32
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeSOA)
	c := new(dns.Client)
	m.MsgHdr.RecursionDesired = true
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		return answer, err
	}
	for _, ain := range in.Answer {
		if soa, ok := ain.(*dns.SOA); ok {
			answer = soa.Serial
		}
	}
	return answer, nil
}
