package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/phi-lani/blockchainApp/controllers"
	"github.com/phi-lani/blockchainApp/database"
	"github.com/phi-lani/blockchainApp/middleware"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the database connection
	database.InitializeDB()

	// Create a new router
	router := mux.NewRouter()

	// Middleware for logging and authentication
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.AuthenticationMiddleware)

	// Define routes
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/profile/{id}", controllers.UpdateProfile).Methods("PUT")

	// Serve the application
	log.Printf("Server starting on port %s ...", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}
