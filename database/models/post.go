package models

// Post represents a content entry authored by a user.
type Post struct {
	BaseModel
	Title  string `gorm:"size:255;not null" json:"title"`
	Slug   string `gorm:"size:255;uniqueIndex" json:"slug"`
	Body   string `gorm:"type:text" json:"body"`
	UserID uint   `gorm:"index;not null" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
