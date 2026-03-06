package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

// RandomBytes returns n cryptographically random bytes.
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, err
	}
	return b, nil
}

// RandomHex returns a random hex string of n bytes (2n hex chars).
func RandomHex(n int) string {
	b, _ := RandomBytes(n)
	return hex.EncodeToString(b)
}

// RandomBase64 returns a URL-safe base64-encoded random string.
func RandomBase64(n int) string {
	b, _ := RandomBytes(n)
	return base64.URLEncoding.EncodeToString(b)
}

// SHA256Hash returns the hex-encoded SHA-256 hash of data.
func SHA256Hash(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

// HMACSign creates an HMAC-SHA256 signature.
func HMACSign(message, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// HMACVerify checks an HMAC-SHA256 signature (constant-time).
func HMACVerify(message, signature, key string) bool {
	expected := HMACSign(message, key)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// Encrypt encrypts plaintext using AES-256-GCM. Key must be 32 bytes.
func Encrypt(plaintext string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be exactly 32 bytes for AES-256")
	}
	block, err := aes.NewCipher(key)
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
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts AES-256-GCM ciphertext. Key must be 32 bytes.
func Decrypt(encrypted string, key []byte) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key must be exactly 32 bytes for AES-256")
	}
	data, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	plaintext, err := gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
