package repository

import (
	"database/sql"
	"log"

	"chatingApp/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) GetAllUsers() ([]models.User, error) {
	rows, err := repo.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Println("❌ שגיאה בשליפת משתמשים:", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
