// models/user.go
package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User - Represents a user in the system
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserCredentials - For login
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserProfileUpdate - Represents profile update information
type UserProfileUpdate struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

// CreateUser - Saves a new user to the database
func CreateUser(db *gorm.DB, user *User) error {
	result := db.Create(user)
	return result.Error
}

// AuthenticateUser - Validates user credentials and returns the user if successful
func AuthenticateUser(db *gorm.DB, credentials UserCredentials) (User, error) {
	var user User
	err := db.Where("email = ?", credentials.Email).First(&user).Error
	if err != nil {
		return user, errors.New("user not found")
	}

	// Verify password
	if !CheckPasswordHash(credentials.Password, user.Password) {
		return user, errors.New("incorrect password")
	}

	return user, nil
}

// UpdateUserProfile - Updates a user's profile in the database
func UpdateUserProfile(db *gorm.DB, userID string, updates *UserProfileUpdate) error {
	return db.Model(&User{}).Where("id = ?", userID).Updates(updates).Error
}

// HashPassword - Hashes the given password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash - Compares the given hash with the password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
