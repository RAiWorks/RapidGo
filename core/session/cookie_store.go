package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"sync"
	"time"
)

// CookieStore stores session data encrypted in the cookie value itself.
// Requires a 32-byte key (AES-256-GCM).
type CookieStore struct {
	Key  []byte
	mu   sync.RWMutex
	data map[string]string // id → encrypted payload
}

// NewCookieStore creates a cookie-based session store with the given AES-256 key.
func NewCookieStore(key []byte) (*CookieStore, error) {
	if len(key) != 32 {
		return nil, errors.New("cookie store requires a 32-byte key")
	}
	return &CookieStore{Key: key, data: make(map[string]string)}, nil
}

func (s *CookieStore) encrypt(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *CookieStore) decrypt(encoded string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	return gcm.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
}

func (s *CookieStore) Read(id string) (map[string]interface{}, error) {
	s.mu.RLock()
	encoded, ok := s.data[id]
	s.mu.RUnlock()
	if !ok {
		return make(map[string]interface{}), nil
	}
	plaintext, err := s.decrypt(encoded)
	if err != nil {
		return make(map[string]interface{}), nil
	}
	var data map[string]interface{}
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *CookieStore) Write(id string, data map[string]interface{}, lifetime time.Duration) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	encrypted, err := s.encrypt(raw)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.data[id] = encrypted
	s.mu.Unlock()
	return nil
}

func (s *CookieStore) Destroy(id string) error {
	s.mu.Lock()
	delete(s.data, id)
	s.mu.Unlock()
	return nil
}

func (s *CookieStore) GC(maxLifetime time.Duration) error {
	return nil // Expiry handled by cookie MaxAge
}
