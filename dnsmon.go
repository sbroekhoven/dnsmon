package main

import (
	"flag"
	"log"
	"os"
	"time"

	"dnsmon/config"
	"dnsmon/cruncher"
)

// This is the main function of dnsmon
func main() {
	// Setup some logging defaults
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Define the flag for the application for opening a config file.
	configFile := flag.String("config", "config.yaml", "What config file to use. (Required)")
	flag.Parse()
	if *configFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read config or create an error and quit.
	conf, err := config.LoadConfiguration(*configFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Check if there are more then 0 domains in the config file
	if len(conf.Domains) < 1 {
		log.Println("No domains found in configfile")
		os.Exit(0)
	}

	// Loop domains from config file
	for _, d := range conf.Domains {
		// Set some vars.
		// firstRun: is for checking if thare are any previous output files. Otherwise there is nothong to compare.
		// eq: is to check if old and new information is the same
		firstRun := false
		var eq bool = false

		// Define filenames for reading and storing.
		filenameLast := conf.Output + d.Name + ".current.json"
		filenameArch := conf.Output + d.Name + "." + time.Now().Format("20060102150405") + ".json"

		// Open stored domain data from json file
		storedData, err := cruncher.ReadJSON(filenameLast)
		if err != nil {
			// Probably first run and files do not exist, so set firstRun to true
			log.Println(err.Error())
			firstRun = true
		}

		// This function collects the information for the DNS.
		data, err := cruncher.Collect(d, conf.Resolver1)
		if err != nil {
			// If there is an error, print it and start using the second resolver to doublecheck.
			log.Println(err.Error())
			log.Println("start trying second resolver now")
			data, err = cruncher.Collect(d, conf.Resolver2)
			if err != nil {
				// If there is still an error in looking this up, skip the rest and continue with next domain.
				log.Println(err.Error())
				continue
			}
		}

		// create this link to pointed data
		var collectedData cruncher.Domain = *data

		// If this is a firstRun, store new current information in output file and then continue
		if firstRun {
			written, err := cruncher.WriteJSON(filenameLast, collectedData)
			if err != nil {
				log.Fatalln(err.Error())
			}
			log.Printf("file written: %s with %d bytes", filenameLast, written)
			log.Println("first run for this domain, continue to next domain")
			continue
		}

		// if this is not a firstRun, compare the old and new information
		if !firstRun {
			eq, err = cruncher.Compare(conf.Alerting, storedData, collectedData)
			if err != nil {
				log.Println(err.Error())
			}
			// log.Println(eq)
			// If not equal, store the changes in a json file
			if !eq {
				// Move current json to an "archive" json
				err = os.Rename(filenameLast, filenameArch)
				if err != nil {
					log.Fatalln(err.Error())
				}

				// Write new current JSON file with new collected data
				written, err := cruncher.WriteJSON(filenameLast, collectedData)
				if err != nil {
					log.Fatalln(err.Error())
				}
				log.Printf("file written: %s with %d bytes", filenameLast, written)
			}
		}
	}
}
