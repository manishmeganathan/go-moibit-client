package moibit

// EncryptionType represents an enumeration for the
// types of Encryption Schemes supported by MOIBit
type EncryptionType int

const (
	// NoEncryption applies no encryption on the file when storing on MOIBit
	NoEncryption EncryptionType = iota - 1

	// DefaultNetworkEncryption applies the default network
	// encryption scheme on the file when storing on MOIBit.
	DefaultNetworkEncryption

	// DeveloperKeyEncryption applies the default encryption scheme defined for
	// the user/developer authenticated with client making the write request
	DeveloperKeyEncryption

	// EndUserKeyEncryption applies the encryption scheme defined by
	// the end user's key for whom the file is being stored.
	EndUserKeyEncryption

	// CustomKeyEncryption applies the custom encryption scheme defined
	// for the application on the file when storing on MOIBit
	CustomKeyEncryption

	// MESEncryption applies the Modern Encryption Standard
	// on the file when storing on MOIBit.
	MESEncryption
)
