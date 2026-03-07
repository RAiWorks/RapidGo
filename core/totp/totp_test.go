package totp

import (
	"encoding/base32"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

// ── T01: GenerateKey returns a valid Key ──────────────────────────────────────

func TestGenerateKey_ReturnsKey(t *testing.T) {
	key, err := GenerateKey("RapidGo", "user@example.com")
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}
	if key == nil {
		t.Fatal("GenerateKey() returned nil key")
	}
	if key.Secret == "" {
		t.Error("Key.Secret is empty")
	}
	if key.URL == "" {
		t.Error("Key.URL is empty")
	}
}

// ── T02: URL contains the issuer ─────────────────────────────────────────────

func TestGenerateKey_URLContainsIssuer(t *testing.T) {
	key, err := GenerateKey("MyApp", "user@example.com")
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}
	if !strings.Contains(key.URL, "MyApp") {
		t.Errorf("URL %q does not contain issuer 'MyApp'", key.URL)
	}
}

// ── T03: URL contains the account ────────────────────────────────────────────

func TestGenerateKey_URLContainsAccount(t *testing.T) {
	key, err := GenerateKey("RapidGo", "alice@example.com")
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}
	if !strings.Contains(key.URL, "alice") {
		t.Errorf("URL %q does not contain account 'alice'", key.URL)
	}
}

// ── T04: Secret is valid Base32 ──────────────────────────────────────────────

func TestGenerateKey_SecretIsBase32(t *testing.T) {
	key, err := GenerateKey("RapidGo", "user@example.com")
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}
	_, err = base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(key.Secret)
	if err != nil {
		t.Errorf("Secret %q is not valid Base32: %v", key.Secret, err)
	}
}

// ── T05: Two calls produce unique secrets ────────────────────────────────────

func TestGenerateKey_UniqueSecrets(t *testing.T) {
	k1, err := GenerateKey("RapidGo", "user@example.com")
	if err != nil {
		t.Fatalf("GenerateKey() #1 error: %v", err)
	}
	k2, err := GenerateKey("RapidGo", "user@example.com")
	if err != nil {
		t.Fatalf("GenerateKey() #2 error: %v", err)
	}
	if k1.Secret == k2.Secret {
		t.Error("two consecutive GenerateKey() calls returned the same secret")
	}
}

// ── T06: ValidateCode with a valid code ──────────────────────────────────────

func TestValidateCode_ValidCode(t *testing.T) {
	key, err := GenerateKey("RapidGo", "user@example.com")
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}
	code, err := totp.GenerateCode(key.Secret, time.Now())
	if err != nil {
		t.Fatalf("totp.GenerateCode() error: %v", err)
	}
	if !ValidateCode(key.Secret, code) {
		t.Errorf("ValidateCode() returned false for valid code %q", code)
	}
}

// ── T07: ValidateCode rejects invalid code ───────────────────────────────────

func TestValidateCode_InvalidCode(t *testing.T) {
	key, err := GenerateKey("RapidGo", "user@example.com")
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}
	if ValidateCode(key.Secret, "000000") {
		t.Error("ValidateCode() returned true for invalid code '000000'")
	}
}

// ── T08: ValidateCode rejects code for wrong secret ──────────────────────────

func TestValidateCode_WrongSecret(t *testing.T) {
	k1, err := GenerateKey("RapidGo", "user1@example.com")
	if err != nil {
		t.Fatalf("GenerateKey() #1 error: %v", err)
	}
	k2, err := GenerateKey("RapidGo", "user2@example.com")
	if err != nil {
		t.Fatalf("GenerateKey() #2 error: %v", err)
	}
	code, err := totp.GenerateCode(k1.Secret, time.Now())
	if err != nil {
		t.Fatalf("totp.GenerateCode() error: %v", err)
	}
	if ValidateCode(k2.Secret, code) {
		t.Error("ValidateCode() returned true for code validated against wrong secret")
	}
}

// ── T09: ValidateCode rejects empty code ─────────────────────────────────────

func TestValidateCode_EmptyCode(t *testing.T) {
	key, err := GenerateKey("RapidGo", "user@example.com")
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}
	if ValidateCode(key.Secret, "") {
		t.Error("ValidateCode() returned true for empty code")
	}
}

// ── T10: GenerateBackupCodes returns correct count ───────────────────────────

func TestGenerateBackupCodes_ReturnsCorrectCount(t *testing.T) {
	codes, err := GenerateBackupCodes(8)
	if err != nil {
		t.Fatalf("GenerateBackupCodes() error: %v", err)
	}
	if len(codes) != 8 {
		t.Errorf("expected 8 codes, got %d", len(codes))
	}
}

// ── T11: Backup codes match XXXX-XXXX format ────────────────────────────────

func TestGenerateBackupCodes_FormatXXXXDashXXXX(t *testing.T) {
	codes, err := GenerateBackupCodes(10)
	if err != nil {
		t.Fatalf("GenerateBackupCodes() error: %v", err)
	}
	pattern := regexp.MustCompile(`^[0-9A-F]{4}-[0-9A-F]{4}$`)
	for _, c := range codes {
		if !pattern.MatchString(c) {
			t.Errorf("code %q does not match XXXX-XXXX hex format", c)
		}
	}
}

// ── T12: All backup codes are unique ─────────────────────────────────────────

func TestGenerateBackupCodes_AllUnique(t *testing.T) {
	codes, err := GenerateBackupCodes(50)
	if err != nil {
		t.Fatalf("GenerateBackupCodes() error: %v", err)
	}
	seen := make(map[string]bool, len(codes))
	for _, c := range codes {
		if seen[c] {
			t.Errorf("duplicate code: %s", c)
		}
		seen[c] = true
	}
}

// ── T13: Zero count returns empty slice ──────────────────────────────────────

func TestGenerateBackupCodes_ZeroCount(t *testing.T) {
	codes, err := GenerateBackupCodes(0)
	if err != nil {
		t.Fatalf("GenerateBackupCodes() error: %v", err)
	}
	if len(codes) != 0 {
		t.Errorf("expected empty slice, got %d codes", len(codes))
	}
}

// ── T14: HashBackupCode returns a hash ───────────────────────────────────────

func TestHashBackupCode_ReturnsHash(t *testing.T) {
	hash, err := HashBackupCode("ABCD-1234")
	if err != nil {
		t.Fatalf("HashBackupCode() error: %v", err)
	}
	if hash == "" {
		t.Error("HashBackupCode() returned empty string")
	}
}

// ── T15: VerifyBackupCode matches correct code ──────────────────────────────

func TestVerifyBackupCode_CorrectCode(t *testing.T) {
	code := "ABCD-1234"
	hash, err := HashBackupCode(code)
	if err != nil {
		t.Fatalf("HashBackupCode() error: %v", err)
	}
	if !VerifyBackupCode(code, hash) {
		t.Error("VerifyBackupCode() returned false for matching code")
	}
}

// ── T16: VerifyBackupCode rejects wrong code ─────────────────────────────────

func TestVerifyBackupCode_WrongCode(t *testing.T) {
	hash, err := HashBackupCode("ABCD-1234")
	if err != nil {
		t.Fatalf("HashBackupCode() error: %v", err)
	}
	if VerifyBackupCode("FFFF-9999", hash) {
		t.Error("VerifyBackupCode() returned true for wrong code")
	}
}

// ── T17: VerifyBackupCode is case-insensitive ────────────────────────────────

func TestVerifyBackupCode_DifferentCase(t *testing.T) {
	code := "ABCD-EF12"
	hash, err := HashBackupCode(code)
	if err != nil {
		t.Fatalf("HashBackupCode() error: %v", err)
	}
	if !VerifyBackupCode("abcd-ef12", hash) {
		t.Error("VerifyBackupCode() should be case-insensitive")
	}
}
