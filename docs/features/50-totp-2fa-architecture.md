# 📐 Architecture: Two-Factor Authentication (TOTP)

> **Feature**: `50` — Two-Factor Authentication (TOTP)
> **Discussion**: [`50-totp-2fa-discussion.md`](50-totp-2fa-discussion.md)
> **Status**: 🟡 IN PROGRESS
> **Date**: 2026-03-07

---

## Overview

Feature #50 adds a `core/totp` package that wraps `pquerna/otp` to provide TOTP secret generation, code validation, provisioning URI creation, and backup code management. The User model gains TOTP fields, and a migration adds the columns. The framework provides building blocks — application developers compose auth flows on top.

---

## File Structure

```
core/
  totp/
    totp.go       ← NEW — TOTP operations (generate, validate, backup codes)
    totp_test.go  ← NEW — unit tests
database/
  models/
    user.go       ← MODIFIED — add TOTP fields
  migrations/
    20260308000002_add_totp_fields.go  ← NEW — migration for users table
```

---

## Component Design

### 1. TOTP Package

**Package**: `core/totp`
**File**: `totp.go`

```go
package totp

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// Key wraps an OTP key with its secret and provisioning URI.
type Key struct {
	Secret string // Base32-encoded secret
	URL    string // otpauth:// provisioning URI for QR codes
}

// GenerateKey creates a new TOTP secret for the given issuer and account.
// The issuer appears in the authenticator app (e.g. "RapidGo").
// The account identifies the user (typically their email).
func GenerateKey(issuer, account string) (*Key, error)

// ValidateCode checks whether a TOTP code is valid for the given secret.
// Uses a time window of ±1 period (30 seconds) for clock drift tolerance.
func ValidateCode(secret, code string) bool

// GenerateBackupCodes creates n random backup codes in "XXXX-XXXX" format.
func GenerateBackupCodes(n int) ([]string, error)

// HashBackupCode returns a bcrypt hash of a backup code.
func HashBackupCode(code string) (string, error)

// VerifyBackupCode checks a plaintext code against a bcrypt hash.
func VerifyBackupCode(code, hash string) bool
```

**Implementation details**:

```go
func GenerateKey(issuer, account string) (*Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: account,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, err
	}
	return &Key{
		Secret: key.Secret(),
		URL:    key.URL(),
	}, nil
}

func ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}
```

**Backup code format**: `XXXX-XXXX` (8 hex chars with dash, uppercase) — human-readable, copy-paste friendly. Codes are normalized to uppercase before hashing/verification for case-insensitive matching.

**Backup code hashing**: Uses `bcrypt` (via `golang.org/x/crypto`) — same pattern as `helpers.HashPassword()`.

---

### 2. User Model Updates

**Package**: `database/models`
**File**: `user.go`

```go
type User struct {
	BaseModel
	Name            string     `gorm:"size:100;not null" json:"name"`
	Email           string     `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password        string     `gorm:"size:255;not null" json:"-"`
	Role            string     `gorm:"size:50;default:user" json:"role"`
	Active          bool       `gorm:"default:true" json:"active"`
	TOTPEnabled     bool       `gorm:"default:false" json:"totp_enabled"`
	TOTPSecret      string     `gorm:"size:512" json:"-"`
	TOTPVerifiedAt  *time.Time `json:"totp_verified_at,omitempty"`
	BackupCodesHash string     `gorm:"type:text" json:"-"`
	Posts           []Post     `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}
```

**Field details**:

| Field | Type | Tags | Purpose |
|---|---|---|---|
| `TOTPEnabled` | `bool` | `gorm:"default:false"` | Whether 2FA is active for this user |
| `TOTPSecret` | `string` | `gorm:"size:512" json:"-"` | AES-256-GCM encrypted base32 secret; excluded from JSON |
| `TOTPVerifiedAt` | `*time.Time` | `json:"totp_verified_at,omitempty"` | Timestamp of initial TOTP enrollment verification |
| `BackupCodesHash` | `string` | `gorm:"type:text" json:"-"` | JSON array of bcrypt hashes; excluded from JSON |

**Security notes**:
- `TOTPSecret` and `BackupCodesHash` use `json:"-"` — never exposed in API responses
- `TOTPSecret` is encrypted with `core/crypto.Encrypt()` before storage
- Backup codes are bcrypt-hashed — irreversible, constant-time comparison

---

### 3. Migration

**Package**: `database/migrations`
**File**: `20260308000002_add_totp_fields.go`

```go
func init() {
	Register(Migration{
		Version: "20260308000002_add_totp_fields",
		Up: func(db *gorm.DB) error {
			type User struct {
				TOTPEnabled     bool       `gorm:"default:false"`
				TOTPSecret      string     `gorm:"size:512"`
				TOTPVerifiedAt  *time.Time
				BackupCodesHash string     `gorm:"type:text"`
			}
			if err := db.Migrator().AddColumn(&User{}, "TOTPEnabled"); err != nil {
				return err
			}
			if err := db.Migrator().AddColumn(&User{}, "TOTPSecret"); err != nil {
				return err
			}
			if err := db.Migrator().AddColumn(&User{}, "TOTPVerifiedAt"); err != nil {
				return err
			}
			return db.Migrator().AddColumn(&User{}, "BackupCodesHash")
		},
		Down: func(db *gorm.DB) error {
			type User struct {
				TOTPEnabled     bool
				TOTPSecret      string
				TOTPVerifiedAt  *time.Time
				BackupCodesHash string
			}
			if err := db.Migrator().DropColumn(&User{}, "BackupCodesHash"); err != nil {
				return err
			}
			if err := db.Migrator().DropColumn(&User{}, "TOTPVerifiedAt"); err != nil {
				return err
			}
			if err := db.Migrator().DropColumn(&User{}, "TOTPSecret"); err != nil {
				return err
			}
			return db.Migrator().DropColumn(&User{}, "TOTPEnabled")
		},
	})
}
```

---

## Usage Example (Application-Level)

The framework provides building blocks. An application would compose a 2FA enrollment flow like:

```go
// In an auth controller — enable TOTP for a user
func EnableTOTP(c *gin.Context) {
	userID := c.GetUint("user_id")
	issuer := os.Getenv("TOTP_ISSUER")

	// 1. Generate TOTP key
	key, _ := totp.GenerateKey(issuer, user.Email)

	// 2. Encrypt secret for storage
	encKey, _ := hex.DecodeString(os.Getenv("TOTP_ENCRYPTION_KEY"))
	encrypted, _ := crypto.Encrypt(key.Secret, encKey)

	// 3. Save to user (2FA not enabled until verified)
	db.Model(&user).Update("totp_secret", encrypted)

	// 4. Return provisioning URI for QR code
	c.JSON(200, gin.H{"url": key.URL})
}

// Verify TOTP during login
func VerifyTOTP(c *gin.Context) {
	code := c.PostForm("code")

	// 1. Decrypt stored secret
	encKey, _ := hex.DecodeString(os.Getenv("TOTP_ENCRYPTION_KEY"))
	secret, _ := crypto.Decrypt(user.TOTPSecret, encKey)

	// 2. Validate code
	if totp.ValidateCode(secret, code) {
		// Issue full auth token
	}
}
```

---

## Dependencies

- **NEW**: `github.com/pquerna/otp` — TOTP/HOTP library (Apache 2.0)
- **Existing**: `golang.org/x/crypto` — bcrypt for backup code hashing (already in go.mod)
- **Existing**: `core/crypto` — AES-256-GCM encryption for TOTP secret storage
- No other new dependencies

---

## Environment Variables

| Variable | Required | Default | Purpose |
|---|---|---|---|
| `TOTP_ISSUER` | No | `"RapidGo"` | App name shown in authenticator app |
| `TOTP_ENCRYPTION_KEY` | Yes (if TOTP used) | — | 64-char hex string (32 bytes decoded, e.g. `openssl rand -hex 32`) for AES-256-GCM encryption of TOTP secrets |

---

## Next

Tasks doc → `50-totp-2fa-tasks.md`
