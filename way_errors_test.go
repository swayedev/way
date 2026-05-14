package way

import (
	"os"
	"testing"
)

func TestGetEncryptionKeyReturnsErrorWhenNotSet(t *testing.T) {
	// Save and clear the env var
	oldVal, wasSet := os.LookupEnv(envCookieEncryptionKey)
	os.Unsetenv(envCookieEncryptionKey)
	defer func() {
		if wasSet {
			os.Setenv(envCookieEncryptionKey, oldVal)
		}
	}()

	key, err := getEncryptionKey()
	if err == nil {
		t.Error("getEncryptionKey() expected error when env var not set, got nil")
	}
	if key != nil {
		t.Errorf("getEncryptionKey() expected nil key on error, got %v", key)
	}
}

func TestGetEncryptionKeyReturnsKeyWhenSet(t *testing.T) {
	testKey := "test-encryption-key-12345"
	os.Setenv(envCookieEncryptionKey, testKey)
	defer os.Unsetenv(envCookieEncryptionKey)

	key, err := getEncryptionKey()
	if err != nil {
		t.Errorf("getEncryptionKey() unexpected error: %v", err)
	}
	if string(key) != testKey {
		t.Errorf("getEncryptionKey() = %s, want %s", string(key), testKey)
	}
}

func TestGetAuthenticationKeyReturnsErrorWhenNotSet(t *testing.T) {
	oldVal, wasSet := os.LookupEnv(envCookieAuthenticationKey)
	os.Unsetenv(envCookieAuthenticationKey)
	defer func() {
		if wasSet {
			os.Setenv(envCookieAuthenticationKey, oldVal)
		}
	}()

	key, err := getAuthenticationKey()
	if err == nil {
		t.Error("getAuthenticationKey() expected error when env var not set, got nil")
	}
	if key != nil {
		t.Errorf("getAuthenticationKey() expected nil key on error, got %v", key)
	}
}

func TestGetAuthenticationKeyReturnsKeyWhenSet(t *testing.T) {
	testKey := "test-auth-key-12345"
	os.Setenv(envCookieAuthenticationKey, testKey)
	defer os.Unsetenv(envCookieAuthenticationKey)

	key, err := getAuthenticationKey()
	if err != nil {
		t.Errorf("getAuthenticationKey() unexpected error: %v", err)
	}
	if string(key) != testKey {
		t.Errorf("getAuthenticationKey() = %s, want %s", string(key), testKey)
	}
}

func TestGetStoreEncryptionKeyReturnsErrorWhenNotSet(t *testing.T) {
	oldVal, wasSet := os.LookupEnv(envStoreEncryptionKey)
	os.Unsetenv(envStoreEncryptionKey)
	defer func() {
		if wasSet {
			os.Setenv(envStoreEncryptionKey, oldVal)
		}
	}()

	key, err := getStoreEncryptionKey()
	if err == nil {
		t.Error("getStoreEncryptionKey() expected error when env var not set, got nil")
	}
	if key != nil {
		t.Errorf("getStoreEncryptionKey() expected nil key on error, got %v", key)
	}
}

func TestGetStoreEncryptionKeyReturnsKeyWhenSet(t *testing.T) {
	testKey := "test-store-key-12345"
	os.Setenv(envStoreEncryptionKey, testKey)
	defer os.Unsetenv(envStoreEncryptionKey)

	key, err := getStoreEncryptionKey()
	if err != nil {
		t.Errorf("getStoreEncryptionKey() unexpected error: %v", err)
	}
	if string(key) != testKey {
		t.Errorf("getStoreEncryptionKey() = %s, want %s", string(key), testKey)
	}
}
