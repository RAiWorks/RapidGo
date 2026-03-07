package services

import (
	"errors"

	"github.com/RAiWorks/RapidGo/database/models"

	"gorm.io/gorm"
)

// UserService handles business logic for user operations.
type UserService struct {
	DB *gorm.DB
}

// NewUserService creates a new UserService with the given database connection.
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// GetByID retrieves a user by their ID.
func (s *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user after checking for duplicate email.
func (s *UserService) Create(name, email, password string) (*models.User, error) {
	existing := &models.User{}
	if err := s.DB.Where("email = ?", email).First(existing).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: password, // hash before saving in real code
	}
	if err := s.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Update updates a user's fields by ID.
func (s *UserService) Update(id uint, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	if err := s.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Delete deletes a user by ID.
func (s *UserService) Delete(id uint) error {
	return s.DB.Delete(&models.User{}, id).Error
}
