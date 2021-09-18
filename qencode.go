package qencode

import "net/http"

// Client is a client to abstract communication
// with qEncode APIs. To create a new Client, use NewClient.
type Client struct {
	baseURL string
	hClient *http.Client
}

// NewClient creates a new Client. It requires a http.Client (if nil
// the default will be used).
func NewClient(cl *http.Client) *Client {
	if cl == nil {
		cl = http.DefaultClient
	}
	return &Client{
		hClient: cl,
		baseURL: "https://api.qencode.com",
	}
}
