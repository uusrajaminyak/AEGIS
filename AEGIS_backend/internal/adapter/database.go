package adapter

import (
		"fmt"
		"log"

		"github.com/uusrajaminyak/aegis-backend/config"
		"gorm.io/driver/postgres"
		"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres(cfg config.Config) {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
				cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

		var err error
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		log.Println("Database connection established")
}