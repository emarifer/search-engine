package db

import (
	"log"
	"os"

	"github.com/emarifer/search-engine/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func InitDB() {
	dbUrl := os.Getenv("DATABASE_URL")

	var err error

	dbConn, err = gorm.Open(postgres.Open(dbUrl))
	if err != nil {
		log.Fatalf("🔥 failed to connect to the database: %s\n", err)
	}

	// Enable uuid-ossp extension
	err = dbConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Fatalln("🔥 failed to enable uuid-ossp extension")
	}

	// Make migrations
	err = dbConn.AutoMigrate(
		&services.User{},
		&services.SearchSettings{},
		&services.CrawledUrl{},
		&services.SearchIndex{},
	)
	if err != nil {
		log.Fatalf("🔥 failed to migrate: %s\n", err)
	}

	log.Println("🚀 connected successfully to the database")
}

func GetDB() *gorm.DB {
	return dbConn
}
