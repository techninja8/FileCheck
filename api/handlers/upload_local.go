package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/techninja8/FileCheck/api/middleware"
	"github.com/techninja8/FileCheck/config"
)

type UploadResponse struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
}

func UploadHandler(c *gin.Context) {
	// Get username from token
	// Append username to database
	db, err := config.InitDB()
	if err != nil {
		fmt.Printf("failed to initialize the database")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer file.Close()

	// Generate SHA-256 hash for the file
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate hash"})
		return
	}
	hashInBytes := hash.Sum(nil)[:32]
	fileHash := hex.EncodeToString(hashInBytes)

	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset file pointer"})
		return
	}

	// Save file locally
	storagePath := os.Getenv("LOCAL_STORAGE_PATH")
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		os.MkdirAll(storagePath, os.ModePerm)
	}

	filePath := filepath.Join(storagePath, header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	fileID := uuid.New().String()

	if db == nil {
		log.Println("Database is not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	token, err := middleware.GetTokenFromRequest(c)
	if err != nil {
		fmt.Printf("failed to get token from request, %v", err)
		return
	}

	useremail, err := middleware.GetEmailFromToken(token)
	if err != nil {
		fmt.Printf("failed to get email from token, %v", err)
		return
	}

	username, err := middleware.GetUsernameFromEmail(c, useremail)

	if err != nil {
		fmt.Printf("failed to get the username from token %v", err)
	}
	_, err = db.Exec(`INSERT INTO files (id, owner, filename, hash, uploaded_at, location) VALUES (?, ?, ?, ?, ?, ?)`,
		fileID, username, header.Filename, fileHash, time.Now(), filePath)

	if err != nil {
		log.Printf("Failed to insert file metadata: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store file metadata"})
		return
	}

	response := UploadResponse{
		ID:       fileID,
		Filename: header.Filename,
		Hash:     fileHash,
	}

	c.JSON(http.StatusOK, response)
}
