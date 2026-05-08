package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PaymentClient struct {
	baseURL string
	client  *http.Client
}

func NewPaymentClient(baseURL string) *PaymentClient {
	return &PaymentClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *PaymentClient) CreatePayment(orderID int, amount float64) error {
	payload := map[string]interface{}{
		"order_id": orderID,
		"amount":   amount,
	}
	jsonData, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/api/payments", c.baseURL)
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("payment service returned status %d", resp.StatusCode)
	}
	return nil
}
