package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RouteClient struct {
	baseURL string
	client  *http.Client
}

func NewRouteClient(baseURL string) *RouteClient {
	return &RouteClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *RouteClient) CreateRoute() (int, error) {
	url := fmt.Sprintf("%s/api/routes", r.baseURL)
	resp, err := r.client.Post(url, "application/json", nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("route service returned status %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	routeID, ok := result["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid response from route service")
	}
	return int(routeID), nil
}

func (r *RouteClient) AddCargoToRoute(routeID, cargoID int) error {
	url := fmt.Sprintf("%s/api/routes/%d/cargo/%d", r.baseURL, routeID, cargoID)
	resp, err := r.client.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("route service returned status %d", resp.StatusCode)
	}
	return nil
}
