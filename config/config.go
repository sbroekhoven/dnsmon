package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Load configuration file and parse the YAML content
// LoadConfiguration loads the configuration from the specified YAML file.
// It takes the path to the configuration file as an argument and returns a pointer
// to the Config struct and an error if any occurs during the process.
//
// The function performs the following steps:
// 1. Creates an empty Config struct.
// 2. Opens the specified configuration file.
// 3. Initializes a new YAML decoder.
// 4. Decodes the YAML content into the Config struct.
//
// If the file cannot be opened or the YAML content cannot be decoded, an error is returned.
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
	Enabled           []string `yaml:"enabled,omitempty"`
	DiscordUsername   string   `yaml:"discord_username,omitempty"`
	DiscordWebhookURL string   `yaml:"discord_webhook_url,omitempty"`
	WebexRoom         string   `yaml:"webex_room,omitempty"`
	WebexToken        string   `yaml:"webex_token,omitempty"`
}
