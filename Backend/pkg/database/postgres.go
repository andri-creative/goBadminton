package database

import (
	"backend/internal/models"
	"backend/pkg/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(cfg *config.Config) *gorm.DB {
	dsn := cfg.DatabaseURL
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("✅ Connected to database successfully")

	// Auto migrate tables
	err = autoMigrate(db)
	if err != nil {
		log.Fatal("Failed to auto migrate:", err)
	}

	return db
}

func autoMigrate(db *gorm.DB) error {
	// Auto migrate all tables
	err := db.AutoMigrate(
		&models.User{},
		&models.Court{},
		&models.Reservation{},
		&models.Payment{},
	)
	if err != nil {
		return err
	}

	fmt.Println("✅ Database tables migrated successfully")
	return nil
}
