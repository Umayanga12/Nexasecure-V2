package main


import (
	"os"
	"fmt"
	"github.com/gin-gonic/gin"
	"nexasecure/logger"
	"nexasecure/handler"
)

func main() {
	config := logger.NewConfigFromEnv()

	// Initialize logger
	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	router := gin.Default()

	router.POST("/login", handler.LoginHandler)
	//router.POST("/verifyotp", handler.MFAHandler)
	router.POST("/logout", handler.LogoutHandler)
	
	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	router.Run(":" + port)
	log.Info("Server started on port " + port)
}
