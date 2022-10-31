package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/miekg/dns"
)

// Load configuration file and parse the json content
func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return config, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, err
}

// Config struct for use in the applications with some general values
type Config struct {
	Contact    string         `json:"contact"`
	Nameserver string         `json:"nameserver"`
	Domains    []ConfigDomain `json:"domains"`
}

// ConfigDomain struct for domains to monitor.
type ConfigDomain struct {
	Name string `json:"name"`
}

func main() {
	// Define the flag for the application for opening a config file.
	configFile := flag.String("config", "config.json", "What config file to use. (Required)")
	flag.Parse()
	if *configFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read config or create an error
	conf, err := LoadConfiguration(*configFile)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// Check if there are more then 0 domains in the config file
	if len(conf.Domains) < 1 {
		println("No domains found in configfile")
		os.Exit(0)
	}

	// Loop domains from config file
	for _, d := range conf.Domains {
		fmt.Println(d.Name)
		soa, err := getSOA(d.Name, conf.Nameserver)
		if err != nil {
			println(err)
		}
		println(soa.Serial)
	}

}

// SOA struct for SOA information aquired from the nameserver.
type SOA struct {
	Ns      string `json:"ns,omitempty"`
	Mbox    string `json:"mbox,omitempty"`
	Serial  uint32 `json:"serial,omitempty"`
	Refresh uint32 `json:"refresh,omitempty"`
	Retry   uint32 `json:"retry,omitempty"`
	Expire  uint32 `json:"expire,omitempty"`
	Minttl  uint32 `json:"minttl,omitempty"`
}

// getSOA function to resolve SOA record from domain
func getSOA(domain string, nameserver string) (*SOA, error) {
	answer := new(SOA)
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
			answer.Serial = soa.Serial   // uint32
			answer.Ns = soa.Ns           // string
			answer.Expire = soa.Expire   // uint32
			answer.Mbox = soa.Mbox       // string
			answer.Minttl = soa.Minttl   // uint32
			answer.Refresh = soa.Refresh // uint32
			answer.Retry = soa.Retry     // uint32
		}
	}
	return answer, nil
}
