package models

// User represents an application user.
type User struct {
	BaseModel
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Role     string `gorm:"size:50;default:user" json:"role"`
	Active   bool   `gorm:"default:true" json:"active"`
	Posts    []Post `gorm:"foreignKey:UserID" json:"posts,omitempty"`
}
