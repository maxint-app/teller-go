package teller

import (
	"crypto/tls"
	"net/http"
)

// Client is the main Teller API client
type Client struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
	certPath    string
	keyPath     string

	// Modules
	Identity     *IdentityModule
	Account      *AccountModule
	Transactions *TransactionModule
	Institutions *InstitutionsModule
}

// NewClient creates a new Teller API client
func NewClient(certPath, keyPath string, accessToken *string) (*Client, error) {
	// Load certificates for mutual TLS
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	token := ""
	if accessToken != nil {
		token = *accessToken
	}

	c := &Client{
		baseURL:     "https://api.teller.io",
		accessToken: token,
		httpClient:  httpClient,
		certPath:    certPath,
		keyPath:     keyPath,
	}

	// Initialize modules
	c.Identity = &IdentityModule{client: c}
	c.Account = &AccountModule{client: c}
	c.Transactions = &TransactionModule{client: c}
	c.Institutions = &InstitutionsModule{client: c}

	return c, nil
}
