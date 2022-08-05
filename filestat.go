package moibit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// FileVersionDescriptor describes the version information of file
type FileVersionDescriptor struct {
	Active bool `json:"active"`
	Enable bool `json:"enable"`

	Hash           string `json:"hash"`
	ProvenanceHash string `json:"provenanceHash"`
	Version        int    `json:"version"`
	Replication    int    `json:"replication"`
	FileSize       int    `json:"filesize"`

	EncryptionKey string `json:"encryptionKey"`
	LastUpdated   string `json:"lastUpdated"`
}

// FileDescriptor describes the status of file
type FileDescriptor struct {
	FileVersionDescriptor // inlined JSON

	Path        string `json:"path"`
	IsDirectory bool   `json:"isDir"`
	Directory   string `json:"directory"`
	NodeAddress string `json:"nodeAddress"`
}

// Exists returns whether a file exists or not.
// Returns true if the file is a directory or has a non nil hash
func (file FileDescriptor) Exists() bool {
	if file.IsDirectory {
		// If file is a directory, return true
		return true
	} else if file.Hash != "" {
		// If file is not directory but hash non-nil hash, return true
		return true
	}

	// Return false otherwise
	return false
}

// String implements the Stringer interface for FileDescriptor
func (file FileDescriptor) String() string {
	if file.IsDirectory {
		return fmt.Sprintf("[Dirc] /%v", file.Directory)
	} else {
		return fmt.Sprintf("[File] %v%v", file.Directory, file.Path)
	}
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

	// Decode the response into a responseFileStatus
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

// requestFileVersions is the request for the FileVersions API of MOIBit
type requestFileVersions struct {
	Path string `json:"path"`
}

// responseFileVersions is the response for the FileVersions API of MOIBit
type responseFileVersions struct {
	Metadata responseMetadata        `json:"meta"`
	Data     []FileVersionDescriptor `json:"data"`
}

// FileVersions returns the version information of the file at the given path.
// Returns a slice of FileVersionDescriptor objects, one for each version.
func (client *Client) FileVersions(path string) ([]FileVersionDescriptor, error) {
	// Generate Request Data
	requestData, err := json.Marshal(requestFileVersions{path})
	if err != nil {
		return nil, fmt.Errorf("request serialization failed: %w", err)
	}

	// Generate Request Object
	requestHTTP, err := http.NewRequest("POST", urlFileVersions, bytes.NewReader(requestData))
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

	// Decode the response into a responseFileVersions
	response := new(responseFileVersions)
	decoder := json.NewDecoder(responseHTTP.Body)
	if err := decoder.Decode(response); err != nil {
		return nil, fmt.Errorf("response decode failed [HTTP %v]: %w", responseHTTP.StatusCode, err)
	}

	// Check the status code of response
	if response.Metadata.StatusCode != 200 {
		return nil, fmt.Errorf("non-ok response [%v]: %v", response.Metadata.StatusCode, response.Metadata.Message)
	}

	// Returns the file version descriptors from the response
	return response.Data, nil
}

// responseMakeDir is the response for the MakeDir API of MOIBit
type responseMakeDir struct {
	Metadata responseMetadata `json:"meta"`
	Data     string           `json:"data"`
}

// MakeDirectory creates a new directory at the given path which can than be used for storing files.
func (client *Client) MakeDirectory(path string) error {
	// Generate Request Object
	requestHTTP, err := http.NewRequest("GET", urlMakeDir, nil)
	if err != nil {
		return fmt.Errorf("request generation failed: %w", err)
	}

	// Set given path to query parameters
	query := requestHTTP.URL.Query()
	query.Add("path", path)
	requestHTTP.URL.RawQuery = query.Encode()

	// Set authentication headers from the client
	client.setHeaders(requestHTTP)

	// Perform the HTTP Request
	responseHTTP, err := client.c.Do(requestHTTP)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	// Decode the response into a responseMakeDir
	response := new(responseMakeDir)
	decoder := json.NewDecoder(responseHTTP.Body)
	if err := decoder.Decode(response); err != nil {
		return fmt.Errorf("response decode failed [HTTP %v]: %w", responseHTTP.StatusCode, err)
	}

	// Check the status code of response
	if response.Metadata.StatusCode != 200 {
		return fmt.Errorf("non-ok response [%v]: %v | %v", response.Metadata.StatusCode, response.Metadata.Message, response.Data)
	}

	return nil
}
