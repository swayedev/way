package way

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	ErrSessionStoreNotFound = errors.New("session store not found")
	ErrSecureCookieNotFound = errors.New("secure cookie not found")
)

// defaultStore is the name of the default session store.
// defaultCookie is the name of the default session cookie.
// stores is a map of session stores.
// cookies is a map of secure cookies.
type Session struct {
	// Name of
	defaultStore string
	// Name of secure cookie
	defaultCookie string
	// Map of session stores
	stores map[string]sessions.Store
	// Map of secure cookies
	cookies map[string]*securecookie.SecureCookie
}

func NewSession() *Session {
	return &Session{
		defaultStore:  "default",
		defaultCookie: "default",
		stores:        make(map[string]sessions.Store),
		cookies:       make(map[string]*securecookie.SecureCookie),
	}
}

func (w *Session) SetDefaultStoreName(name string) {
	w.defaultStore = name
}

func (w *Session) SetDefaultCookieName(name string) {
	w.defaultCookie = name
}

func (w *Session) Stores() map[string]sessions.Store {
	return w.stores
}

func (w *Session) Store(name string) sessions.Store {
	if w == nil {
		return nil
	}
	return w.stores[name]
}

func (w *Session) StoreE(name string) (sessions.Store, error) {
	if w == nil {
		return nil, fmt.Errorf("%w: session manager is nil", ErrSessionStoreNotFound)
	}
	store := w.stores[name]
	if store == nil {
		return nil, fmt.Errorf("%w: %s", ErrSessionStoreNotFound, name)
	}
	return store, nil
}

func (w *Session) SetStore(name string, s sessions.Store) {
	w.stores[name] = s
}

func (w *Session) DeleteStore(name string) {
	delete(w.stores, name)
}

func (w *Session) Cookies() map[string]*securecookie.SecureCookie {
	return w.cookies
}

func (w *Session) Cookie(name string) *securecookie.SecureCookie {
	if w == nil {
		return nil
	}
	return w.cookies[name]
}

func (w *Session) CookieE(name string) (*securecookie.SecureCookie, error) {
	if w == nil {
		return nil, fmt.Errorf("%w: session manager is nil", ErrSecureCookieNotFound)
	}
	cookie := w.cookies[name]
	if cookie == nil {
		return nil, fmt.Errorf("%w: %s", ErrSecureCookieNotFound, name)
	}
	return cookie, nil
}

func (w *Session) SetCookie(name string, s *securecookie.SecureCookie) {
	w.cookies[name] = s
}

func (w *Session) DeleteCookie(name string) {
	delete(w.cookies, name)
}

func (w *Session) DefaultSession() sessions.Store {
	if w == nil {
		return nil
	}
	return w.stores[w.defaultStore]
}

func (w *Session) DefaultSessionE() (sessions.Store, error) {
	if w == nil {
		return nil, fmt.Errorf("%w: session manager is nil", ErrSessionStoreNotFound)
	}
	return w.StoreE(w.defaultStore)
}

func (w *Session) SetDefaultStore(s sessions.Store) {
	w.stores[w.defaultStore] = s
}

func (w *Session) DefaultCookie() *securecookie.SecureCookie {
	if w == nil {
		return nil
	}
	return w.cookies[w.defaultCookie]
}

func (w *Session) DefaultCookieE() (*securecookie.SecureCookie, error) {
	if w == nil {
		return nil, fmt.Errorf("%w: session manager is nil", ErrSecureCookieNotFound)
	}
	return w.CookieE(w.defaultCookie)
}

func (w *Session) SetDefaultCookie(s *securecookie.SecureCookie) {
	w.cookies[w.defaultCookie] = s
}

func (w *Session) CreateEncryptedCookie(
	wr http.ResponseWriter,
	name string,
	cookieName string,
	value map[string]interface{},
	path string,
	maxAge int,
	httpOnly bool,
	secure bool) (*http.Cookie, error) {
	secureCookie, err := w.CookieE(name)
	if err != nil {
		return nil, err
	}
	return CreateEncryptedCookie(wr, *secureCookie, cookieName, value, path, maxAge, httpOnly, secure)
}

func (w *Session) CreateEncryptedCookieWithDefaults(
	wr http.ResponseWriter,
	name string,
	cookieName string,
	value map[string]interface{}) (*http.Cookie, error) {
	secureCookie, err := w.CookieE(name)
	if err != nil {
		return nil, err
	}
	return CreateEncryptedCookieWithDefaults(wr, *secureCookie, cookieName, value)
}

func (w *Session) CreateDefaultEncryptedCookie(
	wr http.ResponseWriter,
	cookieName string,
	value map[string]interface{},
	path string,
	maxAge int,
	httpOnly bool,
	secure bool) (*http.Cookie, error) {
	secureCookie, err := w.DefaultCookieE()
	if err != nil {
		return nil, err
	}
	return CreateEncryptedCookie(wr, *secureCookie, cookieName, value, path, maxAge, httpOnly, secure)
}

func (w *Session) CreateDefaultEncryptedCookieWithDefaults(
	wr http.ResponseWriter,
	cookieName string,
	value map[string]interface{}) (*http.Cookie, error) {
	secureCookie, err := w.DefaultCookieE()
	if err != nil {
		return nil, err
	}
	return CreateEncryptedCookieWithDefaults(wr, *secureCookie, cookieName, value)
}

func (w *Session) ReadEncryptedCookie(r *http.Request, name string, cookieName string) (map[string]string, error) {
	secureCookie, err := w.CookieE(name)
	if err != nil {
		return nil, err
	}
	return ReadEncryptedCookie(r, *secureCookie, cookieName)
}

func (w *Session) ReadDefaultEncryptedCookie(r *http.Request, name string, cookieName string) (map[string]string, error) {
	secureCookie, err := w.DefaultCookieE()
	if err != nil {
		return nil, err
	}
	return ReadEncryptedCookie(r, *secureCookie, cookieName)
}

func CreateEncryptedCookie(
	w http.ResponseWriter,
	secureCookie securecookie.SecureCookie,
	name string,
	value map[string]interface{},
	path string,
	maxAge int,
	httpOnly bool,
	secure bool) (*http.Cookie, error) {
	encoded, err := secureCookie.Encode(name, value)
	if err != nil {
		return nil, err
	}
	cookie := &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     path,
		MaxAge:   maxAge,
		HttpOnly: httpOnly,
		Secure:   secure,
	}
	return cookie, nil
}

func CreateEncryptedCookieWithDefaults(
	w http.ResponseWriter,
	secureCookie securecookie.SecureCookie,
	name string,
	value map[string]interface{}) (*http.Cookie, error) {
	return CreateEncryptedCookie(w, secureCookie, name, value, "/", 36000, true, true)
}

func ReadEncryptedCookie(r *http.Request, secureCookie securecookie.SecureCookie, name string) (map[string]string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return nil, err
	}

	var value map[string]string
	if err = secureCookie.Decode(name, cookie.Value, &value); err == nil {
		return value, nil
	}

	return nil, err
}
