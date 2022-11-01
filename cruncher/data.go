package cruncher

type Domain struct {
	Domainname  string   `json:"domainname"`
	Serial      uint32   `json:"serial"`
	Nameservers []string `json:"nameservers"`
	Mailservers []string `json:"mailservers"`
}
