package moibit

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AppDescriptor describes the status of an application
type AppDescriptor struct {
	IsActive  bool `json:"isActive"`
	IsRemoved bool `json:"isRemoved"`

	AppID          string      `json:"appID"`
	AppName        string      `json:"appName"`
	AppDescription string      `json:"appDescription"`
	EndUsers       interface{} `json:"endUsers"`

	NetworkID   string `json:"networkID"`
	NetworkName string `json:"networkName"`

	Replication    int         `json:"replication"`
	CanEncrypt     interface{} `json:"canEncrypt"`
	EncryptionType int         `json:"encryptionType"`
	CustomKey      interface{} `json:"customKey"`
	RecoveryTime   int64       `json:"recoveryTime"`
}

// responseAppDetails is the response for the AppDetails API of MOIBit
type responseAppDetails struct {
	Metadata responseMetadata `json:"meta"`
	Data     AppDescriptor    `json:"data"`
}

// AppDetails returns the details of the application the client is configured for as a AppDescriptor object
func (client *Client) AppDetails() (AppDescriptor, error) {
	if client.appID == "" {
		return AppDescriptor{}, fmt.Errorf("request failed: no appID set for client")
	}

	// Generate Request Object
	requestHTTP, err := http.NewRequest("GET", client.serviceURL("/appdetails"), nil)
	if err != nil {
		return AppDescriptor{}, fmt.Errorf("request generation failed: %w", err)
	}

	// Set authentication headers from the client
	client.setHeaders(requestHTTP)

	// Perform the HTTP Request
	responseHTTP, err := client.c.Do(requestHTTP)
	if err != nil {
		return AppDescriptor{}, fmt.Errorf("request failed: %w", err)
	}

	// Decode the response into a responseAppDetails
	response := new(responseAppDetails)
	decoder := json.NewDecoder(responseHTTP.Body)
	if err := decoder.Decode(response); err != nil {
		return AppDescriptor{}, fmt.Errorf("response decode failed [HTTP %v]: %w", responseHTTP.StatusCode, err)
	}

	// Check the status code of response
	if response.Metadata.StatusCode != 200 {
		return AppDescriptor{}, fmt.Errorf("non-ok response [%v]: %v", response.Metadata.StatusCode, response.Metadata.Message)
	}

	// Returns the file descriptors from the response
	return response.Data, nil
}

// DevDescriptor describes the status of a developer/user
type DevDescriptor struct {
	Active bool        `json:"active"`
	Key    interface{} `json:"key"`

	Name  string `json:"name"`
	Email string `json:"email"`

	Apps []struct {
		IsActive  bool `json:"isActive"`
		IsRemoved bool `json:"isRemoved"`

		AppID   string `json:"appID"`
		AppName string `json:"appName"`

		Replication    int    `json:"replication"`
		EncryptionType int    `json:"encryptionType"`
		EncryptionAlgo string `json:"encryptionAlgo"`
		RecoveryTime   int64  `json:"recoveryTime"`

		NetworkID   string `json:"networkID"`
		NetworkName string `json:"networkName"`
	} `json:"apps"`

	Networks   interface{} `json:"networks"`
	DevPubKey  interface{} `json:"devPubKey"`
	Encryption interface{} `json:"encryption"`

	IsActive         bool `json:"isActive"`
	Creditcard       bool `json:"creditcard"`
	Canencrypt       bool `json:"canencrypt"`
	Canreplicate     bool `json:"canreplicate"`
	Cancreatenetwork bool `json:"cancreatenetwork"`

	Maxstorage           int    `json:"maxstorage"`
	ReplicationFactor    int    `json:"replicationFactor"`
	Plan                 int    `json:"plan"`
	StripeCustomerID     string `json:"stripeCustomerID"`
	StripeSubscriptionID string `json:"stripeSubscriptionID"`

	NoOfPremiumNodes int         `json:"noOfPremiumNodes"`
	NoOfApps         int         `json:"noOfApps"`
	PremiumNodesList interface{} `json:"premiumNodesList"`

	Credit               int  `json:"credit"`
	FreeTrial            bool `json:"freeTrial"`
	AnnualSubscription   bool `json:"AnnualSubscription"`
	FreeTrialJoiningDate int  `json:"freeTrialJoiningDate"`
}

// responseDevDetails is the response for the DevDetails API of MOIBit
type responseDevDetails struct {
	Metadata responseMetadata `json:"meta"`
	Data     DevDescriptor    `json:"data"`
}

// DevDetails returns the details of developer user the client is configured for as a DevDescriptor object
func (client *Client) DevDetails() (DevDescriptor, error) {
	// Generate Request Object
	requestHTTP, err := http.NewRequest("GET", client.serviceURL("/devstat"), nil)
	if err != nil {
		return DevDescriptor{}, fmt.Errorf("request generation failed: %w", err)
	}

	// Set authentication headers from the client
	client.setHeaders(requestHTTP)

	// Perform the HTTP Request
	responseHTTP, err := client.c.Do(requestHTTP)
	if err != nil {
		return DevDescriptor{}, fmt.Errorf("request failed: %w", err)
	}

	// Decode the response into a responseAppDetails
	response := new(responseDevDetails)
	decoder := json.NewDecoder(responseHTTP.Body)
	if err := decoder.Decode(response); err != nil {
		return DevDescriptor{}, fmt.Errorf("response decode failed [HTTP %v]: %w", responseHTTP.StatusCode, err)
	}

	// Check the status code of response
	if response.Metadata.StatusCode != 200 {
		return DevDescriptor{}, fmt.Errorf("non-ok response [%v]: %v", response.Metadata.StatusCode, response.Metadata.Message)
	}

	// Returns the file descriptors from the response
	return response.Data, nil
}
