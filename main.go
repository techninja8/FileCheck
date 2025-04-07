package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/techninja8/FileCheck/api/auth"
	"github.com/techninja8/FileCheck/api/handlers"
	"github.com/techninja8/FileCheck/api/middleware"
	"github.com/techninja8/FileCheck/config"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env vars")
	} else {
		log.Println(".env loaded")
	}

	// Initialize the database
	config.InitDB()

	// Initialize the router
	router := gin.Default()
	//router.Use(cors.Default()) // Enable CORS

	// Public routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	router.POST("/auth/login", auth.LoginHandler)
	router.POST("/auth/register", auth.RegisterHandler)

	// Protected routes
	authGroup := router.Group("/v1")
	authGroup.Use(middleware.JWTMiddleware())
	{
		authGroup.POST("/upload", handlers.UploadHandler)
		authGroup.GET("/download/:id", handlers.DownloadHandler)
		authGroup.GET("/check/:id", handlers.CheckIntegrityHandler)
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server is running on port", port)
	log.Fatal(router.Run(":" + port))
}
