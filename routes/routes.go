// Package routes provides HTTP route definitions for the user management API.
package routes

import (
	"net/http"

	handler "github.com/frtasoniero/user-management-api/internal/adapters/handler/http"
	"github.com/frtasoniero/user-management-api/internal/core/usecase"
	"github.com/frtasoniero/user-management-api/internal/repository"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, userRepo *repository.UserRepository) {
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	apiGroup := router.Group("/api/v1")
	{
		apiGroup.GET("/health", healthCheck)
		apiGroup.GET("/users", userHandler.GetUsers)
		apiGroup.GET("/users/:id", userHandler.GetUserByID)
		apiGroup.POST("/users/register", userHandler.Register)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
