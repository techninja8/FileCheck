package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techninja8/FileCheck/api/middleware"
	"github.com/techninja8/FileCheck/config"
	"github.com/techninja8/FileCheck/db/models"
)

var db, _ = config.InitDB()

// Auth Handlers
func LoginHandler(c *gin.Context) {
	var loginRequest struct {
		//Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		log.Printf("Error binding login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := models.GetUserByEmail(db, loginRequest.Email)
	if err != nil {
		log.Printf("User not found: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := user.CheckPassword(loginRequest.Password); err != nil {
		//log.Printf("Invalid password for user: %s", loginRequest.Username)
		log.Printf("Invalid password for user: %s", loginRequest.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.Username, user.Email, "user")
	if err != nil {
		log.Printf("Error generating token for user: %s, error: %v", loginRequest.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("token", token, 24*3600, "/", "", false, true) // Cookie set to exists for 24 hours
	fmt.Println("token saved as cookie")

	c.JSON(http.StatusOK, gin.H{"token": token})
}
