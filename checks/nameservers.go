package checks

import (
	"sort"
	"strings"

	"github.com/miekg/dns"
)

// GetNameservers function to get NS records from DNS
func GetNameservers(domain string, nameserver string) ([]string, error) {
	var answer []string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeNS)
	m.MsgHdr.RecursionDesired = true
	m.SetEdns0(4096, true)
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		return answer, err
	}
	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.NS); ok {
			answer = append(answer, strings.ToLower(a.Ns))
		}
	}
	// Need to sort data to be able to compare
	sort.Strings(answer)
	return answer, nil
}
