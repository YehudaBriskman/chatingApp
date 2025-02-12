package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/YehudaBriskman/chatingApp/server/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

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
		log.Fatal("❌ שגיאה בחיבור למסד הנתונים:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ מסד הנתונים אינו מגיב:", err)
	}

	fmt.Println("✅ חיבור למסד הנתונים הצליח!")
}
