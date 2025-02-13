package alerting

import (
	"dnsmon/config"
	"log"
	"strings"
)

// SendAlerts sends an alert to all enabled alerting methods specified in the config.
// It is the caller's responsibility to make sure that the alerting method is enabled
// in the configuration before calling this function.
//
// If the configuration is missing the required information for a particular
// alerting method (e.g. a webhook URL or username), the function will log an
// error but will not return an error.
//
// The message is sent as-is to each alerting method, without any modification.
func SendAlerts(alert config.Alerting, message string) {
	for _, enabled := range alert.Enabled {
		switch strings.ToLower(enabled) {
		case "discord":
			if alert.DiscordWebhookURL != "" && alert.DiscordUsername != "" {
				err := Discord(alert.DiscordWebhookURL, alert.DiscordUsername, message)
				if err != nil {
					log.Printf("Error sending Discord alert: %v\n", err)
				}
			} else {
				log.Println("Discord alerting is enabled but webhook URL or username is missing")
			}
		case "webex":
			if alert.WebexRoom != "" && alert.WebexToken != "" {
				err := Webex(alert.WebexRoom, alert.WebexToken, message)
				if err != nil {
					log.Printf("Error sending Webex alert: %v\n", err)
				}
			} else {
				log.Println("Webex alerting is enabled but room or token is missing")
			}
		default:
			log.Printf("Unknown alerting method: %s\n", enabled)
		}
	}
}
