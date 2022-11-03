package checks

import (
	"sort"

	"github.com/miekg/dns"
)

// GetMailserver function to get MX records from DNS
func GetMailservers(domain string, nameserver string) ([]string, error) {
	var answer []string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeMX)
	m.MsgHdr.RecursionDesired = true
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		return answer, err
	}
	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.MX); ok {
			answer = append(answer, a.Mx)
		}
	}
	// Need to sort data to be able to compare
	sort.Strings(answer)
	return answer, nil
}
