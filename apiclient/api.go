package apiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.example.com/"
)

type APIError struct {
	StatusCode int    `json:"statuscode"`
	StatusDesc string `json:"statusdesc"`
	Message    string `json:"errormessage"`
}

type errorJsonResponse struct {
	Error *APIError `json:"error"`
}

func (r *APIError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}

type Client struct {
	User     string
	Password string
	client   *http.Client
	BaseURL  *url.URL
}

func NewClient(user, password string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:   http.DefaultClient,
		BaseURL:  baseURL,
		User:     user,
		Password: password,
	}

	return c
}

func (c *Client) NewRequest(method string, rsc string, params map[string]string) (*http.Request, error) {
	baseUrl, err := url.Parse(c.BaseURL.String() + rsc)
	if err != nil {
		return nil, err
	}

	if params != nil {
		ps := url.Values{}
		for k, v := range params {
			ps.Set(k, v)
		}
		baseUrl.RawQuery = ps.Encode()
	}

	req, err := http.NewRequest(method, baseUrl.String(), nil)
	req.SetBasicAuth(c.User, c.Password)

	return req, err
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if err := validateResponse(resp); err != nil {
		return resp, err
	}

	err = decodeResponse(resp, v)

	return resp, err
}

func decodeResponse(r *http.Response, v interface{}) error {
	if v == nil {
		return errors.New("nil interface provided to decodeResponse")
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	err := json.Unmarshal([]byte(bodyString), &v)

	return err
}

func validateResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	m := &errorJsonResponse{}
	err := json.Unmarshal([]byte(bodyString), &m)
	if err != nil {
		return err
	}

	return m.Error
}
