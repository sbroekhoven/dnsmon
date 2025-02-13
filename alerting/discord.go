package alerting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Discord sends a message to a Discord webhook URL using the Discord API.
//
// This function returns an error if there is an error sending the request or
// reading the response. If the response status code is not 200 or 204, the
// response body is read and returned as an error.
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

		return fmt.Errorf("error response from discord: %s", string(responseBody))

	}

	return nil
}

type DiscordMessage struct {
	Username string `json:"username,omitempty"`
	Content  string `json:"content,omitempty"`
}
