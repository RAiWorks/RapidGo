package models

import "time"

// BaseModel provides common fields for all models.
// Embed this in your model structs instead of gorm.Model
// to get ID, CreatedAt, and UpdatedAt without soft deletes.
type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
