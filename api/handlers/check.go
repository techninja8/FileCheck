package handlers

/* import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/techninja8/FileCheck/config"
)

type CheckResponse struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	Valid    bool   `json:"valid"`
}

func CheckIntegrityHandler(c *gin.Context) {
	fileID := c.Param("id")

	// Retrieve file metadata from PostgreSQL
	var filename, storedHash string
	err := config.DB.QueryRow("SELECT filename, hash FROM files WHERE id = $1", fileID).Scan(&filename, &storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file metadata"})
		}
		return
	}

	// Download file from S3
	s3Client := s3.NewFromConfig(aws.Config{
		Region: os.Getenv("AWS_REGION"),
		Credentials: aws.Credentials{
			AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	})

	output, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(filename),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file from S3"})
		return
	}
	defer output.Body.Close()

	// Generate SHA-256 hash for the downloaded file
	hash := sha256.New()
	if _, err := io.Copy(hash, output.Body); err != nil {
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
*/
