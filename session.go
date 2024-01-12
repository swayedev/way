package way

import (
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
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
	return w.stores[name]
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
	return w.cookies[name]
}

func (w *Session) SetCookie(name string, s *securecookie.SecureCookie) {
	w.cookies[name] = s
}

func (w *Session) DeleteCookie(name string) {
	delete(w.cookies, name)
}

func (w *Session) DefaultSession() sessions.Store {
	return w.stores[w.defaultStore]
}

func (w *Session) SetDefaultStore(s sessions.Store) {
	w.stores[w.defaultStore] = s
}

func (w *Session) DefaultCookie() *securecookie.SecureCookie {
	return w.cookies[w.defaultCookie]
}

func (w *Session) SetDefaultCookie(s *securecookie.SecureCookie) {
	w.cookies[w.defaultCookie] = s
}

func (w *Session) CreateEncryptedCookie(
	wr http.ResponseWriter,
	name string,
	value map[string]interface{},
	path string,
	maxAge int,
	httpOnly bool,
	secure bool) (*http.Cookie, error) {
	return CreateEncryptedCookie(wr, *w.cookies[name], name, value, path, maxAge, httpOnly, secure)
}

func (w *Session) ReadEncryptedCookie(r *http.Request, name string) (map[string]string, error) {
	return ReadEncryptedCookie(r, *w.cookies[name], name)
}

func (w *Session) CreateEncryptedCookieWithDefaults(
	wr http.ResponseWriter,
	name string,
	value map[string]interface{}) (*http.Cookie, error) {
	return CreateEncryptedCookieWithDefaults(wr, *w.cookies[name], name, value)
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
