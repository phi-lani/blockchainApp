package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/phi-lani/blockchainApp/database"
	"github.com/phi-lani/blockchainApp/models"
	"github.com/phi-lani/blockchainApp/utils"

	"github.com/gorilla/mux"
)

// RegisterUser - Handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Unable to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// Save user in database
	if err := models.CreateUser(database.DB, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// LoginUser - Handles user login and returns JWT
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials models.UserCredentials
	_ = json.NewDecoder(r.Body).Decode(&credentials)

	// Authenticate user
	user, err := models.AuthenticateUser(database.DB, credentials)
	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return token as JSON response
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// UpdateProfile - Allows users to update their profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	var updates models.UserProfileUpdate
	_ = json.NewDecoder(r.Body).Decode(&updates)

	// Update the user's profile in the database
	if err := models.UpdateUserProfile(database.DB, userID, &updates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Profile updated successfully"})
}
