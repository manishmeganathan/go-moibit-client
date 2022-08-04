package moibit

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
