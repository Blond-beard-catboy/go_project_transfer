package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_project_transfer/route-service/internal/models"
)

type CargoClient struct {
	baseURL string
	client  *http.Client
}

func NewCargoClient(baseURL string) *CargoClient {
	return &CargoClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *CargoClient) GetCargo(cargoID int) (*models.CargoData, error) {
	url := fmt.Sprintf("%s/api/cargo/%d", c.baseURL, cargoID)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cargo service returned status %d", resp.StatusCode)
	}

	var cargo models.CargoData
	if err := json.NewDecoder(resp.Body).Decode(&cargo); err != nil {
		return nil, err
	}
	return &cargo, nil
}
