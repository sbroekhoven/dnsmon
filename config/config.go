package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Load configuration file and parse the json content
func LoadConfiguration(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// Config struct for use in the applications with some general values
type Config struct {
	Contact   string         `yaml:"contact,omitempty"`
	Resolver1 string         `yaml:"resolver1,omitempty"`
	Resolver2 string         `yaml:"resolver2,omitempty"`
	Alerting  ConfigAlerting `yaml:"alerting,omitempty"`
	Domains   []ConfigDomain `yaml:"domains,omitempty"`
	Output    string         `yaml:"output,omitempty"`
}

// ConfigDomain struct for domains to monitor.
type ConfigDomain struct {
	Name string `yaml:"name,omitempty"`
}

// ConfigAlerting struct for domains to monitor.
type ConfigAlerting struct {
	DiscordUsername   string `yaml:"discord_username,omitempty"`
	DiscordWebhookURL string `yaml:"discord_webhook_url,omitempty"`
}
