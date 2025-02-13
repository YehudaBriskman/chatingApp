package main

import (
	"fmt"
	"os"

	"chatingApp/db"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("📌 Usage: go run scripts/run_queries.go <command>")
		fmt.Println("📌 Available commands: drop, truncate, reset, custom")
		return
	}

	switch os.Args[1] {
	case "drop":
		fmt.Println("🚨 Dropping all tables...")
		db.DropTables()
		fmt.Println("✅ Success: All tables have been dropped.")

	case "truncate":
		fmt.Println("🚨 Truncating all tables (resetting data)...")
		db.TruncateTables()
		fmt.Println("✅ Success: All tables have been truncated.")

	case "reset":
		fmt.Println("🚨 Resetting the database (drop & recreate tables)...")
		db.ResetDatabase()
		fmt.Println("✅ Success: Database has been reset.")

	case "custom":
		if len(os.Args) < 3 {
			fmt.Println("❌ Error: You must provide an SQL query to execute.")
			return
		}
		query := os.Args[2]
		fmt.Println("🚀 Executing custom SQL query...")
		db.ExecuteSQL(query)
		fmt.Println("✅ Success: Custom SQL query executed.")

	default:
		fmt.Println("❌ Error: Invalid command. Use one of: drop, truncate, reset, custom")
	}
}
