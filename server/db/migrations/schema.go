package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateTables(db *sql.DB) {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Error: error creating table %v", err)
		}
	}

	fmt.Println("✅ Success: all of the tables was created successfully")
}
