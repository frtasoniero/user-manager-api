// Package http provides HTTP handlers for user management operations.
package http

import (
	"net/http"

	"github.com/frtasoniero/user-management-api/internal/core/domain"
	"github.com/frtasoniero/user-management-api/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUC   ports.UserUseCase
	userRepo ports.UserRepository
}

type RegisterRequest struct {
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required,min=6"`
	Profile  domain.Profile `json:"profile" binding:"required"`
}

func NewUserHandler(userUC ports.UserUseCase, userRepo ports.UserRepository) *UserHandler {
	return &UserHandler{
		userUC:   userUC,
		userRepo: userRepo,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userUC.Register(c.Request.Context(), req.Email, req.Password, req.Profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
