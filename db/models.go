package db

import (
	"bookstore/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect database: %v", err)
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&Book{}, &Author{}); err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return db, nil
}
