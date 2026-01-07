package teller

import (
	"encoding/json"
	"net/http"
)

// TellerInstitution represents a financial institution
type TellerInstitution struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Products []string `json:"products"`
}

// InstitutionsModule handles institution-related API calls
type InstitutionsModule struct {
	client *Client
}

// List retrieves all institutions
func (m *InstitutionsModule) List() ([]TellerInstitution, error) {
	req, err := http.NewRequest("GET", m.client.baseURL+"/institutions", nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []TellerInstitution
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
