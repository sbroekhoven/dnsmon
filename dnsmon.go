package main

import (
	"bufio"
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

	// Check if domains_file exists
	if conf.DomainsFile != "" {
		_, err := os.Stat(conf.DomainsFile)
		if os.IsNotExist(err) {
			log.Fatalf("Error: The domains file '%s' specified in config does not exist.\n", conf.DomainsFile)
		} else if err != nil {
			log.Fatalf("Error checking domains file: %v\n", err)
		}
	} else {
		log.Fatalln("Error: No domains file specified in config.")
	}

	// Open the domains file
	file, err := os.Open(conf.DomainsFile)
	log.Printf("Processing domains file: %s\n", conf.DomainsFile)

	if err != nil {
		log.Fatalf("Error opening domains file: %v\n", err)
	}
	defer file.Close()

	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := scanner.Text()
		// firstRun: is for checking if thare are any previous output files. Otherwise there is nothong to compare.
		// eq: is to check if old and new information is the same
		firstRun := false
		var eq bool = false

		// Define filenames for reading and storing.
		filenameLast := conf.Output + domain + ".current.json"
		filenameArch := conf.Output + domain + "." + time.Now().Format("20060102150405") + ".json"

		// Open stored domain data from json file
		storedData, err := cruncher.ReadJSON(filenameLast)
		if err != nil {
			// Probably first run and files do not exist, so set firstRun to true
			log.Println(err.Error())
			firstRun = true
		}

		// This function collects the information for the DNS.
		data, err := cruncher.Collect(domain, conf.Resolver1)
		if err != nil {
			// If there is an error, print it and start using the second resolver to doublecheck.
			log.Printf("Error with primary resolver: %v. Trying secondary resolver.", err)

			// Check if a secondary resolver is configured
			if conf.Resolver2 != "" {
				log.Println("Trying secondary resolver.")
				// Try again with the second resolver
				data, err = cruncher.Collect(domain, conf.Resolver2)
				if err != nil {
					// If there's still an error with the second resolver, log it and skip this domain
					log.Printf("Error with secondary resolver for domain %s: %v. Skipping this domain.", domain, err)
					continue
				}
			} else {
				// If no secondary resolver is configured, log it and skip this domain
				log.Printf("No secondary resolver configured. Skipping domain %s.", domain)
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

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading domains file: %v\n", err)
	}
}
