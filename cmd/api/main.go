// Package main is the entry point for the User Management API server
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/frtasoniero/user-management-api/database"
	"github.com/frtasoniero/user-management-api/internal/repository"
	"github.com/frtasoniero/user-management-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// main initializes and starts the User Management API server
func main() {
	// Load environment variables from .env file (optional for development)
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading .env file")
	}

	// Initialize database connection to MongoDB
	database.ConnectToMongoDB()
	// Ensure database connection is closed when the application terminates
	defer database.DisconnectFromMongoDB()

	// Get database name from environment variable
	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		log.Fatal("MONGODB_DB_NAME environment variable is not set")
	}

	// Initialize repository layer with MongoDB database connection
	dbClient := database.MongoDBClient.Database(dbName)
	userRepo := repository.NewUserRepository(dbClient, "users")

	// Initialize Gin HTTP router with default middleware (logger and recovery)
	router := gin.Default()

	// Register all API routes and handlers
	routes.RegisterRoutes(router, userRepo)

	// Get server port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configure HTTP server with timeout and handler
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start HTTP server in a goroutine to allow for graceful shutdown
	go func() {
		log.Printf("Server running on port: %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Setup graceful shutdown - wait for interrupt signal (Ctrl+C, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until signal is received
	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown - finish existing requests within timeout
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server shutdown.")
}
