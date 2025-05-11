// internal/utils/crypto.go
package utils

import (
	"auth-server/database"
	"auth-server/logger"
	"context"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
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


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	fmt.Println("Checking password hash....")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


func GetReqPubAddrFromUsername(username string) (string, error) {
	var reqPubAddr string
	err := database.DB.QueryRow(context.Background(),
		"SELECT reqpubaddr FROM users WHERE username = $1",
		username,
	).Scan(&reqPubAddr)
	if err != nil {
		log.Error("Failed to fetch reqpubaddr", "username", username, "error", err)
		return "", fmt.Errorf("failed to fetch reqpubaddr for username %s: %v", username, err)
	}
	return reqPubAddr, nil
}

func GetAuthPubAddrFromUsername(username string) (string, error) {
	var authPubAddr string
	err := database.DB.QueryRow(context.Background(),
		"SELECT authpubaddr FROM users WHERE username = $1",
		username,
	).Scan(&authPubAddr)
	if err != nil {
		log.Error("Failed to fetch authpubaddr", "username", username, "error", err)
		return "", fmt.Errorf("failed to fetch authpubaddr for username %s: %v", username, err)
	}
	return authPubAddr, nil
}
