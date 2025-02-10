package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Load configuration file and parse the YAML content
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
	Contact     string   `yaml:"contact,omitempty"`
	Resolver1   string   `yaml:"resolver1,omitempty"`
	Resolver2   string   `yaml:"resolver2,omitempty"`
	Alerting    Alerting `yaml:"alerting,omitempty"`
	DomainsFile string   `yaml:"domains_file,omitempty"`
	Output      string   `yaml:"output,omitempty"`
}

// Alerting struct for alerting purpose.
type Alerting struct {
	DiscordUsername   string `yaml:"discord_username,omitempty"`
	DiscordWebhookURL string `yaml:"discord_webhook_url,omitempty"`
}
