package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"chatingApp/models"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository is responsible for database operations related to users.
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository initializes a new UserRepository instance.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetAllUsers retrieves all users from the database.
func (repo *UserRepository) GetAllUsers() ([]models.User, error) {
	// Execute an SQL query to fetch all users
	rows, err := repo.DB.Query("SELECT id, name, email, role, created_at, updated_at FROM users")
	if err != nil {
		log.Println("Error: Failed to retrieve users", err)
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed when function exits

	var users []models.User

	// Iterate through the result set and populate the users slice
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// AddUser inserts a new user into the database.
func (repo *UserRepository) AddUser(name, email, password, role string) error {
	// SQL statement to insert a new user with hashed password
	query := "INSERT INTO users (name, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	// Get the current timestamp
	timestamp := time.Now()

	// Execute the SQL statement using Exec (not Query)
	_, err := repo.DB.Exec(query, name, email, password, role, timestamp, timestamp)
	if err != nil {
		log.Println("Error: Failed to insert user", err)
		return err
	}

	log.Println("✅ User added successfully:", name)
	return nil
}

// Login verifies user credentials and returns user info if valid.
func (repo *UserRepository) Login(email, password string) (*models.User, error) {
	// SQL query to find the user by email
	query := "SELECT id, name, email, password, role FROM users WHERE email = $1"

	// Execute query
	row := repo.DB.QueryRow(query, email)

	// Map result to a user struct
	var user models.User
	var hashedPassword string

	err := row.Scan(&user.ID, &user.Name, &user.Email, &hashedPassword, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("❌ Error: User not found")
			return nil, errors.New("user not found")
		}
		log.Println("❌ Error: Failed to find user", err)
		return nil, err
	}

	// Compare hashed password with provided password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		log.Println("❌ Error: Invalid password")
		return nil, errors.New("invalid credentials")
	}

	log.Println("✅ User logged in successfully:", user.Name)
	return &user, nil
}
