package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NotificationClient struct {
	baseURL string
	client  *http.Client
}

func NewNotificationClient(baseURL string) *NotificationClient {
	return &NotificationClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *NotificationClient) SendNotification(userID int, subject, body string) error {
	payload := map[string]interface{}{
		"user_id": userID,
		"type":    "email",
		"subject": subject,
		"body":    body,
	}
	jsonData, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/api/notify", c.baseURL)
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("notification service returned status %d", resp.StatusCode)
	}
	return nil
}
