package main

import (
	"flag"
	"log"
	"os"
	"time"

	"dnsmon/config"
	"dnsmon/cruncher"
)

// This is the main function
func main() {
	// Setup Logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Define the flag for the application for opening a config file.
	configFile := flag.String("config", "config.yaml", "What config file to use. (Required)")
	flag.Parse()
	if *configFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read config or create an error
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
		firstRun := false
		var eq bool = false

		// Define filenames
		filenameLast := conf.Output + d.Name + ".current.json"
		filenameArch := conf.Output + d.Name + "." + time.Now().Format("20060102150405") + ".json"

		// Open stored domain data from json file
		storedData, err := cruncher.ReadJSON(filenameLast)
		if err != nil {
			// Probably first run and files do not exist
			log.Println(err.Error())
			firstRun = true
		}

		data, err := cruncher.Collect(d, conf.Resolver1)
		if err != nil {
			// Error collecting information regarding this domain
			// Contintue with the next domain
			log.Println(err.Error())
			log.Println("start trying second resolver now")
			data, err = cruncher.Collect(d, conf.Resolver1)
			if err != nil {
				log.Println(err.Error())
				continue
			}
		}
		var collectedData cruncher.Domain = *data

		// Store new current information
		written, err := cruncher.WriteJSON(filenameLast, collectedData)
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("file written: %s with %d bytes", filenameLast, written)

		// Also store a file for archiving on date
		writtenArch, err := cruncher.WriteJSON(filenameArch, collectedData)
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("file written: %s with %d bytes", filenameArch, writtenArch)

		if !firstRun {
			// Own compare functionallity
			eq, err = cruncher.Compare(conf.Alerting, storedData, collectedData)
			if err != nil {
				log.Println(err.Error())
			}
			log.Println(eq)
		}
	}
}
