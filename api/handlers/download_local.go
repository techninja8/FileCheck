package handlers

import (
	"database/sql"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/techninja8/FileCheck/config"
)

func DownloadHandler(c *gin.Context) {
	fileID := c.Param("id")

	// Retrieve file metadata from SQLite3
	var filename, filePath string
	err := config.DB.QueryRow("SELECT filename, location FROM files WHERE id = ?", fileID).Scan(&filename, &filePath)
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
