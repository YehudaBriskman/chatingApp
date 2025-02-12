package db

import (
	"database/sql"
	"fmt"
	"log"

	"chatingApp/config"
	"chatingApp/db/migrations"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// ConnectDB initializes the database connection and applies migrations
func ConnectDB() {
	config.LoadConfig()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.AppConfig.DBHost, config.AppConfig.DBPort, config.AppConfig.DBUser,
		config.AppConfig.DBPassword, config.AppConfig.DBName, config.AppConfig.DBSSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Error: Failed to connect to the database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Error: Database is not responding:", err)
	}

	fmt.Println("✅ Success: Connected to the database.")

	// Execute database migrations to ensure tables exist
	migrations.CreateTables(DB)
	fmt.Println("✅ Success: Database migrations applied.")
}
