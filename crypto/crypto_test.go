package crypto

import (
	"bytes"
	"errors"
	"testing"

	"github.com/swayedev/fcrypt"
)

func TestEncryptDecryptUsesFcrypt(t *testing.T) {
	key := string(bytes.Repeat([]byte{1}, fcrypt.DefaultKeyLength))
	plaintext := []byte("hello way")

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("Decrypt() = %q, want %q", decrypted, plaintext)
	}
}

func TestDecryptShortCiphertextReturnsFcryptError(t *testing.T) {
	key := string(bytes.Repeat([]byte{1}, fcrypt.DefaultKeyLength))
	_, err := Decrypt("00", key)
	if !errors.Is(err, fcrypt.ErrCiphertextTooShort) {
		t.Fatalf("Decrypt(short) error = %v, want ErrCiphertextTooShort", err)
	}
}

func TestHashStringToStringUsesFcryptSHA3(t *testing.T) {
	got := HashStringToString("hello")
	want := fcrypt.HashStringToStringSHA3("hello")
	if got != want {
		t.Fatalf("HashStringToString() = %s, want %s", got, want)
	}
}
