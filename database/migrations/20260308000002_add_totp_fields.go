package migrations

import "gorm.io/gorm"

func init() {
	Register(Migration{
		Version: "20260308000002_add_totp_fields",
		Up: func(db *gorm.DB) error {
			type User struct {
				TOTPEnabled     bool   `gorm:"default:false"`
				TOTPSecret      string `gorm:"size:512"`
				TOTPVerifiedAt  *string
				BackupCodesHash string `gorm:"type:text"`
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
				TOTPVerifiedAt  *string
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
