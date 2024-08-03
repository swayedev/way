package way

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// Session manages session stores and secure cookies.
type Session struct {
	defaultStore  string
	defaultCookie string
	stores        map[string]sessions.Store
	cookies       map[string]*securecookie.SecureCookie
	logger        *log.Logger
}

// NewSession initializes a new Session instance.
func NewSession(logger *log.Logger) *Session {
	return &Session{
		defaultStore:  "default",
		defaultCookie: "default",
		stores:        make(map[string]sessions.Store),
		cookies:       make(map[string]*securecookie.SecureCookie),
		logger:        logger,
	}
}

// SetDefaultStoreName sets the default session store name.
func (s *Session) SetDefaultStoreName(name string) {
	s.defaultStore = name
}

// SetDefaultCookieName sets the default cookie name.
func (s *Session) SetDefaultCookieName(name string) {
	s.defaultCookie = name
}

// Stores returns the map of session stores.
func (s *Session) Stores() map[string]sessions.Store {
	return s.stores
}

// Store returns a session store by name.
func (s *Session) Store(name string) sessions.Store {
	return s.stores[name]
}

// SetStore sets a session store by name.
func (s *Session) SetStore(name string, store sessions.Store) {
	s.stores[name] = store
}

// DeleteStore deletes a session store by name.
func (s *Session) DeleteStore(name string) {
	delete(s.stores, name)
}

// Cookies returns the map of secure cookies.
func (s *Session) Cookies() map[string]*securecookie.SecureCookie {
	return s.cookies
}

// Cookie returns a secure cookie by name.
func (s *Session) Cookie(name string) *securecookie.SecureCookie {
	return s.cookies[name]
}

// SetCookie sets a secure cookie by name.
func (s *Session) SetCookie(name string, cookie *securecookie.SecureCookie) {
	s.cookies[name] = cookie
}

// DeleteCookie deletes a secure cookie by name.
func (s *Session) DeleteCookie(name string) {
	delete(s.cookies, name)
}

// DefaultSession returns the default session store.
func (s *Session) DefaultSession() sessions.Store {
	return s.stores[s.defaultStore]
}

// SetDefaultStore sets the default session store.
func (s *Session) SetDefaultStore(store sessions.Store) {
	s.stores[s.defaultStore] = store
}

// DefaultCookie returns the default secure cookie.
func (s *Session) DefaultCookie() *securecookie.SecureCookie {
	return s.cookies[s.defaultCookie]
}

// SetDefaultCookie sets the default secure cookie.
func (s *Session) SetDefaultCookie(cookie *securecookie.SecureCookie) {
	s.cookies[s.defaultCookie] = cookie
}

// CreateEncryptedCookie creates an encrypted cookie.
func (s *Session) CreateEncryptedCookie(
	w http.ResponseWriter,
	name, cookieName string,
	value map[string]interface{},
	path string,
	maxAge int,
	httpOnly, secure bool) (*http.Cookie, error) {
	return createEncryptedCookie(w, *s.cookies[name], cookieName, value, path, maxAge, httpOnly, secure)
}

// CreateEncryptedCookieWithDefaults creates an encrypted cookie with default settings.
func (s *Session) CreateEncryptedCookieWithDefaults(
	w http.ResponseWriter,
	name, cookieName string,
	value map[string]interface{}) (*http.Cookie, error) {
	return createEncryptedCookieWithDefaults(w, *s.cookies[name], cookieName, value)
}

// CreateDefaultEncryptedCookie creates an encrypted cookie using the default secure cookie.
func (s *Session) CreateDefaultEncryptedCookie(
	w http.ResponseWriter,
	cookieName string,
	value map[string]interface{},
	path string,
	maxAge int,
	httpOnly, secure bool) (*http.Cookie, error) {
	return createEncryptedCookie(w, *s.cookies[s.defaultCookie], cookieName, value, path, maxAge, httpOnly, secure)
}

// CreateDefaultEncryptedCookieWithDefaults creates an encrypted cookie with default settings using the default secure cookie.
func (s *Session) CreateDefaultEncryptedCookieWithDefaults(
	w http.ResponseWriter,
	cookieName string,
	value map[string]interface{}) (*http.Cookie, error) {
	return createEncryptedCookieWithDefaults(w, *s.cookies[s.defaultCookie], cookieName, value)
}

// ReadEncryptedCookie reads an encrypted cookie.
func (s *Session) ReadEncryptedCookie(r *http.Request, name, cookieName string) (map[string]string, error) {
	return readEncryptedCookie(r, *s.cookies[name], cookieName)
}

// ReadDefaultEncryptedCookie reads an encrypted cookie using the default secure cookie.
func (s *Session) ReadDefaultEncryptedCookie(r *http.Request, cookieName string) (map[string]string, error) {
	return readEncryptedCookie(r, *s.cookies[s.defaultCookie], cookieName)
}

// createEncryptedCookie creates an encrypted cookie.
func createEncryptedCookie(
	w http.ResponseWriter,
	secureCookie securecookie.SecureCookie,
	name string,
	value map[string]interface{},
	path string,
	maxAge int,
	httpOnly, secure bool) (*http.Cookie, error) {
	encoded, err := secureCookie.Encode(name, value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode cookie %s: %w", name, err)
	}
	cookie := &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     path,
		MaxAge:   maxAge,
		HttpOnly: httpOnly,
		Secure:   secure,
	}
	http.SetCookie(w, cookie)
	return cookie, nil
}

// createEncryptedCookieWithDefaults creates an encrypted cookie with default settings.
func createEncryptedCookieWithDefaults(
	w http.ResponseWriter,
	secureCookie securecookie.SecureCookie,
	name string,
	value map[string]interface{}) (*http.Cookie, error) {
	return createEncryptedCookie(w, secureCookie, name, value, "/", 36000, true, true)
}

// readEncryptedCookie reads an encrypted cookie.
func readEncryptedCookie(r *http.Request, secureCookie securecookie.SecureCookie, name string) (map[string]string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read cookie %s: %w", name, err)
	}

	var value map[string]string
	if err = secureCookie.Decode(name, cookie.Value, &value); err != nil {
		return nil, fmt.Errorf("failed to decode cookie %s: %w", name, err)
	}

	return value, nil
}
