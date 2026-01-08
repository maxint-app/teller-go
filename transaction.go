package teller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TellerTransactionProcessingType = string

const (
	TellerTransactionProcessingTypePending  TellerTransactionProcessingType = "pending"
	TellerTransactionProcessingTypeComplete TellerTransactionProcessingType = "complete"
)

type TellerTransactionCounterPartyType = string

const (
	TellerTransactionCounterPartyTypePerson       TellerTransactionCounterPartyType = "person"
	TellerTransactionCounterPartyTypeOrganization TellerTransactionCounterPartyType = "organization"
)

type TellerTransactionStatusType = string

const (
	TellerTransactionStatusTypePosted  TellerTransactionStatusType = "posted"
	TellerTransactionStatusTypePending TellerTransactionStatusType = "pending"
)

// TellerTransaction represents a financial transaction
type TellerTransaction struct {
	AccountID   string `json:"account_id"`
	Amount      string `json:"amount"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Details     struct {
		ProcessingStatus TellerTransactionProcessingType `json:"processing_status"` // "pending" or "complete"
		Category         string                          `json:"category"`
		Counterparty     struct {
			Name *string                           `json:"name"`
			Type TellerTransactionCounterPartyType `json:"type"` // "person" or "organization"
		} `json:"counterparty"`
	} `json:"details"`
	Status TellerTransactionStatusType `json:"status"` // "posted" or "pending"
	ID     string                      `json:"id"`
	Links  struct {
		Self    string `json:"self"`
		Account string `json:"account"`
	} `json:"links"`
	RunningBalance *string `json:"running_balance"`
	Type           string  `json:"type"`
}

// TransactionModule handles transaction-related API calls
type TransactionModule struct {
	client *Client
}

// List retrieves transactions for an account
func (m *TransactionModule) List(accountID string, options *TellerOptionsPagination) ([]TellerTransaction, error) {
	endpoint := fmt.Sprintf("%s/accounts/%s/transactions", m.client.baseURL, accountID)

	// Add pagination parameters if provided
	if options != nil {
		params := url.Values{}
		if options.Cursor != nil {
			params.Set("from_id", *options.Cursor)
		}
		if options.Limit != nil && *options.Limit > 0 {
			params.Set("count", strconv.Itoa(*options.Limit))
		}
		if options.StartDate != nil {
			params.Set("start_date", *options.StartDate)
		}
		if options.EndDate != nil {
			params.Set("end_date", *options.EndDate)
		}
		if len(params) > 0 {
			endpoint += "?" + params.Encode()
		}
	}

	req, err := http.NewRequest("GET", endpoint, nil)
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

	var result []TellerTransaction
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// Get retrieves a single transaction
func (m *TransactionModule) Get(accountID string, id string, options *TellerOptionsBase) (*TellerTransaction, error) {
	endpoint := fmt.Sprintf("%s/accounts/%s/transactions/%s", m.client.baseURL, accountID, id)

	req, err := http.NewRequest("GET", endpoint, nil)
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

	var result TellerTransaction
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
