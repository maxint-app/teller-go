package teller

import (
	"encoding/json"
	"net/http"
)

// TellerAddress represents an address associated with an identity
type TellerAddress struct {
	Primary     bool   `json:"primary"`
	Street      string `json:"street"`
	City        string `json:"city"`
	Region      string `json:"region"`
	PostalCode  string `json:"postal_code"`
	CountryCode string `json:"country_code"`
}

// TellerOwner represents the owner of an account
type TellerOwner struct {
	Type  string `json:"type"` // "person" or "business"
	Names []struct {
		Type string `json:"type"`
		Data string `json:"data"`
	} `json:"names"`
	Addresses    []TellerAddress `json:"addresses"`
	PhoneNumbers []struct {
		Type string `json:"type"` // "mobile", "home", "work", "unknown"
		Data string `json:"data"`
	} `json:"phone_numbers"`
	Emails []struct {
		Data string `json:"data"`
	} `json:"emails"`
}

// TellerIdentity represents identity information for an account
type TellerIdentity struct {
	Account TellerAccount `json:"account"`
	Owners  []TellerOwner `json:"owners"`
}

// IdentityModule handles identity-related API calls
type IdentityModule struct {
	client *Client
}

// Get retrieves identity information
func (m *IdentityModule) Get(options *TellerOptionsBase) ([]TellerIdentity, error) {
	req, err := http.NewRequest("GET", m.client.baseURL+"/identity", nil)
	if err != nil {
		return nil, err
	}

	if options != nil && options.AccessToken != "" {
		req.SetBasicAuth(options.AccessToken, "")
	} else if m.client.accessToken != "" {
		req.SetBasicAuth(m.client.accessToken, "")
	}

	resp, err := m.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []TellerIdentity
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
