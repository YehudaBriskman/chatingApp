package services

import (
	"chatingApp/models"
	"chatingApp/repository"
	"errors"
	"os"
	"strconv"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = os.Getenv("SECRET_KEY")

// UserService provides business logic for user operations.
type UserService struct {
	UserRepo *repository.UserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}

// GetAllUsers retrieves all users from the repository.
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.UserRepo.GetAllUsers()
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.UserRepo.GetUserByEmail(email)
}

// AddUser hashes the password and adds a new user to the repository.
func (s *UserService) AddUser(name, email, password, role string) error {
	// Validate role
	if role != "admin" && role != "super-admin" && role != "user" {
		return errors.New("invalid role: must be 'super-admin' or 'admin' or 'user'")
	}

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	return s.UserRepo.AddUser(name, email, string(hashedPassword), role)
}

// Login authenticates a user and returns a JWT token.
func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.Login(email, password)
	if err != nil {
		return "", errors.New(err.Error())
	}

	// Generate JWT token
	token, err := s.GenerateToken(user.Email, user.ID, user.Role)
	if err != nil {
		return "", errors.New("failed to generate authentication token")
	}

	return token, nil
}

// GenerateToken generates a JWT token for user authentication
func (s *UserService) GenerateToken(email string, id interface{}, role string) (string, error) {
	var userID int

	switch v := id.(type) {
	case string:
		parsedID, err := strconv.Atoi(v)
		if err != nil {
			return "", errors.New("invalid user ID format")
		}
		userID = parsedID
	case int:
		userID = v
	default:
		return "", errors.New("unsupported user ID type")
	}

	claims := jwt.MapClaims{
		"email":   email,
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return signedToken, nil
}

// ValidateToken verifies the JWT token and extracts claims
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure the token is signed with HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Extract claims and ensure the token is valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Convert user_id to int if needed
	if userID, exists := claims["user_id"]; exists {
		switch v := userID.(type) {
		case string:
			parsedID, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.New("invalid user_id format")
			}
			claims["user_id"] = parsedID
		case float64:
			// JSON unmarshal converts numbers to float64, so we cast to int
			claims["user_id"] = int(v)
		}
	}

	return claims, nil
}
