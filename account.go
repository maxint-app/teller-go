package teller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TellerAccountType = string

const (
	TellerAccountTypeDepository TellerAccountType = "depository"
	TellerAccountTypeCredit     TellerAccountType = "credit"
)

type TellerAccountSubtype = string

const (
	TellerAccountSubtypeChecking             TellerAccountSubtype = "checking"
	TellerAccountSubtypeSavings              TellerAccountSubtype = "savings"
	TellerAccountSubtypeMoneyMarket          TellerAccountSubtype = "money_market"
	TellerAccountSubtypeCertificateOfDeposit TellerAccountSubtype = "certificate_of_deposit"
	TellerAccountSubtypeTreasury             TellerAccountSubtype = "treasury"
	TellerAccountSubtypeCreditCard           TellerAccountSubtype = "credit_card"
	TellerAccountSubtypeSweep                TellerAccountSubtype = "sweep"
)

type TellerAccountStatusType = string

const (
	TellerAccountStatusTypeOpen   TellerAccountStatusType = "open"
	TellerAccountStatusTypeClosed TellerAccountStatusType = "closed"
)

// TellerAccount represents a bank account
type TellerAccount struct {
	Currency     string            `json:"currency"`
	EnrollmentID string            `json:"enrollment_id"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         TellerAccountType `json:"type"` // "depository" or "credit"
	Institution  struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"institution"`
	LastFour string `json:"last_four"`
	Links    struct {
		Self         string `json:"self"`
		Details      string `json:"details"`
		Balances     string `json:"balances"`
		Transactions string `json:"transactions"`
	} `json:"links"`
	Subtype TellerAccountSubtype    `json:"subtype"` // "checking", "savings", etc.
	Status  TellerAccountStatusType `json:"status"`  // "open" or "closed"
}

// TellerAccountDetails represents detailed account information
type TellerAccountDetails struct {
	AccountID     string `json:"account_id"`
	AccountNumber string `json:"account_number"`
	Links         struct {
		Self    string `json:"self"`
		Account string `json:"account"`
	} `json:"links"`
	RoutingNumbers struct {
		ACH  *string `json:"ach"`
		Wire *string `json:"wire"`
		BACS *string `json:"bacs"`
	} `json:"routing_numbers"`
}

// TellerAccountBalances represents account balance information
type TellerAccountBalances struct {
	AccountID string `json:"account_id"`
	Ledger    string `json:"ledger"`
	Available string `json:"available"`
	Links     struct {
		Self    string `json:"self"`
		Account string `json:"account"`
	} `json:"links"`
}

// AccountModule handles account-related API calls
type AccountModule struct {
	client *Client
}

// List retrieves all accounts
func (m *AccountModule) List(options *TellerOptionsBase) ([]TellerAccount, error) {
	req, err := http.NewRequest("GET", m.client.baseURL+"/account", nil)
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

	var result []TellerAccount
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// Get retrieves a single account by ID
func (m *AccountModule) Get(id string, options *TellerOptionsBase) (*TellerAccount, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s", m.client.baseURL, id), nil)
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

	var result TellerAccount
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Remove deletes a single account by ID
func (m *AccountModule) Remove(id string, options *TellerOptionsBase) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/accounts/%s", m.client.baseURL, id), nil)
	if err != nil {
		return err
	}

	if options != nil && options.AccessToken != "" {
		req.SetBasicAuth(options.AccessToken, "")
	} else if m.client.accessToken != "" {
		req.SetBasicAuth(m.client.accessToken, "")
	}

	resp, err := m.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// RemoveAll deletes all accounts
func (m *AccountModule) RemoveAll(options *TellerOptionsBase) error {
	req, err := http.NewRequest("DELETE", m.client.baseURL+"/accounts", nil)
	if err != nil {
		return err
	}

	if options != nil && options.AccessToken != "" {
		req.SetBasicAuth(options.AccessToken, "")
	} else if m.client.accessToken != "" {
		req.SetBasicAuth(m.client.accessToken, "")
	}

	resp, err := m.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Details retrieves detailed information for an account
func (m *AccountModule) Details(id string, options *TellerOptionsBase) (*TellerAccountDetails, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s/details", m.client.baseURL, id), nil)
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

	var result TellerAccountDetails
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Balances retrieves balance information for an account
func (m *AccountModule) Balances(id string, options *TellerOptionsBase) (*TellerAccountBalances, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s/balances", m.client.baseURL, id), nil)
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

	var result TellerAccountBalances
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
