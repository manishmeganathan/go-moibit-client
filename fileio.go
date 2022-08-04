package moibit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// requestReadFile is the request for the ReadFile API of MOIBit
type requestReadFile struct {
	FileName string `json:"fileName"`
	Version  int    `json:"version"`
}

// ReadFile reads a file from MOIBit at the given path for the given version.
// Returns the []byte data of the file and an error.
func (client *Client) ReadFile(path string, version int) ([]byte, error) {
	// Generate Request Data
	requestData, err := json.Marshal(requestReadFile{path, version})
	if err != nil {
		return nil, fmt.Errorf("request serialization failed: %w", err)
	}

	// Generate Request Object
	requestHTTP, err := http.NewRequest("POST", urlReadFile, bytes.NewReader(requestData))
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

	// Check the status code of response
	if responseHTTP.StatusCode != 200 {
		return nil, fmt.Errorf("non-ok response [%v]", responseHTTP.StatusCode)
	}

	// Read all bytes from the response body
	data, err := io.ReadAll(responseHTTP.Body)
	if err != nil {
		return nil, fmt.Errorf("response data spool: %w", err)
	}

	return data, nil
}

// requestWriteFile is the request for the WriteFile API of MOIBit
type requestWriteFile struct {
	FileText string `json:"text"`
	FileName string `json:"fileName"`

	KeepPrevious  bool `json:"keepPrevious"`
	CreateFolders bool `json:"createFolders"`
	IsProvenance  bool `json:"isProvenance"`

	Replication int            `json:"replication,omitempty"`
	Encryption  EncryptionType `json:"encryptionType,omitempty"`
}

// defaultWriteFileRequest generates a new requestWriteFile object for the given file name and data
func defaultWriteFileRequest(data []byte, name string) *requestWriteFile {
	return &requestWriteFile{
		FileName: name, FileText: string(data),
		KeepPrevious: false, CreateFolders: true, IsProvenance: false,
	}
}

// responseWriteFile is the response for the WriteFile API of MOIBit
type responseWriteFile struct {
	Metadata responseMetadata `json:"meta"`
	Data     []FileDescriptor `json:"data"`
}

// WriteFile writes a given file to MOIBit. Accepts the file data as raw bytes and the file name.
// It also accepts a variadic number of WriteOption to modify the write request.
// Returns a []FileDescriptor (and error) containing the status of the file after successful write.
func (client *Client) WriteFile(data []byte, name string, opts ...WriteOption) ([]FileDescriptor, error) {
	// Generate Request Data
	request := defaultWriteFileRequest(data, name)
	for _, opt := range opts {
		if err := opt(request); err != nil {
			return nil, fmt.Errorf("request creation failed while applying options: %w", err)
		}
	}

	// Serialize Request Data
	requestData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("request serialization failed: %w", err)
	}

	// Generate Request Object
	requestHTTP, err := http.NewRequest("POST", urlWriteFile, bytes.NewReader(requestData))
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

	// Decode the response into a responseWriteFiles
	response := new(responseWriteFile)
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

// UnmarshalJSON implements the json.Unmarshaler interface for responseWriteFile.
// The custom unmarshaler is required because of a bug in the WriteTextToFile
// API causing to return the file descriptors in a deformed string form.
func (resp *responseWriteFile) UnmarshalJSON(data []byte) error {
	// Declare an intermediate representation for the deformed response data.
	// The meta key of the response is not deformed but the data key is returned
	// as an array of strings, with each string representing the string version
	// of the JSON output for an array of FileDescriptors.
	var ir struct {
		Meta responseMetadata `json:"meta"`
		Data []string         `json:"data"`
	}

	// Attempt to unmarshal data into the IR
	if err := json.Unmarshal(data, &ir); err != nil {
		return fmt.Errorf("failed decode 'responseWriteFile' into intermediate representation: %w", err)
	}

	// Iterate over the strings in the deformed data
	fileDescriptors := make([]FileDescriptor, 0)
	for _, descriptor := range ir.Data {
		// Attempt to unmarshal each string into a slice of FileDescriptors
		// Ideally, this should only be one string but this accommodates any further deformations.
		var fd []FileDescriptor
		if err := json.Unmarshal([]byte(descriptor), &fd); err != nil {
			return fmt.Errorf("failed to decode 'data' into []FileDescriptor: %w", err)
		}

		// Append the file descriptors into the super set
		fileDescriptors = append(fileDescriptors, fd...)
	}

	// Update the response fields
	resp.Metadata, resp.Data = ir.Meta, fileDescriptors
	return nil
}

// WriteOption is a request option for the WriteFile method of Client.
type WriteOption func(*requestWriteFile) error

// KeepPrevious returns a WriteOption that can be used to preserve the
// versioning of the file that is written, in case the file already exists.
func KeepPrevious() WriteOption {
	return func(request *requestWriteFile) error {
		request.KeepPrevious = true
		return nil
	}
}

// CreateFolders returns a WriteOption that can be used to specify that the file write
// should create all folders (that do not exist) specified in the path of the file name
func CreateFolders() WriteOption {
	return func(request *requestWriteFile) error {
		request.CreateFolders = true
		return nil
	}
}

// CreateOnlyFile returns a WriteOption that can be used to specify that the file write
// should fail in case folders specified in the path of the file do not exist already.
func CreateOnlyFile() WriteOption {
	return func(request *requestWriteFile) error {
		request.CreateFolders = false
		return nil
	}
}

// Provenance returns a WriteOption that can be used to specify that the
// proof of the file needs to be stored on MOI's Indus Provenance Network.
func Provenance() WriteOption {
	return func(request *requestWriteFile) error {
		request.IsProvenance = true
		return nil
	}
}

// ReplicationFactor returns a WriteOption that can be used to specify
// the number of replications for the written file on its network.
func ReplicationFactor(n int) WriteOption {
	return func(request *requestWriteFile) error {
		request.Replication = n
		return nil
	}
}

// ApplyEncryption returns a WriteOption that can specify the encryption
// scheme for the file while being written to MOIBit.
func ApplyEncryption(encryption EncryptionType) WriteOption {
	return func(request *requestWriteFile) error {
		request.Encryption = encryption
		return nil
	}
}

// requestRemoveFile is the request for the Remove API of MOIBit
type requestRemoveFile struct {
	FilePath    string `json:"path"`
	Version     int    `json:"version"`
	IsDirectory bool   `json:"isdir"`
	Operation   int    `json:"operationType"`
}

// defaultRemoveFileRequest generates a new requestRemoveFile object for the given file path and version
func defaultRemoveFileRequest(path string, version int) *requestRemoveFile {
	return &requestRemoveFile{
		FilePath: path, Version: version,
		IsDirectory: false, Operation: 0,
	}
}

// RemoveFile removes a file at the given path of the specified version.
// It also accepts a variadic number of RemoveOption to modify the remove request.
// 		- To remove directories, use the path to the directory and pass the RemoveDirectory option.
// 		- To restore files, pass the file path and version to restore with the PerformRestore option.
func (client *Client) RemoveFile(path string, version int, opts ...RemoveOption) error {
	// Generate Request Data
	request := defaultRemoveFileRequest(path, version)
	for _, opt := range opts {
		if err := opt(request); err != nil {
			return fmt.Errorf("request creation failed while applying options: %w", err)
		}
	}

	// Serialize Request Data
	requestData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("request serialization failed: %w", err)
	}

	// Generate Request Object
	requestHTTP, err := http.NewRequest("POST", urlRemoveFile, bytes.NewReader(requestData))
	if err != nil {
		return fmt.Errorf("request generation failed: %w", err)
	}

	// Set authentication headers from the client
	client.setHeaders(requestHTTP)

	// Perform the HTTP Request
	responseHTTP, err := client.c.Do(requestHTTP)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	// Check the status code of response
	if responseHTTP.StatusCode != 200 {
		return fmt.Errorf("non-ok response [%v]", responseHTTP.StatusCode)
	}

	return nil
}

// RemoveOption is a request option for the RemoveFile method of Client.
type RemoveOption func(*requestRemoveFile) error

// RemoveDirectory returns a RemoveOption that will specify that the file to delete is a directory.
// Note: This will cause RemoveFile to fail if it is a file and not a directory.
func RemoveDirectory() RemoveOption {
	return func(request *requestRemoveFile) error {
		request.IsDirectory = true
		return nil
	}
}

// PerformRestore returns a RemoveOption that will set the operation mode of RemoveFile
// to restoration, which will result in MOIBit attempting to restore the file version
func PerformRestore() RemoveOption {
	return func(request *requestRemoveFile) error {
		request.Operation = 1
		return nil
	}
}
