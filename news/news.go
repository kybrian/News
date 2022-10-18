package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Article struct {
	Source struct {
		ID   interface{} `json:"id"`
		Name string      `json:"name"`
	} `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

type Results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

// Client working with the news API
type Client struct {
	http     *http.Client // points to the HTTP Client that should be used to make requests
	key      string       // Holds the api key
	PageSize int          // Holds the number of results to return per page
}

func(c *Client) FetchEverything(query, page string)(*Results, error){
	endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en", url.QueryEscape(query), c.PageSize, page, c.key)
	resp, err := c.http.Get(endpoint)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err 
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &Results{} 
	return res, json.Unmarshal(body, res)

}

// Create a new client instant for making requests to the news API
func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{httpClient, key, pageSize}
}
