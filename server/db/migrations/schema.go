package migrations

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

func CreateTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT CHECK (role IN ('user', 'admin', 'super-admin')) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS rooms (
	        id SERIAL PRIMARY KEY,
	        name TEXT NOT NULL,
	        description TEXT,
	        created_by INT REFERENCES users(id) ON DELETE CASCADE,
	        room_admins INT[] DEFAULT '{}', -- List of user IDs as admins
	        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`,

		`CREATE TABLE IF NOT EXISTS room_users (
			room_id INT REFERENCES rooms(id) ON DELETE CASCADE,
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			PRIMARY KEY (room_id, user_id)
		);`,

		`CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			room_id INT REFERENCES rooms(id) ON DELETE CASCADE,
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS system_logs (
			id SERIAL PRIMARY KEY,
			method TEXT NOT NULL,
			endpoint TEXT NOT NULL,
			user_id INT NULL,
			status_code INT NOT NULL,
			message TEXT NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Error: failed to create table: %v", err)
		}
	}

	fmt.Println("✅ Success: all tables have been created successfully")

	// Insert default super-admin user if not exists
	email := os.Getenv("SUPERADMIN_EMAIL")
	password := os.Getenv("SUPERADMIN_PASSWORD")
	if email == "" || password == "" {
		log.Println("⚠️ Warning: SUPERADMIN_EMAIL or SUPERADMIN_PASSWORD is not set in environment variables.")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error: failed to hash super-admin password: %v", err)
	}

	_, err = db.Exec(
		`INSERT INTO users (name, email, password, role, created_at, updated_at)
		 VALUES ($1, $2, $3, 'super-admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		 ON CONFLICT (email) DO NOTHING;`,
		"YehudaSuper", email, string(hashedPassword),
	)
	if err != nil {
		log.Fatalf("Error: failed to insert super-admin user: %v", err)
	}

	// Fetch the created super-admin user
	var id int
	var name, fetchedEmail, role string
	var createdAt, updatedAt string
	row := db.QueryRow("SELECT id, name, email, role, created_at, updated_at FROM users WHERE email = $1", email)
	if err := row.Scan(&id, &name, &fetchedEmail, &role, &createdAt, &updatedAt); err != nil {
		log.Fatalf("Error: failed to fetch super-admin user: %v", err)
	}

	fmt.Printf("✅ Success: Super Admin user found - ID: %d, Name: %s, Email: %s, Role: %s, Created At: %s, Updated At: %s\n", id, name, fetchedEmail, role, createdAt, updatedAt)
}
