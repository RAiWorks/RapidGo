package session

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// DBStore stores sessions in a database table via GORM.
type DBStore struct {
	DB *gorm.DB
}

// SessionRecord represents a row in the sessions table.
type SessionRecord struct {
	ID         string    `gorm:"primaryKey;size:255"`
	Data       string    `gorm:"type:text;not null"`
	UserID     *uint     `gorm:"index"`
	IPAddress  *string   `gorm:"size:45"`
	UserAgent  *string   `gorm:"type:text"`
	LastActive time.Time `gorm:"autoUpdateTime"`
	CreatedAt  time.Time
}

// TableName returns the database table name.
func (SessionRecord) TableName() string { return "sessions" }

func (s *DBStore) Read(id string) (map[string]interface{}, error) {
	var rec SessionRecord
	if err := s.DB.Where("id = ?", id).First(&rec).Error; err != nil {
		return make(map[string]interface{}), nil
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(rec.Data), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *DBStore) Write(id string, data map[string]interface{}, lifetime time.Duration) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	rec := SessionRecord{
		ID:         id,
		Data:       string(raw),
		LastActive: time.Now(),
	}
	return s.DB.Save(&rec).Error
}

func (s *DBStore) Destroy(id string) error {
	return s.DB.Where("id = ?", id).Delete(&SessionRecord{}).Error
}

func (s *DBStore) GC(maxLifetime time.Duration) error {
	cutoff := time.Now().Add(-maxLifetime)
	return s.DB.Where("last_active < ?", cutoff).Delete(&SessionRecord{}).Error
}
