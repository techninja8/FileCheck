package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/techninja8/FileCheck/config"
)

type UploadResponse struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
}

func UploadHandler(c *gin.Context) {
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

	// Store file metadata in SQLite3
	var fileID int
	err = config.DB.QueryRow("INSERT INTO files (filename, hash, uploaded_at, location) VALUES (?, ?, ?, ?) RETURNING id",
		header.Filename, fileHash, time.Now(), filePath).Scan(&fileID)
	if err != nil {
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
