# go-moibit-client
Golang Client Library for interacting with MOIBit's Decentralized Storage API's

## Types
- <a  href="#Client"><code>Client</code></a>
- <a  href="#FileDescriptor"><code>FileDescriptor</code></a>
- <a  href="#AppDescriptor"><code>AppDescriptor</code></a>
- <a  href="#DevDescriptor"><code>DevDescriptor</code></a>


## Methods for File/Directory stats
- <a  href="#ListFiles"><code>ListFiles</code></a>
- <a  href="#FileStatus"><code>FileStatus</code></a>
- <a  href="#FileVersions"><code>FileVersions</code></a>

## Methods for Read/Write operations
- <a  href="#ReadFile"><code>ReadFile</code></a>
- <a  href="#WriteFile"><code>WriteFile</code></a>
- <a  href="#RemoveFile"><code>RemoveFile</code></a>
- <a  href="#MakeDirectory"><code>MakeDirectory</code></a>

## Methods for App/Dev details
- <a  href="#AppDetails"><code>AppDetails</code></a>
- <a  href="#DevDetails"><code>DevDetails</code></a>

<a name="Client"></a>
## Client
Client provides various methods to interact with MOIBit. 
``` go
type Client struct {
// contains filtered or unexported fields
}
```

<a name="AppDescriptor"></a>
## App Descriptor
App Descriptors holds metadata of an app registered with MOIBit.
```go
type AppDescriptor struct {
    IsActive  bool `json:"isActive"`
	IsRemoved bool `json:"isRemoved"`
    
    // app meta data
	AppID          string      `json:"appID"`
	AppName        string      `json:"appName"`
	AppDescription string      `json:"appDescription"`
	EndUsers       interface{} `json:"endUsers"`
    
    // network meta data
	NetworkID   string `json:"networkID"`
	NetworkName string `json:"networkName"`
    
	Replication    int         `json:"replication"`
	CanEncrypt     interface{} `json:"canEncrypt"`
	EncryptionType int         `json:"encryptionType"`
	CustomKey      interface{} `json:"customKey"`
	RecoveryTime   int64       `json:"recoveryTime"`
}
```
<a name="FileDescriptor"></a>
## File Descriptor
File Descriptor holds metadata of a file.
```go
type FileDescriptor struct {
	FileVersionDescriptor // inlined JSON

	Path        string `json:"path"`
	IsDirectory bool   `json:"isDir"`
	Directory   string `json:"directory"`
	NodeAddress string `json:"nodeAddress"`
}
```
<a name="DevDescriptor"></a>
## Dev Descriptor
Dev Descriptor holds metadata of a developer.
```go
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
```



<a name="ListFiles"></a>
### ListFiles(path string) ([]FileDescriptor, error)
ListFiles lists the files for a specified path,The files are returned as a slice of FileDescriptor objects.
An error is returned if the API fails or the client cannot authenticate with MOIBit.
```go
func (client *Client) ListFiles(path string) ([]FileDescriptor, error)
```

<a name="FileStatus"></a>
### FileStatus(path string) (FileDescriptor, error)
FileStatus returns the status of a file at a specified path, The returned FileStatus is empty if the file does not exist, which can be checked with Exists().
An error is returned if the API fails or the client cannot authenticate with MOIBit.
```go
func (client *Client) ListFiles(path string) ([]FileDescriptor, error)
```


<a name="FileDirectory"></a>
### FileVersions(path string) ([]FileVersionDescriptor, error)
FileVersions returns the version information of the file at the given path.
Returns a slice of FileVersionDescriptor objects, one for each version.
```go
func (client *Client) FileVersions(path string) ([]FileVersionDescriptor, error)
```


<a name="ReadFile"></a>
### ReadFile(path string, version int) ([]byte, error)
ReadFile reads a file from MOIBit at the given path for the given version.
Returns the []byte data of the file and an error.
```go
func (client *Client) ReadFile(path string, version int) ([]byte, error)
```

<a name="WriteFile"></a>
### WriteFile(data []byte, name string, opts ...WriteOption) ([]FileDescriptor, error)
WriteFile writes a given file to MOIBit. Accepts the file data as raw bytes and the file name.
It also accepts a variadic number of WriteOption to modify the write request.
Returns a []FileDescriptor (and error) containing the status of the file after successful write.
```go
func (client *Client) WriteFile(data []byte, name string, opts ...WriteOption) ([]FileDescriptor, error) 
```


<a name="RemoveFile"></a>
###  RemoveFile(path string, version int, opts ...RemoveOption) error
RemoveFile removes a file at the given path of the specified version.
It also accepts a variadic number of RemoveOption to modify the remove request.
- To remove directories, use the path to the directory and pass the RemoveDirectory option.
- To restore files, pass the file path and version to restore with the PerformRestore option.
This call is used to create and assign an activity to a user, this call return an error in the following cases.
```go
func (client *Client) RemoveFile(path string, version int, opts ...RemoveOption) error 
```

<a name="MakeDirectory"></a>
###  MakeDirectory(path string) error
MakeDirectory creates a new directory at the given path which can than be used for storing files.
```go
func (client *Client) MakeDirectory(path string) error 
```

<a name="AppDetails"></a>
###  AppDetails() (AppDescriptor, error) {
AppDetails returns the details of the application the client is configured for as a AppDescriptor object
```go
func (client *Client) AppDetails() (AppDescriptor, error)
```

<a name="DevDetails"></a>
###  DevDetails() (DevDescriptor, error) {
DevDetails returns the details of developer user the client is configured for as a DevDescriptor object
```go
func (client *Client) DevDetails() (DevDescriptor, error) 
```
