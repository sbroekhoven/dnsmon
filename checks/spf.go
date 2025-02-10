package checks

import (
	"strings"

	"github.com/miekg/dns"
)

// GetSPFRecord function to get FPF records from DNS
func GetSPFRecord(domain string, nameserver string) (string, error) {
	var answer string
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)
	m.MsgHdr.RecursionDesired = true
	c := new(dns.Client)
	in, _, err := c.Exchange(m, nameserver+":53")
	if err != nil {
		return answer, err
	}

	for _, ain := range in.Answer {
		if a, ok := ain.(*dns.TXT); ok {
			txtecord := strings.Join(a.Txt, "")
			lower := strings.ToLower(txtecord)
			// println(lower)
			if strings.HasPrefix(lower, "v=spf1") {
				return lower, nil
			}
		}
	}

	return answer, nil
}
