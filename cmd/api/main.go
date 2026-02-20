package main

import (
	"log"
	"strings"

	"github.com/alanloffler/go-calth-api/internal/auth"
	"github.com/alanloffler/go-calth-api/internal/business"
	"github.com/alanloffler/go-calth-api/internal/config"
	"github.com/alanloffler/go-calth-api/internal/database"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/health"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/alanloffler/go-calth-api/internal/permission"
	"github.com/alanloffler/go-calth-api/internal/role"
	"github.com/alanloffler/go-calth-api/internal/user"
	"github.com/gin-contrib/cors"
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
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return strings.HasSuffix(origin, cfg.CorsOrigin)
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	var authService *auth.AuthService = auth.NewAuthService(cfg)

	// Public routes
	health.RegisterRoutes(router, pool)
	auth.RegisterRoutes(router, queries, cfg)
	user.RegisterRoutes(router, queries)
	role.RegisterRoutes(router, queries, pool)
	permission.RegisterRoutes(router, queries)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))

	business.RegisterRoutes(protected, queries)

	router.Run(":" + cfg.Port)
}
