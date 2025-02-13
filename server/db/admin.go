package db

import (
	"database/sql"
	"fmt"
	"log"

	"chatingApp/config"

	_ "github.com/lib/pq"
)

// ExecuteSQL executes an SQL query that is sent to it
func ExecuteSQL(query string) {
	config.LoadConfig()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.AppConfig.DBHost, config.AppConfig.DBPort, config.AppConfig.DBUser,
		config.AppConfig.DBPassword, config.AppConfig.DBName, config.AppConfig.DBSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ ERROR: Failed to connect to the database:", err)
	}
	defer db.Close()

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("❌ ERROR: Failed to execute the SQL query:", err)
	}

	fmt.Println("✅ SUCCESS: SQL query executed successfully!")
}

// DropTables drops all tables from the database
func DropTables() {
	fmt.Println("🚨 Dropping all tables...")
	ExecuteSQL(`
		DROP TABLE IF EXISTS messages, users CASCADE;
	`)
	fmt.Println("✅ SUCCESS: All tables have been dropped.")
}

// TruncateTables removes all data from tables but keeps their structure
func TruncateTables() {
	fmt.Println("🚨 Truncating all tables...")
	ExecuteSQL(`
		TRUNCATE TABLE messages, users RESTART IDENTITY CASCADE;
	`)
	fmt.Println("✅ SUCCESS: All tables have been truncated.")
}

// ResetDatabase drops and recreates all tables
func ResetDatabase() {
	fmt.Println("🚨 Resetting the database...")
	DropTables()
	fmt.Println("✅ SUCCESS: Tables dropped.")

	ConnectDB() // Re-run migrations
	fmt.Println("✅ SUCCESS: Database has been reset and reinitialized.")
}