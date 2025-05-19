// internal/handlers/auth.go
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	//"time"

	"auth-server/database"
	"auth-server/logger"

	//"auth-server/logger"
	"auth-server/models"
	"auth-server/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	config := logger.NewConfigFromEnv()

	// Initialize logger
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	//fmt.Println(creds) // Corrected to use fmt.Println instead of fmt.println
	// No need to check if creds is nil as structs cannot be nil as structs cannot be nil
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		fmt.Println("Invalid request payload")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received credentials: %+v\n", creds)

	var user models.User
	err = database.DB.QueryRow(context.Background(),
		"SELECT id, username, email, password FROM users WHERE username = $1",
		creds.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	//fmt.Println("User details fetched from database:", user)
	if err != nil {
		//log.Warn("User not found or error occurred for username: %s, error: %v", creds.Username, err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// fmt.Printf("Fetched user details from database: %+v\n", user)

	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		//log.Warn("Invalid password for user: %s", creds.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	fmt.Printf("User authenticated: %s\n", creds.Username)

	// Generate secure OTP
	// otp, err := generateSecureOTP()
	// if err != nil {
	// 	// log.Error("Failed to generate OTP: %v", err)
	// 	http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
	// 	return
	// }
	// expiresAt := time.Now().Add(5 * time.Minute)

	//fmt.Println("Generated OTP: %s, Expires At: %v", otp, expiresAt)

	// Store OTP in the database
	// _, err = database.DB.Exec(context.Background(),
	// 	"INSERT INTO otps (user_id, code, expires_at) VALUES ($1, $2, $3)",
	// 	user.ID, otp, expiresAt,
	// )
	// if err != nil {
	// 	//log.Error("Failed to store OTP for user: %s, error: %v", creds.Username, err) // Re-enabled logging for error
	// 	http.Error(w, "Failed to generate OTP", http.StatusInternalServerError)
	// 	return
	// }

	// // //log.Info("OTP generated and stored for user: %s", creds.Username)

	// if err := utils.SendOTPEmail(user.Email, otp); err != nil {
	// 	//log.Error("Failed to send OTP email to user: %s, error: %v", creds.Username, err)
	// 	http.Error(w, "Failed to send OTP", http.StatusInternalServerError)
	// 	fmt.Print(err)
	// 	return
	// }

	// fmt.Printf("OTP sent to email for user: %s\n", creds.Username)

	// Respond indicating OTP was sent (adjusted for clarity)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User authenticated"})
	//log.Info("Login process completed for user: %s", creds.Username)
}

// generateSecureOTP generates a 6-digit OTP using crypto/rand
func generateSecureOTP() (string, error) {
	b := make([]byte, 4) // 4 bytes for 32 bits, sufficient for 6 digits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Convert bytes to a 32-bit unsigned integer
	num := uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
	otpNum := num % 1000000
	return fmt.Sprintf("%06d", otpNum), nil
}
func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		OTP      string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var otp models.OTP
	err := database.DB.QueryRow(context.Background(),
		`SELECT user_id, code, expires_at 
		FROM otps 
		WHERE user_id = (SELECT id FROM users WHERE username = $1) 
		AND code = $2 
		AND expires_at > NOW()`,
		req.Username, req.OTP,
	).Scan(&otp.UserID, &otp.Code, &otp.ExpiresAt)

	if err != nil {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	// Cleanup OTP
	database.DB.Exec(context.Background(),
		"DELETE FROM otps WHERE user_id = $1",
		otp.UserID,
	)

	json.NewEncoder(w).Encode(map[string]bool{"verified": true})
}
