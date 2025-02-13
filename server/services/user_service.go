package services

import (
	"chatingApp/models"
	"chatingApp/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const secretKey = "your-secret-key"

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
	token, err := s.GenerateToken(user.Email, user.Role)
	if err != nil {
		return "", errors.New("failed to generate authentication token")
	}

	return token, nil
}

// GenerateToken generates a JWT token for user authentication
func (s *UserService) GenerateToken(email, role string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
