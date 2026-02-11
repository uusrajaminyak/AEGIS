package main

import (
		"fmt"
		"log"
		"github.com/gin-gonic/gin"
		"github.com/uusrajaminyak/aegis-backend/config"
		"github.com/uusrajaminyak/aegis-backend/internal/adapter"
)

func main() {
		cfg, err := config.LoadConfig(".")
		if err != nil {
				log.Fatalf("Failed to load config: %v", err)
		}
		fmt.Println("Configuration loaded successfully")
		adapter.ConnectPostgres(cfg)
		fmt.Println("Database connected successfully")
		r := gin.Default()
		r.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
						"message": "pong",
						"status": "success",
				})
		})
		port := cfg.ServerPort
		if port == "" {
				port = "8080"
		}
		fmt.Printf("Starting server on port %s...\n", port)
		if err := r.Run(":" + port); err != nil {
				log.Fatalf("Failed to start server: %v", err)
		}
}