package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/techninja8/FileCheck/config"
)

var db, _ = config.InitDB()

// GetEmailFromToken extracts the email from a JWT token
func GetEmailFromToken(tokenString string) (string, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwtSecret")), nil
	})
	if err != nil {
		return "", fmt.Errorf("error parsing token: %v", err)
	}

	// Extract the email from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("email not found in token")
	}

	return email, nil
}

func GetTokenFromRequest(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}

func GetUsernameFromEmail(c *gin.Context, email string) (string, error) {

	// I have a variable (string known as email), how do i query my DB to get the associated username with that email, they are both in table users
	query := `SELECT username FROM users WHERE email = ?`
	var username string
	err := db.QueryRow(query, email).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user found with email: %v", email)
		}
		return "", fmt.Errorf("failed to query database, %v", err)
	}

	return username, nil
}
