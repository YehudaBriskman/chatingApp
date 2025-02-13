package handlers

import (
	"net/http"
	"strings"
	"chatingApp/services"
	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	UserService *services.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{UserService: service}
}

// GetUsers handles the GET request to retrieve all users.
func (h *UserHandler) GetUsers(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
		return
	}

	claims, err := services.ValidateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	if claims["role"] != "admin" && claims["role"] != "super-admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	users, err := h.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// AddUser handles the POST request to add a new user.
func (h *UserHandler) AddUser(c *gin.Context) {
	// Struct to bind the incoming JSON request
	var userInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	// Bind JSON request to userInput struct
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// If the role is "admin" or "super-admin", require a token
	if userInput.Role == "admin" || userInput.Role == "super-admin" {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
			return
		}

		claims, err := services.ValidateToken(strings.TrimPrefix(token, "Bearer "))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if userInput.Role == "admin" && claims["role"] != "admin" && claims["role"] != "super-admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		if userInput.Role == "super-admin" && claims["role"] != "super-admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// Call the service to add a new user
	err := h.UserService.AddUser(userInput.Name, userInput.Email, userInput.Password, userInput.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}

// Login handles the POST request for user authentication.
func (h *UserHandler) Login(c *gin.Context) {
	// Struct to bind the incoming JSON request
	var loginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind JSON request to loginInput struct
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service to authenticate user
	token, err := h.UserService.Login(loginInput.Email, loginInput.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
