package way

import "io"

// Crypto is an interface that defines cryptographic operations.
type Crypto interface {
	// HashString calculates the hash of a string and returns a fixed-size byte array.
	HashString(data string) []byte

	// HashStringWithSalt calculates the hash of a string with a salt and returns a fixed-size byte array.
	HashStringWithSalt(data, salt string) []byte

	// HashByte calculates the hash of a byte slice and returns a fixed-size byte array.
	HashByte(data []byte) []byte

	// Encrypt encrypts the given byte slice using the provided passphrase and returns the encrypted data as a string.
	Encrypt(data []byte, passphrase string) ([]byte, error)

	// EncryptStream encrypts the given byte slice using the provided passphrase and writes the encrypted data to the provided writer.
	EncryptStream(data []byte, passphrase string, writer io.Writer) error

	// Decrypt decrypts the given encrypted string using the provided passphrase and returns the decrypted data as a byte slice.
	Decrypt(encrypted string, passphrase string) ([]byte, error)

	// DecryptStream decrypts the given encrypted string using the provided passphrase and writes the decrypted data to the provided writer.
	DecryptStream(encrypted string, passphrase string, writer io.Writer) error
}
