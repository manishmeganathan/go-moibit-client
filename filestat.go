package moibit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// FileStatus describes the status of file
type FileStatus struct {
	Active bool `json:"active"`
	Enable bool `json:"enable"`

	Hash          string `json:"hash"`
	Version       int    `json:"version"`
	Replication   int    `json:"replication"`
	FileSize      int    `json:"filesize"`
	EncryptionKey string `json:"encryptionKey"`
	LastUpdated   string `json:"lastUpdated"`

	Directory   string `json:"directory"`
	Path        string `json:"path"`
	NodeAddress string `json:"nodeAddress"`
}

// requestListFiles is the request for the ListFiles API of MOIBit
type requestListFiles struct {
	Path string `json:"path"`
}

// responseListFiles is the response for the ListFiles API of MOIBit
type responseListFiles struct {
	Metadata responseMetadata `json:"meta"`
	Data     []FileStatus     `json:"data"`
}

// ListFiles lists the files for a specified path.
// The files are returned as a slice of FileStatus objects.
// An error is returned if the API fails or the client cannot authenticate with MOIBit.
func (client *Client) ListFiles(path string) ([]FileStatus, error) {
	// Generate Request Data
	requestData, err := json.Marshal(requestListFiles{path})
	if err != nil {
		return nil, fmt.Errorf("request serialization failed: %w", err)
	}

	// Generate Request Object
	request, err := http.NewRequest("POST", urlListFiles, bytes.NewReader(requestData))
	if err != nil {
		return nil, fmt.Errorf("request generation failed: %w", err)
	}

	// Set authentication headers from the client
	client.setHeaders(request)

	// Perform the HTTP Request
	response, err := client.c.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check the status code of response
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("non-ok response: %v", response.StatusCode)
	}

	// Decode the response into a responseListFiles
	resp := new(responseListFiles)
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(resp); err != nil {
		return nil, fmt.Errorf("response decode failed: %w", err)
	}

	// Returns the file descriptors from the response
	return resp.Data, nil
}
