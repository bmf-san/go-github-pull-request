package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const defaultBaseURL = "https://api.github.com/"

// Client is an api client.
type Client struct {
	Pulls *HandlerPull
	Repos *HandlerRepo
}

// HTTPClient is a HTTPClient.
type HTTPClient struct {
	baseURL *url.URL
	client  *http.Client
}

// NewClient creates a client.
func NewClient(token string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	httpClient := NewHTTPClient(tc)
	ph := NewHandlerPull(httpClient)
	r := NewHandlerRepo(httpClient)
	return &Client{
		Pulls: ph,
		Repos: r,
	}
}

// NewHTTPClient creates a HTTPClient
func NewHTTPClient(httpClient *http.Client) *HTTPClient {
	baseURL, _ := url.Parse(defaultBaseURL)
	return &HTTPClient{
		baseURL: baseURL,
		client:  httpClient,
	}
}

// NewClient creates a request.
func (hc *HTTPClient) NewRequest(method string, header http.Header, body []byte, url string) (*http.Request, error) {
	u, err := hc.baseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}

	return req, nil
}

// Do
func (hc *HTTPClient) Do(req *http.Request) ([]byte, error) {
	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
