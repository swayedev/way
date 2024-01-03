package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/sha3"
)

type Crypto interface {
	HashString(data string) [32]byte
	HashByte(data []byte) [32]byte
	Encrypt(data []byte, passphrase string) (string, error)
	Decrypt(encrypted string, passphrase string) ([]byte, error)
}

// HashStringToString takes a string as input and returns its SHA3-256 hash as a hexadecimal string.
func HashStringToString(data string) string {
	hashArray := sha3.Sum256([]byte(data))
	return hex.EncodeToString(hashArray[:])
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
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given encrypted string using the provided passphrase.
// It returns the decrypted data as a byte slice.
// If an error occurs during decryption, it returns nil and the corresponding error.
func Decrypt(encrypted string, passphrase string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
