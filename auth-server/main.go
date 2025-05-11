// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	
	"auth-server/database"
	"auth-server/handlers"
	
	"github.com/gorilla/mux"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.CreateTables(); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/verify-otp", handlers.VerifyOTP).Methods("POST")
	r.HandleFunc("/getalluser", handlers.GetAllUsers).Methods("GET")
	r.HandleFunc("/createuser", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/getuser", handlers.GetUser).Methods("GET")
	r.HandleFunc("/updateuser", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/deleteuser", handlers.DeleteUser).Methods("DELETE")
	r.HandleFunc("/getauthpubaddr", handlers.GetAuthPubAddr).Methods("POST")
	r.HandleFunc("/getreqpubaddr", handlers.GetReqPubAddr).Methods("POST")
	
	// Add user CRUD routes here
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}