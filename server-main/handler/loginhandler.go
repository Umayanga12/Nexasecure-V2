package handler

import (
	"encoding/json"
	"fmt"

	"net/http"
	"nexasecure/logger"
	"os"
	"nexasecure/util"

	"github.com/gin-gonic/gin"
)


func LoginHandler(c *gin.Context) {
	config := logger.NewConfigFromEnv()

	// Initialize logger
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Get credentials from request
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		log.Error("Invalid JSON payload: %v", err)
		return
	}

	// Validate input
	if credentials.Username == "" || credentials.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		log.Error("Username or password is empty")
		return
	}

	// Get external endpoint
	externalEndpoint := "http://127.0.0.1:18909/login" // or os.Getenv("LOGIN_ENDPOINT")
	if externalEndpoint == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		log.Error("LOGIN_ENDPOINT not set in environment variables")
		return
	}

	// Prepare payload
	payload := map[string]string{
		"username": credentials.Username,
		"password": credentials.Password,
	}
	jsonPayload, _ := json.Marshal(payload)

	// Send request using MakeAPICall
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	statusCode, body, err := util.MakeAPICall("POST", externalEndpoint, headers, jsonPayload)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Error connecting to external service"})
		log.Error("Error connecting to external service: %v", err)
		return
	}

	if statusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed"})
		log.Error("Login failed with status code: %d", statusCode)
		return
	}

	// Parse response for "message" field
	var response struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing response"})
		log.Error("Error parsing response: %v", err)
		return
	}

	// Log and return based on response message
	if response.Message == "User authenticated" {
		if util.ManageNFTforLogin(credentials.Username) {
			c.JSON(http.StatusOK, gin.H{"message": "User Authenticated Successfully"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to manage NFT for login"})
			log.Error("Failed to manage NFT for user: %s", credentials.Username)
			return
		}

		// c.JSON(http.StatusOK, gin.H{"message": "User Authenticated Successfully"})
		// log.Info("Initial Login successful for user: %s.", credentials.Username, response.Message)
		
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed"})
		log.Error("Login failed for user: %s", credentials.Username)
	}
}

func MFAHandler(c *gin.Context) {
	config := logger.NewConfigFromEnv()

	// Initialize logger
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	var MFA struct {
		OTP      string `json:"otp"`
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&MFA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		log.Error("Invalid JSON payload: %v", err)
		return
	}
	if MFA.OTP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP is required"})
		log.Error("OTP is empty")
		return
	}

	// Get external endpoint
	externalEndpoint := "http://localhost:8000/verify-otp" // or os.Getenv("MFA_ENDPOINT")
	if externalEndpoint == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		log.Error("MFA_ENDPOINT not set in environment variables")
		return
	}

	payload := map[string]string{
		"otp":      MFA.OTP,
		"username": MFA.Username,
	}

	jsonPayload, _ := json.Marshal(payload)
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	statusCode, body, err := util.MakeAPICall("POST", externalEndpoint, headers, jsonPayload)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Error connecting to external service"})
		log.Error("Error connecting to external service: %v", err)
		return
	}
	if statusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "MFA failed"})
		log.Error("MFA failed with status code: %d", statusCode)
		return
	}

	// Parse response for "verified" field
	var response struct {
		Verified bool `json:"verified"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing response"})
		log.Error("Error parsing response: %v", err)
		return
	}

	if response.Verified {
		
		log.Info("MFA successful for user: %s", MFA.Username)
		//Get Auth NFT 
		go util.ManageNFTforLogin(MFA.Username)
		c.JSON(http.StatusOK, gin.H{"message": "MFA successful"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "MFA failed"})
		log.Error("MFA failed for user: %s", MFA.Username)
	}
}



