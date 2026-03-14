package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Manager ties together the store, cookie handling, and ID generation.
type Manager struct {
	Store      Store
	CookieName string
	Lifetime   time.Duration
	Path       string
	Domain     string
	Secure     bool
	HTTPOnly   bool
	SameSite   http.SameSite
}

// NewManager creates a session manager from environment config.
func NewManager(store Store) *Manager {
	lifetime, _ := strconv.Atoi(os.Getenv("SESSION_LIFETIME"))
	if lifetime == 0 {
		lifetime = 120
	}
	secure := os.Getenv("SESSION_SECURE") == "true"
	httpOnly := os.Getenv("SESSION_HTTPONLY") != "false"

	sameSite := http.SameSiteLaxMode
	switch os.Getenv("SESSION_SAMESITE") {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	}

	cookieName := os.Getenv("SESSION_COOKIE")
	if cookieName == "" {
		cookieName = "framework_session"
	}

	return &Manager{
		Store:      store,
		CookieName: cookieName,
		Lifetime:   time.Duration(lifetime) * time.Minute,
		Path:       os.Getenv("SESSION_PATH"),
		Domain:     os.Getenv("SESSION_DOMAIN"),
		Secure:     secure,
		HTTPOnly:   httpOnly,
		SameSite:   sameSite,
	}
}

func (m *Manager) generateID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Start retrieves or creates a session for the request.
func (m *Manager) Start(r *http.Request) (string, map[string]interface{}, error) {
	cookie, err := r.Cookie(m.CookieName)
	if err == nil && cookie.Value != "" {
		data, err := m.Store.Read(cookie.Value)
		if err != nil {
			return "", nil, err
		}
		if len(data) > 0 {
			return cookie.Value, data, nil
		}
	}
	id := m.generateID()
	return id, make(map[string]interface{}), nil
}

// SetCookie writes the session cookie to the response without persisting data.
// Use this to set the cookie before the response body is written.
func (m *Manager) SetCookie(w http.ResponseWriter, id string) {
	http.SetCookie(w, &http.Cookie{
		Name:     m.CookieName,
		Value:    id,
		Path:     m.Path,
		Domain:   m.Domain,
		MaxAge:   int(m.Lifetime.Seconds()),
		Secure:   m.Secure,
		HttpOnly: m.HTTPOnly,
		SameSite: m.SameSite,
	})
}

// Save persists session data and writes the session cookie.
func (m *Manager) Save(w http.ResponseWriter, id string, data map[string]interface{}) error {
	if err := m.Store.Write(id, data, m.Lifetime); err != nil {
		return err
	}
	m.SetCookie(w, id)
	return nil
}

// Destroy removes a session and clears the cookie.
func (m *Manager) Destroy(w http.ResponseWriter, id string) error {
	if err := m.Store.Destroy(id); err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:   m.CookieName,
		Value:  "",
		Path:   m.Path,
		Domain: m.Domain,
		MaxAge: -1,
	})
	return nil
}

// Flash writes a flash message that will be available on the next request only.
func (m *Manager) Flash(data map[string]interface{}, key string, value interface{}) {
	flashes, _ := data["_flashes"].(map[string]interface{})
	if flashes == nil {
		flashes = make(map[string]interface{})
	}
	flashes[key] = value
	data["_flashes"] = flashes
}

// GetFlash reads and removes a flash message.
func (m *Manager) GetFlash(data map[string]interface{}, key string) (interface{}, bool) {
	flashes, _ := data["_flashes"].(map[string]interface{})
	if flashes == nil {
		return nil, false
	}
	val, ok := flashes[key]
	if ok {
		delete(flashes, key)
		if len(flashes) == 0 {
			delete(data, "_flashes")
		} else {
			data["_flashes"] = flashes
		}
	}
	return val, ok
}

// FlashErrors stores validation errors for the next request.
func (m *Manager) FlashErrors(data map[string]interface{}, errors map[string][]string) {
	m.Flash(data, "_errors", errors)
}

// FlashOldInput stores form input for re-populating after validation failure.
func (m *Manager) FlashOldInput(data map[string]interface{}, input map[string]string) {
	m.Flash(data, "_old_input", input)
}
