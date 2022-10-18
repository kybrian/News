package news

import "net/http"

// Client working with the news API
type Client struct {
	http     *http.Client // points to the HTTP Client that should be used to make requests
	key      string       // Holds the api key
	PageSize int          // Holds the number of results to return per page
}

// Create a new client instant for making requests to the news API
func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{httpClient, key, pageSize}
}
