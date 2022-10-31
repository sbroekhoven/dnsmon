package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

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

		// Get the serial number of the zone file
		domainSerial, err := getSerial(d.Name, conf.Nameserver)
		if err != nil {
			println(err.Error())
		}
		println(domainSerial)

		// Get the nameservers for the domains
		domainNameservers, err := getNameservers(d.Name, conf.Nameserver)
		if err != nil {
			println(err.Error())
		}
		fmt.Printf("%v\n", domainNameservers)

		// Get the mailservers for the domain
		domainMailservers, err := getMailservers(d.Name, conf.Nameserver)
		if err != nil {
			println(err.Error())
		}
		fmt.Printf("%v\n", domainMailservers)

	}

}

// getSerial function to resolve SOA record from domain
func getSerial(domain string, nameserver string) (uint32, error) {
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

func getNameservers(domain string, nameserver string) ([]string, error) {
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
	return answer, nil
}

func getMailservers(domain string, nameserver string) ([]string, error) {
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
	return answer, nil
}
