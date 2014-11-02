package apiclient

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.foo.com/"
)

type Client struct {
	User     string
	Password string
	client   *http.Client
	BaseURL  *url.URL
}

func NewClient(httpClient *http.Client) *Client {

	if httpClient != nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:   httpClient,
		BaseURL:  baseURL,
		User:     "",
		Password: "",
	}

	return c
}
