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

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page     int
	PageSize int
}

// FilterParams holds filtering and sorting parameters
type FilterParams struct {
	Fields []string
	Search string
	SortBy string
	Order  string
}

// Constants for validation
const (
	DefaultPage      = 1
	DefaultPageSize  = 10
	MaxPageSize      = 100
	DefaultSortField = "created_at"
	DefaultSortOrder = "asc"
)

// Valid sort fields to prevent injection attacks
var validSortFields = map[string]bool{
	"email":      true,
	"created_at": true,
	"updated_at": true,
	"first_name": true,
	"last_name":  true,
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
	// Parse pagination parameters
	pagination := h.parsePaginationParams(c)

	// Parse and validate filter parameters
	filter, errResp := h.parseFilterParams(c)
	if errResp != nil {
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	// Build filter options for use case
	options := h.buildUserFilterOptions(pagination, filter)

	// Execute use case
	result, err := h.userUC.GetUsers(c.Request.Context(), options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *UserHandler) parsePaginationParams(c *gin.Context) PaginationParams {
	params := PaginationParams{
		Page:     DefaultPage,
		PageSize: DefaultPageSize,
	}

	// Parse page parameter
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	// Parse page_size parameter
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= MaxPageSize {
			params.PageSize = pageSize
		}
	}

	return params
}

// parseFilterParams extracts and validates filter parameters from the request
func (h *UserHandler) parseFilterParams(c *gin.Context) (*FilterParams, *ErrorResponse) {
	params := &FilterParams{
		SortBy: DefaultSortField,
		Order:  DefaultSortOrder,
	}

	// Parse field selection
	if fieldsParam := c.Query("fields"); fieldsParam != "" {
		fields := strings.Split(fieldsParam, ",")
		params.Fields = make([]string, len(fields))
		for i, field := range fields {
			params.Fields[i] = strings.TrimSpace(field)
		}
	}

	// Parse search parameter
	params.Search = strings.TrimSpace(c.Query("search"))

	// Parse and validate sort field
	if sortBy := strings.TrimSpace(c.Query("sort")); sortBy != "" {
		if !validSortFields[sortBy] {
			return nil, &ErrorResponse{
				Error: "Invalid sort field. Valid options: email, created_at, updated_at, first_name, last_name",
			}
		}
		params.SortBy = sortBy
	}

	// Parse and validate sort order
	if order := strings.ToLower(strings.TrimSpace(c.Query("order"))); order != "" {
		if order != "asc" && order != "desc" {
			params.Order = DefaultSortOrder
		} else {
			params.Order = order
		}
	}

	return params, nil
}

// buildUserFilterOptions creates the filter options for the use case
func (h *UserHandler) buildUserFilterOptions(pagination PaginationParams, filter *FilterParams) *ports.GetUsersOptions {
	return &ports.GetUsersOptions{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
		Fields:   filter.Fields,
		Search:   filter.Search,
		SortBy:   filter.SortBy,
		Order:    filter.Order,
	}
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
