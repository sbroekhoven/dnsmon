package cruncher

// Domain struct to store information. Need omitempty for empty results.
type Domain struct {
	Domainname  string   `json:"domainname,omitempty"`
	Serial      uint32   `json:"serial,omitempty"`
	Nameservers []string `json:"nameservers,omitempty"`
	Mailservers []string `json:"mailservers,omitempty"`
	Records     []Record `json:"records,omitempty"`
}

type Record struct {
	Hostname string   `json:"hostname,omitempty"`
	IPv4     []string `json:"ipv4,omitempty"`
	IPv6     []string `json:"ipv6,omitempty"`
	CNAME    string   `json:"cname,omitempty"`
}
