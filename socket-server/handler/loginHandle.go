package handler

// import (
// 	"net/http"
// 	"encoding/json"
// 	"github.com/gin-gonic/gin"
// )

// type LoginRequest struct {
// 	Username string `json:"username" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// }


// //login handler

// func (h *Handler) LoginHandle(c *gin.Context) {
// 	// Get the request body
// 	var loginRequest LoginRequest
// 	if err := c.ShouldBindJSON(&loginRequest); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	// Validate the request
// 	if loginRequest.Username == "" || loginRequest.Password == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
// 		return
// 	}

// 	// Check if the user exists in the database
// 	user, err := h.db.GetUserByUsername(loginRequest.Username)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		return
// 	}
// 	if user == nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
// 		return
// 	}

// 	// Check if the password is correct
// 	if user.Password != loginRequest.Password {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
// 		return
// 	}

// 	// Generate a JWT token for the user
// 	token, err := h.auth.GenerateToken(user.ID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": token})
// }