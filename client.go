package moibit

import "net/http"

// DefaultNetworkID represents the network ID of the default MOIBit Network.
// Can be overridden using the NetworkID() option while constructing the client
const DefaultNetworkID = "12D3KooWSMAGyrB9TG45AAWaQNJmMdfJpnLQ5e1XM21hkm3FokHk"

// Client represents a MOIBit API Client
type Client struct {
	c http.Client

	pubkey    string
	nonce     string
	signature string

	appID string
	netID string
}

// NewClient creates a new MOIBit API Client for the given public key, nonce and signature.
// Accepts a variadic number of ClientOption arguments to set the App or Network ID.
// Uses the DefaultNetworkID and no App ID, by default.
func NewClient(pubkey, nonce, signature string, opts ...ClientOption) (*Client, error) {
	// Generate the default client for the pubkey, nonce and signature
	client := defaultClient(pubkey, nonce, signature)
	// Apply the options on the client config
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	// TODO: Authenticate User

	return client, nil
}
