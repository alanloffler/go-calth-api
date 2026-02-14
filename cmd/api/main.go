package main

import (
	"log"

	"github.com/alanloffler/go-calth-api/internal/config"
	"github.com/alanloffler/go-calth-api/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load config
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Database connection
	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer pool.Close()

	// Gin router
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Calth API running", "status": "success"})
	})
	router.GET("/health/db", func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(500, gin.H{"message": err.Error(), "status": "error"})
			return
		}
		c.JSON(200, gin.H{
			"message": "Database connected successfully", "status": "connected",
		})
	})

	router.Run(":" + cfg.Port)
}
