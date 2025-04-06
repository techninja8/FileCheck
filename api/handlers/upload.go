package handlers

/*import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
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

	// Upload file to S3
	s3Client := s3.NewFromConfig(aws.Config{
		Region: os.Getenv("AWS_REGION"),
		Credentials: aws.Credentials{
			AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	})

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(header.Filename),
		Body:   file,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to S3"})
		return
	}

	// Store file metadata in PostgreSQL
	var fileID int
	err = config.DB.QueryRow("INSERT INTO files (filename, hash, uploaded_at) VALUES ($1, $2, $3) RETURNING id",
		header.Filename, fileHash, time.Now()).Scan(&fileID)
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
*/
