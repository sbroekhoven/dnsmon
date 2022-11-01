package main

import (
	"flag"
	"os"
	"reflect"

	"dnsmon/checks"
	"dnsmon/config"
	"dnsmon/cruncher"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Set some logging defaults
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// This is the main function
func main() {
	log.Info("Start")

	// Define the flag for the application for opening a config file.
	configFile := flag.String("config", "config.json", "What config file to use. (Required)")
	flag.Parse()
	if *configFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Infof("Reading config file %v", *configFile)

	// Read config or create an error
	conf, err := config.LoadConfiguration(*configFile)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// Check if there are more then 0 domains in the config file
	if len(conf.Domains) < 1 {
		println("No domains found in configfile")
		log.Info("No domains in config file found")
		os.Exit(0)
	}

	// Loop domains from config file
	for _, d := range conf.Domains {
		// A common pattern is to re-use fields between logging statements by re-using
		// the logrus.Entry returned from WithFields()
		domainLogger := log.WithFields(log.Fields{
			"domain": d.Name,
		})

		// Open stored domain data from json file
		domainLogger.Info("Open last stored JSON file")
		storedData, err := cruncher.ReadJSON(d.Name + ".last.json")
		if err != nil {
			domainLogger.Error(err.Error())
		}

		data := new(cruncher.Domain)
		data.Domainname = d.Name

		// Get the serial number of the zone file
		domainLogger.Info("Get serial")
		domainSerial, err := checks.GetSerial(d.Name, conf.Nameserver)
		if err != nil {
			domainLogger.Error(err.Error())
		}
		data.Serial = domainSerial

		// Get the nameservers for the domains
		domainLogger.Info("Get nameservers")
		domainNameservers, err := checks.GetNameservers(d.Name, conf.Nameserver)
		if err != nil {
			domainLogger.Error(err.Error())
		}
		data.Nameservers = domainNameservers

		// Get the mailservers for the domain
		domainLogger.Info("Get mailservers")
		domainMailservers, err := checks.GetMailservers(d.Name, conf.Nameserver)
		if err != nil {
			domainLogger.Error(err.Error())
		}
		data.Mailservers = domainMailservers

		// Store the json to file
		// prepare filename
		domainLogger.Info("Write JSON file")
		filename := data.Domainname + ".last.json"
		written, err := cruncher.WriteJSON(filename, *data)
		if err != nil {
			domainLogger.Error(err.Error())
		}
		domainLogger.Infof("Written %d bytes", written)

		// compare the stuff
		var resolvedData cruncher.Domain = *data
		if reflect.DeepEqual(storedData, resolvedData) {
			domainLogger.Info("storedData is equal to resolvedData")
		} else {
			domainLogger.Info("storedData is not equal to resolvedData")
		}
	}
}
