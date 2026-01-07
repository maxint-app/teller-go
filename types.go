package teller

// TellerOptionsBase represents base options for Teller API requests
type TellerOptionsBase struct {
	AccessToken string
}

// TellerOptionsPagination represents pagination options for Teller API requests
type TellerOptionsPagination struct {
	TellerOptionsBase
	Cursor string
	Limit  int
}
