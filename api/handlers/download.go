package handlers

/* import (
	"database/sql"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/techninja8/FileCheck/config"
)

func DownloadHandler(c *gin.Context) {
	fileID := c.Param("id")

	// Retrieve file metadata from PostgreSQL
	var filename string
	err := config.DB.QueryRow("SELECT filename FROM files WHERE id = $1", fileID).Scan(&filename)
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

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	io.Copy(c.Writer, output.Body)
}
*/
