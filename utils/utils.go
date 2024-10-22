package utils

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/phi-lani/blockchainApp/models"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword - Hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash - Compares the password with the hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT - Generates a JWT token for the user
func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get JWT secret from environment variables
	secret := os.Getenv("JWT_SECRET")

	// Sign the token with the secret
	return token.SignedString([]byte(secret))
}

// ValidateJWT - Validates a JWT token and returns the token if it's valid
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// Parse the token
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
