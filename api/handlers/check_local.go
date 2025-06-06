package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/techninja8/FileCheck/api/middleware"
	"github.com/techninja8/FileCheck/config"
)

type CheckResponse struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	Valid    bool   `json:"valid"`
}

func CheckIntegrityHandler(c *gin.Context) {
	fileID := c.Param("id")

	db, err := config.InitDB()
	if err != nil {
		log.Printf("failed to initialize database: %v", err)
		return
	}

	token, err := middleware.GetTokenFromRequest(c)
	if err != nil {
		log.Printf("failed to get token from request, %v", err)
		return
	}

	validAccess, err := middleware.VerifyFileOwnership(c, db, fileID, token)
	if !validAccess {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Access to file denied, pls login to get proper token for user"})
		return
	}
	if err != nil {
		log.Printf("failed to check access, %v", err)
		return
	}

	var filename, storedHash, filePath string

	log.Printf("Retrieving file metadata for ID: %s", fileID)

	err = db.QueryRow("SELECT filename, hash, location FROM files WHERE id = ?", fileID).Scan(&filename, &storedHash, &filePath)
	if err != nil {
		log.Printf("Error retrieving file metadata: %v", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file metadata"})
		}
		return
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
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

	// Compare the generated hash with the stored hash
	valid := fileHash == storedHash

	response := CheckResponse{
		ID:       fileID,
		Filename: filename,
		Hash:     fileHash,
		Valid:    valid,
	}

	c.JSON(http.StatusOK, response)
}
