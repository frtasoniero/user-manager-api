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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading .env file")
	}

	database.ConnectToMongoDB()
	defer database.DisconnectFromMongoDB()

	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		log.Fatal("MONGODB_DB_NAME environment variable is not set")
	}

	userRepo := repository.NewUserRepository(database.MongoDBClient.Database(dbName), "users")

	router := gin.Default()

	routes.RegisterRoutes(router, userRepo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Servidor rodando na porta: %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Falha ao iniciar o servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Servidor forçado a desligar: ", err)
	}

	log.Println("Servidor desligado.")
}
