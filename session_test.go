package way

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/securecookie"
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

func TestSessionMissingStoreReturnsError(t *testing.T) {
	session := NewSession()

	store, err := session.StoreE("missing")
	if !errors.Is(err, ErrSessionStoreNotFound) {
		t.Fatalf("error = %v, want ErrSessionStoreNotFound", err)
	}
	if store != nil {
		t.Fatalf("store = %v, want nil", store)
	}
}

func TestSessionMissingSecureCookieReturnsError(t *testing.T) {
	session := NewSession()

	cookie, err := session.CookieE("missing")
	if !errors.Is(err, ErrSecureCookieNotFound) {
		t.Fatalf("error = %v, want ErrSecureCookieNotFound", err)
	}
	if cookie != nil {
		t.Fatalf("cookie = %v, want nil", cookie)
	}
}

func TestCreateEncryptedCookieMissingCookieDoesNotPanic(t *testing.T) {
	session := NewSession()
	rec := httptest.NewRecorder()

	cookie, err := session.CreateEncryptedCookieWithDefaults(rec, "missing", "session", map[string]interface{}{"user": "1"})
	if !errors.Is(err, ErrSecureCookieNotFound) {
		t.Fatalf("error = %v, want ErrSecureCookieNotFound", err)
	}
	if cookie != nil {
		t.Fatalf("cookie = %v, want nil", cookie)
	}
}

func TestCreateEncryptedCookieWithConfiguredCookie(t *testing.T) {
	session := NewSession()
	session.SetCookie("secure", securecookie.New([]byte("01234567890123456789012345678901"), []byte("01234567890123456789012345678901")))
	rec := httptest.NewRecorder()

	cookie, err := session.CreateEncryptedCookieWithDefaults(rec, "secure", "session", map[string]interface{}{"user": "1"})
	if err != nil {
		t.Fatalf("CreateEncryptedCookieWithDefaults() error = %v", err)
	}
	if cookie == nil || cookie.Name != "session" {
		t.Fatalf("cookie = %#v, want named session cookie", cookie)
	}
}
