package alerting

import (
	"dnsmon/config"
	"log"
	"strings"
)

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
