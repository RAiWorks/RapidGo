package models

import (
	"strings"
	"time"

	"github.com/RAiWorks/RapidGo/app/helpers"
	"gorm.io/gorm"
)

// User represents an application user.
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

// BeforeCreate hashes the password before inserting into the database.
// Skips if the password is already bcrypt-hashed (starts with "$2a$" or "$2b$").
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" && !strings.HasPrefix(u.Password, "$2a$") && !strings.HasPrefix(u.Password, "$2b$") {
		hashed, err := helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashed
	}
	return nil
}
