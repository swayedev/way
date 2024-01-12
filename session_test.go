package way

import (
	"testing"
)

func TestNewSession(t *testing.T) {
	// Call the NewSession function
	session := NewSession()

	// Check if the defaultStore, defaultCookie, stores, and cookies fields are set correctly
	if session.defaultStore != "default" {
		t.Errorf("NewSession() defaultStore field = %v; want %v", session.defaultStore, "default")
	}
	if session.defaultCookie != "default" {
		t.Errorf("NewSession() defaultCookie field = %v; want %v", session.defaultCookie, "default")
	}
	if session.stores == nil {
		t.Errorf("NewSession() stores field is nil; want a non-nil map")
	}
	if session.cookies == nil {
		t.Errorf("NewSession() cookies field is nil; want a non-nil map")
	}
}

func TestSetDefaultStoreName(t *testing.T) {
	// Create a new session
	session := NewSession()

	// Set the default store name
	session.SetDefaultStoreName("custom")

	// Check if the defaultStore field is set correctly
	if session.defaultStore != "custom" {
		t.Errorf("SetDefaultStoreName() defaultStore field = %v; want %v", session.defaultStore, "custom")
	}
}
