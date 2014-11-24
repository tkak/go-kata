package apiclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

// setup sets up a test HTTP server along with a apiclient.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// api client configured to use test server
	client = NewClient("foo", "bar")
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	c := NewClient("user", "password")

	if got, want := c.client, http.DefaultClient; got != want {
		t.Errorf("NewClient client = %v, want %v", got, want)
	}

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}

	if got, want := c.User, "user"; got != want {
		t.Errorf("NewClient User is %v, want %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	setup()
	defer teardown()

	req, err := client.NewRequest("GET", "/checks", nil)
	if err != nil {
		t.Errorf("NewRequest returned error: %v", err)
	}

	if got, want := req.Method, "GET"; got != want {
		t.Errorf("NewRequest Method returned %+v, want %+v", got, want)
	}

	if got, want := req.URL.String(), client.BaseURL.String()+"/checks"; got != want {
		t.Errorf("NewRequest URL returned %+v, want %+v", got, want)
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("Request method = %v, want %v", got, want)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/foo", nil)
	body := new(foo)
	client.Do(req, body)

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}
