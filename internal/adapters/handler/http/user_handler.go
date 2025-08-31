// Package http provides HTTP handlers for user management operations.
package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/frtasoniero/user-management-api/internal/core/domain"
	"github.com/frtasoniero/user-management-api/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUC ports.UserUseCase
}

type RegisterRequest struct {
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required,min=6"`
	Profile  domain.Profile `json:"profile" binding:"required"`
}

func NewUserHandler(userUC ports.UserUseCase) *UserHandler {
	return &UserHandler{
		userUC: userUC,
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

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	user, err := h.userUC.GetUserByID(c.Request.Context(), idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	// Parse pagination parameters from URL query
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	// Parse field selection from URL query
	var fields []string
	if fieldsParam := c.Query("fields"); fieldsParam != "" {
		fields = strings.Split(fieldsParam, ",")
		// Clean up field names (remove spaces)
		for i, field := range fields {
			fields[i] = strings.TrimSpace(field)
		}
	}

	// Parse search parameter
	search := strings.TrimSpace(c.Query("search"))

	// Parse sorting parameters
	sortBy := strings.TrimSpace(c.Query("sort"))
	if sortBy == "" {
		sortBy = "created_at" // Default sort field
	}

	// Validate sort field to prevent injection
	validSortFields := map[string]bool{
		"email":      true,
		"created_at": true,
		"updated_at": true,
		"first_name": true,
		"last_name":  true,
	}
	if !validSortFields[sortBy] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort field"})
		return
	}

	// Parse order parameter (asc or desc)
	order := strings.ToLower(strings.TrimSpace(c.Query("order")))
	if order != "asc" && order != "desc" {
		order = "asc" // Default order
	}

	// Build filter options
	filter := &ports.GetUsersOptions{
		Page:     page,
		PageSize: pageSize,
		Fields:   fields,
		Search:   search,
		SortBy:   sortBy,
		Order:    order,
	}

	users, err := h.userUC.GetUsers(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
