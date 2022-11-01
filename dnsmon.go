package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"

	"dnsmon/checks"
	"dnsmon/config"
	"dnsmon/cruncher"
)

// This is the main function
func main() {
	// Define the flag for the application for opening a config file.
	configFile := flag.String("config", "config.json", "What config file to use. (Required)")
	flag.Parse()
	if *configFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read config or create an error
	conf, err := config.LoadConfiguration(*configFile)
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
		// Open stored domain data from json file
		storedData, err := cruncher.ReadJSON(d.Name + ".last.json")
		if err != nil {
			println(err.Error())
		}
		jsonStored, err := json.MarshalIndent(storedData, "", "   ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", jsonStored)

		data := new(cruncher.Domain)
		data.Domainname = d.Name

		// Get the serial number of the zone file
		domainSerial, err := checks.GetSerial(d.Name, conf.Nameserver)
		if err != nil {
			println(err.Error())
		}
		data.Serial = domainSerial

		// Get the nameservers for the domains
		domainNameservers, err := checks.GetNameservers(d.Name, conf.Nameserver)
		if err != nil {
			println(err.Error())
		}
		data.Nameservers = domainNameservers

		// Get the mailservers for the domain
		domainMailservers, err := checks.GetMailservers(d.Name, conf.Nameserver)
		if err != nil {
			println(err.Error())
		}
		data.Mailservers = domainMailservers

		json, err := json.MarshalIndent(data, "", "   ")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", json)

		// Store the json to file
		// prepare filename
		filename := data.Domainname + ".last.json"

		f, err := os.OpenFile(filename, os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		n, err := f.Write(json)
		if err != nil {
			panic(err)
		}
		fmt.Printf("wrote %d bytes\n", n)

		// compare the stuff
		var resolvedData cruncher.Domain = *data
		if reflect.DeepEqual(storedData, resolvedData) {
			fmt.Println("storedData is equal to resolvedData")
		} else {
			fmt.Println("storedData is not equal to resolvedData")
		}
	}

}
