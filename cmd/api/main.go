package main

import (
	"log"

	"github.com/alanloffler/go-calth-api/internal/business"
	"github.com/alanloffler/go-calth-api/internal/config"
	"github.com/alanloffler/go-calth-api/internal/database"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/health"
	"github.com/alanloffler/go-calth-api/internal/permission"
	"github.com/alanloffler/go-calth-api/internal/role"
	"github.com/alanloffler/go-calth-api/internal/user"
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

	// sqlc queries
	var queries *sqlc.Queries = sqlc.New(pool)

	// Gin router
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	health.RegisterRoutes(router, pool)
	user.RegisterRoutes(router, queries)
	business.RegisterRoutes(router, queries)
	role.RegisterRoutes(router, queries)
	permission.RegisterRoutes(router, queries)

	router.Run(":" + cfg.Port)
}
