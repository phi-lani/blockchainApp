package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/phi-lani/blockchainApp/database"
	"github.com/phi-lani/blockchainApp/models"
	"github.com/phi-lani/blockchainApp/utils"
)

// RegisterUser - Handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	fmt.Println("Received user input for registration:", user)

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	if err := models.CreateUser(database.DB, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// LoginUser - Handles user login and sends MFA token
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials models.UserCredentials
	_ = json.NewDecoder(r.Body).Decode(&credentials)

	fmt.Println("Login attempt with credentials:", credentials)

	user, err := models.AuthenticateUser(database.DB, credentials)
	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	mfaToken := utils.GenerateMFAToken()

	err = models.StoreMFAToken(database.DB, user.ID, mfaToken)
	if err != nil {
		http.Error(w, "Failed to store MFA token", http.StatusInternalServerError)
		return
	}

	err = utils.SendMFAToken(user.Email, mfaToken)
	if err != nil {
		http.Error(w, "Failed to send MFA token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "MFA token has been sent. Please verify to complete login."})
}

// ValidateMFAToken - Handles MFA token validation
func ValidateMFAToken(w http.ResponseWriter, r *http.Request) {
	var mfaData struct {
		UserID   uint   `json:"user_id"`
		MFAToken string `json:"mfa_token"`
	}

	_ = json.NewDecoder(r.Body).Decode(&mfaData)

	var storedMFAToken models.MFAToken
	err := database.DB.Where("user_id = ? AND token = ?", mfaData.UserID, mfaData.MFAToken).First(&storedMFAToken).Error
	if err != nil || storedMFAToken.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Invalid or expired MFA token", http.StatusUnauthorized)
		return
	}

	user := models.User{}
	database.DB.First(&user, mfaData.UserID)

	jwtToken, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"jwt_token": jwtToken})
}

// UpdateProfile - Allows users to update their profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var updates models.UserProfileUpdate
	params := mux.Vars(r)
	userID := params["id"]

	_ = json.NewDecoder(r.Body).Decode(&updates)

	if err := models.UpdateUserProfile(database.DB, userID, &updates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Profile updated successfully"})
}
