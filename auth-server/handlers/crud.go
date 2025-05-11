package handlers

import (
	"auth-server/database"
	"auth-server/logger"
	"auth-server/models"
	"auth-server/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var log logger.Logger


//init logger
func init() {
	config := logger.NewConfigFromEnv()
	var logerr error
	log, logerr = logger.NewLogger(config)
	if logerr != nil {
		fmt.Printf("Failed to initialize logger: %v\n", logerr)
		os.Exit(1)
	}
}


//create new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	var userDetails struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		ReqPubAddr  string `json:"reqpubaddr"`
		AuthPubAddr string `json:"authpubaddr"`
	}
	if err := json.NewDecoder(r.Body).Decode(&userDetails); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Error("Failed to decode request body", "error", err)
		return
	}

	// Check if the username already exists
	var existingUserID int
	err := database.DB.QueryRow(context.Background(),
		"SELECT id FROM users WHERE username = $1",
		userDetails.Username,
	).Scan(&existingUserID)
	if err == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		log.Error("Username already exists", "username", userDetails.Username)
		return
	}

	// Insert the new user
	var user models.User
	hashedPassword, err := utils.HashPassword(userDetails.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		log.Error("Failed to hash password", "error", err)
		return
	}
	userDetails.Password = hashedPassword
	err = database.DB.QueryRow(context.Background(),
		"INSERT INTO users (username, email, password, reqpubaddr, authpubaddr) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		userDetails.Username, userDetails.Email, userDetails.Password, userDetails.ReqPubAddr, userDetails.AuthPubAddr,
	).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Error("Failed to create user", "error", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}


//get user
func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	fmt.Println("User ID:", userID)
	var user models.User
	err := database.DB.QueryRow(context.Background(),
		"SELECT id, username, email FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		log.Error("User not found", "userID", userID, "error", err)
		return
	}
	json.NewEncoder(w).Encode(user)
}


//update user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var userDetails struct {
		ID         string `json:"id"`
		Username   string `json:"username"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		AuthPubAddr string `json:"authpubaddr"`
		ReqPubAddr  string `json:"reqpubaddr"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userDetails); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Error("Failed to decode request body", "error", err)
		return
	}

	// Check if the user ID exists
	var existingUserID string
	err := database.DB.QueryRow(context.Background(),
		"SELECT id FROM users WHERE id = $1",
		userDetails.ID,
	).Scan(&existingUserID)
	if err != nil {
		http.Error(w, "User ID does not exist", http.StatusNotFound)
		log.Error("User ID does not exist", "userID", userDetails.ID, "error", err)
		return
	}

	// Hash the password if it is provided
	var hashedPassword string
	if userDetails.Password != "" {
		hashedPassword, err = utils.HashPassword(userDetails.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			log.Error("Failed to hash password", "error", err)
			return
		}
	}

	// Update the user data
	_, err = database.DB.Exec(context.Background(),
		"UPDATE users SET username = $1, email = $2, password = COALESCE(NULLIF($3, ''), password), authpubaddr = $4, reqpubaddr = $5 WHERE id = $6",
		userDetails.Username, userDetails.Email, hashedPassword, userDetails.AuthPubAddr, userDetails.ReqPubAddr, userDetails.ID,
	)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		log.Error("Failed to update user", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "User updated successfully"}
	json.NewEncoder(w).Encode(response)
	log.Info("User updated successfully", "userID", userDetails.ID)
}


//delete user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Check if the user ID exists
	var existingUserID int
	err := database.DB.QueryRow(context.Background(),
		"SELECT id FROM users WHERE id = $1",
		userID,
	).Scan(&existingUserID)
	if err != nil {
		http.Error(w, "User ID does not exist", http.StatusNotFound)
		log.Error("User ID does not exist", "userID", userID, "error", err)
		return
	}

	_, err = database.DB.Exec(context.Background(),
		"DELETE FROM users WHERE id = $1",
		userID,
	)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		log.Error("Failed to delete user", "error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "User deleted successfully"}
	json.NewEncoder(w).Encode(response)
	log.Info("User deleted successfully", "userID", userID)
}

//get all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	rows, err := database.DB.Query(context.Background(),
		"SELECT id, username, email FROM users",
	)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		log.Error("Failed to fetch users", "error", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			log.Error("Failed to scan user", "error", err)
			return
		}
		users = append(users, user)
	}

	log.Info("Fetched all users successfully", "count", len(users))
	json.NewEncoder(w).Encode(users)
}