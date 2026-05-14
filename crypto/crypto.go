package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/swayedev/fcrypt"
	"golang.org/x/crypto/sha3"
)

// Crypto is an interface that defines cryptographic operations.
type Crypto interface {
	// HashString calculates the hash of a string and returns a fixed-size byte array.
	HashString(data string) [32]byte

	// HashByte calculates the hash of a byte slice and returns a fixed-size byte array.
	HashByte(data []byte) [32]byte

	// Encrypt encrypts the given byte slice using the provided passphrase and returns the encrypted data as a string.
	Encrypt(data []byte, passphrase string) (string, error)

	// Decrypt decrypts the given encrypted string using the provided passphrase and returns the decrypted data as a byte slice.
	Decrypt(encrypted string, passphrase string) ([]byte, error)
}

// HashStringToString takes a string as input and returns its SHA3-256 hash as a hexadecimal string.
func HashStringToString(data string) string {
	return fcrypt.HashStringToStringSHA3(data)
}

// HashString calculates the SHA3-256 hash of the input string.
// It takes a string as input and returns a fixed-size array of 32 bytes.
func HashString(data string) [32]byte {
	return sha3.Sum256([]byte(data))
}

// HashByte calculates the SHA3-256 hash of the given byte slice.
// It returns a fixed-size array of 32 bytes representing the hash.
func HashByte(data []byte) [32]byte {
	return sha3.Sum256([]byte(data))
}

// Encrypt encrypts the given data using the provided passphrase.
// It returns the encrypted data as a hexadecimal string and any error encountered.
func Encrypt(data []byte, passphrase string) (string, error) {
	ciphertext, err := fcrypt.Encrypt(data, []byte(passphrase))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given encrypted string using the provided passphrase.
// It returns the decrypted data as a byte slice.
// If an error occurs during decryption, it returns nil and the corresponding error.
func Decrypt(encrypted string, passphrase string) ([]byte, error) {
	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	return fcrypt.Decrypt(data, []byte(passphrase))
}

// GenerateRandomKey generates a random key of the specified length.
// It uses the crypto/rand package to generate random bytes and returns the key as a byte slice.
// If an error occurs during the generation process, it returns an error.
func GenerateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}
