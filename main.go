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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize the database
	config.InitDB()

	// Initialize the router
	router := gin.Default()

	router.POST("/auth/login", auth.LoginHandler)
	router.POST("/auth/register", auth.RegisterHandler)

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
