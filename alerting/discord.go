package alerting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Discord function to send out alerts to a Discord channel webhook
func Discord(url string, username string, content string) error {
	payload := new(DiscordMessage)
	payload.Username = username
	payload.Content = content

	jsonValue, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf(string(responseBody))
	}

	return nil
}

type DiscordMessage struct {
	Username string `json:"username,omitempty"`
	Content  string `json:"content,omitempty"`
}
