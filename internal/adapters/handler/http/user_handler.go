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

func NewUserHandler(userUC ports.UserUseCase, userRepo ports.UserRepository) *UserHandler {
	return &UserHandler{
		userUC:   userUC,
		userRepo: userRepo,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userUC.Register(c.Request.Context(), user.Email, user.PasswordHash, user.Profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
