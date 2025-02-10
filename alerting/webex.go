package alerting

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func Webex(webexRoomID string, webexToken string, content string) error {

	// Create an HTTP client with or without proxy
	client := &http.Client{}

	// Check if the required configuration is present
	if webexToken == "" || webexRoomID == "" {
		return fmt.Errorf("webex token or room ID is missing in the configuration")
	}

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the roomId field
	_ = writer.WriteField("roomId", webexRoomID)

	// Add the markdown field
	_ = writer.WriteField("markdown", content)

	// Close the multipart writer
	err := writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://webexapis.com/v1/messages", body)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+webexToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(responseBody))
	}

	fmt.Println("Message sent successfully")
	return nil
}
