package crypto

import (
	"testing"
)

func TestGenerateRandomKeyReturnsKeyWhenSuccessful(t *testing.T) {
	length := 32
	key, err := GenerateRandomKey(length)
	if err != nil {
		t.Errorf("GenerateRandomKey() unexpected error: %v", err)
	}
	if len(key) != length {
		t.Errorf("GenerateRandomKey() returned key of length %d, want %d", len(key), length)
	}
	if key == nil {
		t.Error("GenerateRandomKey() returned nil key")
	}
}

func TestGenerateRandomKeyReturnsDifferentKeysOnMultipleCalls(t *testing.T) {
	key1, err1 := GenerateRandomKey(16)
	key2, err2 := GenerateRandomKey(16)
	if err1 != nil || err2 != nil {
		t.Errorf("GenerateRandomKey() unexpected error: %v, %v", err1, err2)
	}
	if len(key1) != len(key2) {
		t.Errorf("GenerateRandomKey() returned different lengths: %d vs %d", len(key1), len(key2))
	}
	// Check that keys are different (with extremely high probability)
	same := true
	for i := range key1 {
		if key1[i] != key2[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("GenerateRandomKey() returned the same key on multiple calls (extremely unlikely)")
	}
}

func TestGenerateRandomKeyZeroLength(t *testing.T) {
	key, err := GenerateRandomKey(0)
	if err != nil {
		t.Errorf("GenerateRandomKey(0) unexpected error: %v", err)
	}
	if len(key) != 0 {
		t.Errorf("GenerateRandomKey(0) returned key of length %d, want 0", len(key))
	}
}

func TestGenerateRandomKeyLargeLength(t *testing.T) {
	length := 65536
	key, err := GenerateRandomKey(length)
	if err != nil {
		t.Errorf("GenerateRandomKey(%d) unexpected error: %v", length, err)
	}
	if len(key) != length {
		t.Errorf("GenerateRandomKey(%d) returned key of length %d", length, len(key))
	}
}
