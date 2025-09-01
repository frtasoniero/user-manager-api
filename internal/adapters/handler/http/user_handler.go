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

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email    string         `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string         `json:"password" binding:"required,min=6" example:"securePassword123"`
	Profile  domain.Profile `json:"profile" binding:"required"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input"`
}

func NewUserHandler(userUC ports.UserUseCase) *UserHandler {
	return &UserHandler{
		userUC: userUC,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user account with email, password, and profile information
// @Description The password will be securely hashed before storage
// @Tags users
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration data"
// @Success 201 {object} RegisterResponse "User registered successfully"
// @Failure 400 {object} ErrorResponse "Bad request - invalid input data"
// @Failure 409 {object} ErrorResponse "Conflict - email already exists"
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.userUC.Register(c.Request.Context(), req.Email, req.Password, req.Profile); err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "already in use") {
			status = http.StatusConflict
		}
		c.JSON(status, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a specific user by their UUID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User UUID" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 200 {object} domain.User "User details"
// @Failure 400 {object} ErrorResponse "Bad request - invalid UUID format"
// @Failure 404 {object} ErrorResponse "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	user, err := h.userUC.GetUserByID(c.Request.Context(), idParam)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUsers godoc
// @Summary Get users with advanced filtering
// @Description Retrieve a paginated list of users with optional search, sorting, and field selection
// @Description Supports full-text search across email, first name, and last name
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number (1-based)" default(1) minimum(1)
// @Param page_size query int false "Number of users per page" default(10) minimum(1) maximum(100)
// @Param search query string false "Search term for email, first name, or last name" example("john")
// @Param sort query string false "Sort field" Enums(email, created_at, updated_at, first_name, last_name) example("created_at")
// @Param order query string false "Sort order" Enums(asc, desc) default(asc) example("desc")
// @Param fields query string false "Comma-separated list of fields to include in response" example("email,profile.first_name,created_at")
// @Success 200 {object} GetUsersResponse "List of users with pagination info"
// @Failure 400 {object} ErrorResponse "Bad request - invalid parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /users [get]
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
	if sortBy != "" && !validSortFields[sortBy] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid sort field. Valid options: email, created_at, updated_at, first_name, last_name",
		})
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

	result, err := h.userUC.GetUsers(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Description Remove a specific user by their UUID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User UUID" example("550e8400-e29b-41d4-a716-446655440000")
// @Success 200 {object} domain.User "User details"
// @Failure 400 {object} ErrorResponse "Bad request - invalid UUID format"
// @Failure 404 {object} ErrorResponse "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	err := h.userUC.DeleteUser(c.Request.Context(), idParam)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// func (h *UserHandler) UpdateUser(c *gin.Context) {
// 	idParam := c.Param("id")
// 	var user domain.User
// 	if err := c.ShouldBindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
// 		return
// 	}
// 	err := h.userUC.UpdateUser(c.Request.Context(), &user)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "not found") {
// 			c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
// 		}
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
// }
