package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techninja8/FileCheck/db/models"
)

// RegisterHandler handles user registration
func RegisterHandler(c *gin.Context) {
	var registerRequest struct {
		Username             string `json:"username" binding:"required"`
		Email                string `json:"email" binding:"required"`
		Password             string `json:"password" binding:"required"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Binding failed"})
		return
	}

	if registerRequest.Password != registerRequest.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Password does not match confirmation"})
		return
	}

	//db := c.MustGet("db").(*sql.DB)
	user := &models.User{
		Username: registerRequest.Username,
		Email:    registerRequest.Email,
	}

	if err := user.HashPassword(registerRequest.Password); err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user, password"})
		return
	}

	if err := models.CreateUser(db, user); err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user, can't create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
