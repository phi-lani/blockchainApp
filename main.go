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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	database.InitializeDB()

	router := mux.NewRouter()

	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	router.HandleFunc("/validate-mfa", controllers.ValidateMFAToken).Methods("POST")

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthenticationMiddleware)
	protected.HandleFunc("/profile/{id}", controllers.UpdateProfile).Methods("PUT")

	log.Printf("Server starting on port %s ...", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}
