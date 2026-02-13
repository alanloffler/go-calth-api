package main

import (
	"log"

	"github.com/alanloffler/go-calth-api/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	var router *gin.Engine = gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Calth API running", "status": "success"})
	})

	router.Run(":" + cfg.Port)
}
