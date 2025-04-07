package handlers

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/techninja8/FileCheck/api/middleware"
	"github.com/techninja8/FileCheck/config"
)

func DownloadHandler(c *gin.Context) {
	fileID := c.Param("id")

	db, err := config.InitDB()
	if err != nil {
		log.Printf("failed to initialize database")
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

	// Retrieve file metadata from SQLite3
	var filename, filePath string
	err = db.QueryRow("SELECT filename, location FROM files WHERE id = ?", fileID).Scan(&filename, &filePath)
	if err != nil {
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

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	io.Copy(c.Writer, file)
}
