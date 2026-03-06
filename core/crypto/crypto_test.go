package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

// TC-01: RandomBytes returns n bytes
func TestRandomBytesLength(t *testing.T) {
	b, err := RandomBytes(16)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b) != 16 {
		t.Fatalf("expected 16 bytes, got %d", len(b))
	}
}

// TC-02: RandomBytes returns unique values
func TestRandomBytesUnique(t *testing.T) {
	a, _ := RandomBytes(16)
	b, _ := RandomBytes(16)
	if string(a) == string(b) {
		t.Fatal("expected different byte slices")
	}
}

// TC-03: RandomHex correct length
func TestRandomHexLength(t *testing.T) {
	h := RandomHex(16)
	if len(h) != 32 {
		t.Fatalf("expected 32-char hex, got %d", len(h))
	}
	if _, err := hex.DecodeString(h); err != nil {
		t.Fatalf("invalid hex: %v", err)
	}
}

// TC-04: RandomBase64 non-empty and valid
func TestRandomBase64Valid(t *testing.T) {
	s := RandomBase64(24)
	if s == "" {
		t.Fatal("expected non-empty string")
	}
	if _, err := base64.URLEncoding.DecodeString(s); err != nil {
		t.Fatalf("invalid base64: %v", err)
	}
}

// TC-05: SHA256Hash deterministic
func TestSHA256HashDeterministic(t *testing.T) {
	a := SHA256Hash("hello")
	b := SHA256Hash("hello")
	if a != b {
		t.Fatal("expected identical hashes")
	}
}

// TC-06: SHA256Hash length
func TestSHA256HashLength(t *testing.T) {
	h := SHA256Hash("test")
	if len(h) != 64 {
		t.Fatalf("expected 64-char hex, got %d", len(h))
	}
}

// TC-07: HMACSign deterministic
func TestHMACSignDeterministic(t *testing.T) {
	a := HMACSign("msg", "key")
	b := HMACSign("msg", "key")
	if a != b {
		t.Fatal("expected identical signatures")
	}
}

// TC-08: HMACVerify valid
func TestHMACVerifyValid(t *testing.T) {
	sig := HMACSign("msg", "key")
	if !HMACVerify("msg", sig, "key") {
		t.Fatal("expected verify to return true")
	}
}

// TC-09: HMACVerify tampered message
func TestHMACVerifyTampered(t *testing.T) {
	sig := HMACSign("a", "key")
	if HMACVerify("b", sig, "key") {
		t.Fatal("expected verify to return false for tampered message")
	}
}

// TC-10: HMACVerify wrong key
func TestHMACVerifyWrongKey(t *testing.T) {
	sig := HMACSign("msg", "key1")
	if HMACVerify("msg", sig, "key2") {
		t.Fatal("expected verify to return false for wrong key")
	}
}

// TC-11: Encrypt/Decrypt round-trip
func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := make([]byte, 32)
	copy(key, []byte("12345678901234567890123456789012"))
	plain := "secret data"
	enc, err := Encrypt(plain, key)
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}
	dec, err := Decrypt(enc, key)
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	if dec != plain {
		t.Fatalf("expected %q, got %q", plain, dec)
	}
}

// TC-12: Encrypt produces unique output (random nonce)
func TestEncryptUnique(t *testing.T) {
	key := make([]byte, 32)
	copy(key, []byte("12345678901234567890123456789012"))
	a, _ := Encrypt("same", key)
	b, _ := Encrypt("same", key)
	if a == b {
		t.Fatal("expected different ciphertext due to random nonce")
	}
}

// TC-13: Decrypt fails wrong key
func TestDecryptWrongKey(t *testing.T) {
	key1 := make([]byte, 32)
	copy(key1, []byte("12345678901234567890123456789012"))
	key2 := make([]byte, 32)
	copy(key2, []byte("abcdefghijklmnopqrstuvwxyz123456"))
	enc, _ := Encrypt("secret", key1)
	_, err := Decrypt(enc, key2)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

// TC-14: Decrypt fails short ciphertext
func TestDecryptShortCiphertext(t *testing.T) {
	key := make([]byte, 32)
	copy(key, []byte("12345678901234567890123456789012"))
	short := base64.URLEncoding.EncodeToString([]byte("abc"))
	_, err := Decrypt(short, key)
	if err == nil {
		t.Fatal("expected error for short ciphertext")
	}
	if !strings.Contains(err.Error(), "ciphertext too short") {
		t.Fatalf("expected 'ciphertext too short', got: %v", err)
	}
}

// TC-15: Encrypt fails invalid key length
func TestEncryptInvalidKey(t *testing.T) {
	key := make([]byte, 16)
	_, err := Encrypt("data", key)
	if err == nil {
		t.Fatal("expected error for 16-byte key")
	}
}
