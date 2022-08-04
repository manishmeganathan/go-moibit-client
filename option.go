package moibit

import "net/http"

// ClientOption is a MOIBit client option provided to the Client Constructor
type ClientOption func(*Client) error

// defaultClient generates a new Client for a given public key, nonce and signature.
func defaultClient(pk, n, sig string) *Client {
	return &Client{http.Client{}, pk, n, sig, "", DefaultNetworkID}
}

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
