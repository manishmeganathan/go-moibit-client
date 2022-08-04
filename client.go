package moibit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DefaultNetworkID represents the network ID of the default MOIBit Network.
// Can be overridden using the NetworkID() option while constructing the client
const DefaultNetworkID = "12D3KooWSMAGyrB9TG45AAWaQNJmMdfJpnLQ5e1XM21hkm3FokHk"

// URLs for all MOIBit API Endpoints
const (
	urlAuthUser  = "https://api.moinet.io/moibit/v1/user/auth"
	urlListFiles = "https://api.moinet.io/moibit/v1/listfiles"
)

// Client represents a MOIBit API Client
type Client struct {
	c http.Client

	pubkey    string
	nonce     string
	signature string

	appID string
	netID string
}

// NewClient creates a new MOIBit API Client for the given signature and nonce
// Accepts a variadic number of ClientOption arguments to set the App or Network ID.
// Uses the DefaultNetworkID and no App ID, by default.
func NewClient(signature, nonce string, opts ...ClientOption) (*Client, error) {
	// Generate the default client with the nonce and signature
	client := defaultClient(signature, nonce)
	// Apply the options on the client config
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("option could not be applied: %w", err)
		}
	}

	// Authenticate credentials
	if err := client.authenticate(); err != nil {
		return nil, fmt.Errorf("user could not be authenticated: %w", err)
	}

	return client, nil
}

// defaultClient generates a new Client for a given public key, nonce and signature.
func defaultClient(sig, n string) *Client {
	return &Client{http.Client{}, "", n, sig, "", DefaultNetworkID}
}

// authenticate attempts to authenticate the Client credentials with MOIBit.
// Returns an error if either the authentication routine fails or if the credentials are invalid.
func (client *Client) authenticate() error {
	// Create new POST request for user authentication
	request, err := http.NewRequest("POST", urlAuthUser, nil)
	if err != nil {
		return fmt.Errorf("request generation failed: %w", err)
	}

	// Set the auth headers (signature and nonce)
	client.setAuthHeaders(request)

	// Perform the request
	response, err := client.c.Do(request)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	// Check the status code of response
	if response.StatusCode != 200 {
		return fmt.Errorf("user not authenticated: %v", response.StatusCode)
	}

	// Decode the response into a responseUserAuth
	auth := new(responseUserAuth)
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(auth); err != nil {
		return fmt.Errorf("response decode failed: %w", err)
	}

	// Set the pubkey of the client from the address in the response
	client.pubkey = auth.Data.Address
	return nil
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

// setAuthHeaders accepts a HTTP Request and sets the auth headers
// "nonce" and "signature" for authenticating with MOIBit.
func (client *Client) setAuthHeaders(request *http.Request) {
	request.Header.Set("nonce", client.nonce)
	request.Header.Set("signature", client.signature)
}

// setHeaders accepts a HTTP Request and sets the headers
// "developerKey", "networkID" and "appID" from the client.
func (client *Client) setHeaders(request *http.Request) {
	client.setAuthHeaders(request)

	request.Header.Set("developerKey", client.pubkey)
	request.Header.Set("networkID", client.netID)
	request.Header.Set("appID", client.appID)
}
