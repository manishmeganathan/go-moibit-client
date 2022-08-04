package moibit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// FileDescriptor describes the status of file
type FileDescriptor struct {
	Active bool `json:"active"`
	Enable bool `json:"enable"`

	Hash          string `json:"hash"`
	Version       int    `json:"version"`
	Replication   int    `json:"replication"`
	FileSize      int    `json:"filesize"`
	EncryptionKey string `json:"encryptionKey"`
	LastUpdated   string `json:"lastUpdated"`

	IsDirectory bool   `json:"isDir"`
	Directory   string `json:"directory"`
	Path        string `json:"path"`
	NodeAddress string `json:"nodeAddress"`
}

// Exists returns a boolean indicating if the
// file exists based on if the hash is empty.
func (file *FileDescriptor) Exists() bool {
	return file.Hash == "" && !file.IsDirectory
}

// requestListFiles is the request for the ListFiles API of MOIBit
type requestListFiles struct {
	Path string `json:"path"`
}

// responseListFiles is the response for the ListFiles API of MOIBit
type responseListFiles struct {
	Metadata responseMetadata `json:"meta"`
	Data     []FileDescriptor `json:"data"`
}

// ListFiles lists the files for a specified path.
// The files are returned as a slice of FileDescriptor objects.
// An error is returned if the API fails or the client cannot authenticate with MOIBit.
func (client *Client) ListFiles(path string) ([]FileDescriptor, error) {
	// Generate Request Data
	requestData, err := json.Marshal(requestListFiles{path})
	if err != nil {
		return nil, fmt.Errorf("request serialization failed: %w", err)
	}

	// Generate Request Object
	requestHTTP, err := http.NewRequest("POST", urlListFiles, bytes.NewReader(requestData))
	if err != nil {
		return nil, fmt.Errorf("request generation failed: %w", err)
	}

	// Set authentication headers from the client
	client.setHeaders(requestHTTP)

	// Perform the HTTP Request
	responseHTTP, err := client.c.Do(requestHTTP)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Decode the response into a responseListFiles
	response := new(responseListFiles)
	decoder := json.NewDecoder(responseHTTP.Body)
	if err := decoder.Decode(response); err != nil {
		return nil, fmt.Errorf("response decode failed [HTTP %v]: %w", responseHTTP.StatusCode, err)
	}

	// Check the status code of response
	if response.Metadata.StatusCode != 200 {
		return nil, fmt.Errorf("non-ok response [%v]: %v", response.Metadata.StatusCode, response.Metadata.Message)
	}

	// Returns the file descriptors from the response
	return response.Data, nil
}

// requestFileStatus is the request for the FileStatus API of MOIBit
type requestFileStatus struct {
	Path string `json:"path"`
}

// responseFileStatus is the response for the FileStatus API of MOIBit
type responseFileStatus struct {
	Metadata responseMetadata `json:"meta"`
	Data     FileDescriptor   `json:"data"`
}

// FileStatus returns the status of a file at a specified path.
// The returned FileStatus is empty if the file does not exist, which can be checked with Exists().
// An error is returned if the API fails or the client cannot authenticate with MOIBit.
func (client *Client) FileStatus(path string) (FileDescriptor, error) {
	// Generate Request Data
	requestData, err := json.Marshal(requestFileStatus{path})
	if err != nil {
		return FileDescriptor{}, fmt.Errorf("request serialization failed: %w", err)
	}

	// Generate Request Object
	requestHTTP, err := http.NewRequest("POST", urlFileStatus, bytes.NewReader(requestData))
	if err != nil {
		return FileDescriptor{}, fmt.Errorf("request generation failed: %w", err)
	}

	// Set authentication headers from the client
	client.setHeaders(requestHTTP)

	// Perform the HTTP Request
	responseHTTP, err := client.c.Do(requestHTTP)
	if err != nil {
		return FileDescriptor{}, fmt.Errorf("request failed: %w", err)
	}

	// Decode the response into a responseListFiles
	response := new(responseFileStatus)
	decoder := json.NewDecoder(responseHTTP.Body)
	if err := decoder.Decode(response); err != nil {
		return FileDescriptor{}, fmt.Errorf("response decode failed [HTTP %v]: %w", responseHTTP.StatusCode, err)
	}

	// Check the status code of response
	if response.Metadata.StatusCode != 200 {
		return FileDescriptor{}, fmt.Errorf("non-ok response [%v]: %v", response.Metadata.StatusCode, response.Metadata.Message)
	}

	// Returns the file descriptors from the response
	return response.Data, nil
}
