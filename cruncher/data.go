package cruncher

// Domain struct to store information. Need omitempty for empty results.
type Domain struct {
	Domainname  string   `json:"domainname,omitempty"`
	Serial      uint32   `json:"serial,omitempty"`
	Nameservers []string `json:"nameservers,omitempty"`
	Mailservers []string `json:"mailservers,omitempty"`
}
