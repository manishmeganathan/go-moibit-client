package moibit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// DefaultNetworkID represents the network ID of the default MOIBit Network.
	// Can be overridden using the NetworkID() option while constructing the client
	DefaultNetworkID = "12D3KooWSMAGyrB9TG45AAWaQNJmMdfJpnLQ5e1XM21hkm3FokHk"

	// DefaultBaseURL represents the default base URL for the MOIBit Client.
	// It is the primary production service endpoint. API endpoints are built on this base URL
	DefaultBaseURL = "https://api.moinet.io/moibit/v1"
)

// ClientOption is a MOIBit client option provided to the Client Constructor
type ClientOption func(*Client) error

// AppID returns ClientOption that can be used to set the App ID for a Client
func AppID(app string) ClientOption {
	return func(config *Client) error {
		config.appID = app
		return nil
	}
}

// NetworkID returns a ClientOption that can be used to set the Network ID for a Client
func NetworkID(net string) ClientOption {
	return func(config *Client) error {
		config.netID = net
		return nil
	}
}

// BaseURL returns a ClientOption that can be used to set a Base URL for a Client
// The Base URL can be used to have the Client dial a local development instance or custom network.
func BaseURL(url string) ClientOption {
	return func(client *Client) error {
		client.url = url
		return nil
	}
}

// Client represents a MOIBit API Client
type Client struct {
	c   http.Client
	url string

	pubkey    string
	nonce     string
	signature string

	appID string
	netID string
}

// NewClient creates a new MOIBit API Client for the given signature and nonce
// Accepts a variadic number of ClientOption arguments to set the App ID, Network ID or Base URL
// Uses the DefaultNetworkID, DefaultBaseURL and no App ID, by default.
func NewClient(signature, nonce string, opts ...ClientOption) (*Client, error) {
	// Generate the default client with the nonce and signature
	client := defaultClient(signature, nonce)
	// Apply the options on the client config
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("option could not be applied: %w", err)
		}
	}

	// Authenticate credentials and get public key
	pubkey, err := Authenticate(client)
	if err != nil {
		return nil, fmt.Errorf("user could not be authenticated: %w", err)
	}

	// Set the pubkey of the client from the authenticated public key
	client.pubkey = pubkey
	return client, nil
}

// defaultClient generates a new Client for a given public key, nonce and signature.
func defaultClient(sig, n string) *Client {
	return &Client{
		http.Client{}, DefaultBaseURL,
		"", n, sig,
		"", DefaultNetworkID,
	}
}

// serviceURL generates a URL with the given endpoint concatenated with the base URL of the Client.
// Note: Endpoint should start with "/" as this as a pure concat operation on the
func (client *Client) serviceURL(endpoint string) string {
	return client.url + endpoint
}

// Authenticate attempts to authenticate a set of credentials with MOIBit.
// Accepts the nonce and signature of the developer and returns the public or an
// error if either the authentication routine fails or if the credentials are invalid.
func Authenticate(client *Client) (string, error) {
	// Create new POST request for user authentication
	request, err := http.NewRequest("POST", client.serviceURL("/user/auth"), nil)
	if err != nil {
		return "", fmt.Errorf("request generation failed: %w", err)
	}

	// Set the auth headers (signature and nonce)
	request.Header.Set("nonce", client.nonce)
	request.Header.Set("signature", client.signature)

	// Perform the request
	response, err := client.c.Do(request)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	// Check the status code of response
	if response.StatusCode != 200 {
		return "", fmt.Errorf("user not authenticated: %v", response.StatusCode)
	}

	// Decode the response into a responseUserAuth
	auth := new(responseUserAuth)
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(auth); err != nil {
		return "", fmt.Errorf("response decode failed: %w", err)
	}

	return auth.Data.Address, nil
}

// responseMetadata is a generic response data structure
// attached to every response for its metadata field
type responseMetadata struct {
	// HTTP Status Code of the response
	StatusCode int `json:"code"`
	// Request of the corresponding response
	RequestID string `json:"requestID"`
	// Metadata message returned to the response
	Message string `json:"message"`
}

// responseUserAuth is the response from the User/Auth API of MOIBit
type responseUserAuth struct {
	Metadata responseMetadata `json:"meta"`
	Data     struct {
		Address string `json:"address"`
		Entropy string `json:"entropy"`
	} `json:"data"`
}

// setHeaders accepts a HTTP Request and sets the headers
// "developerKey", "networkID" and "appID" from the client.
func (client *Client) setHeaders(request *http.Request) {
	request.Header.Set("nonce", client.nonce)
	request.Header.Set("signature", client.signature)
	request.Header.Set("developerKey", client.pubkey)
	request.Header.Set("networkID", client.netID)
	request.Header.Set("appID", client.appID)
}
