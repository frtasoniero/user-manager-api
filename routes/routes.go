// Package routes provides HTTP route definitions for the user management API.
package routes

import (
	"net/http"

	handler "github.com/frtasoniero/user-management-api/internal/adapters/handler/http"
	"github.com/frtasoniero/user-management-api/internal/core/usecase"
	"github.com/frtasoniero/user-management-api/internal/repository"
	"github.com/gin-gonic/gin"

	// Swagger imports
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router *gin.Engine, userRepo *repository.UserRepository) {
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	// Swagger documentation endpoint
	// Access at: http://localhost:8080/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiGroup := router.Group("/api/v1")
	{
		apiGroup.GET("/health", healthCheck)

		// User routes
		apiGroup.GET("/users", userHandler.GetUsers)
		apiGroup.GET("/users/:id", userHandler.GetUserByID)
		apiGroup.POST("/users/register", userHandler.Register)
	}
}

// healthCheck godoc
// @Summary Health check endpoint
// @Description Check if the API server is running and healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "API is healthy"
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
