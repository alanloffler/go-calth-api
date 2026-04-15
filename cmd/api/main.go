package main

import (
	"log"
	"strings"

	"github.com/alanloffler/go-calth-api/internal/auth"
	"github.com/alanloffler/go-calth-api/internal/business"
	"github.com/alanloffler/go-calth-api/internal/config"
	"github.com/alanloffler/go-calth-api/internal/database"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/event"
	"github.com/alanloffler/go-calth-api/internal/health"
	"github.com/alanloffler/go-calth-api/internal/medical_history"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/alanloffler/go-calth-api/internal/permission"
	"github.com/alanloffler/go-calth-api/internal/role"
	"github.com/alanloffler/go-calth-api/internal/setting"
	"github.com/alanloffler/go-calth-api/internal/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
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

	// Redis
	redisArr := cfg.RedisAddr
	if redisArr == "" {
		redisArr = "127.0.0.1:6379"
	}
	redisClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisArr})
	defer redisClient.Close()

	// sqlc queries
	var queries *sqlc.Queries = sqlc.New(pool)

	// Gin router
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if rootSuffix, ok := strings.CutPrefix(cfg.CorsOrigin, "."); ok {
				return strings.HasSuffix(origin, cfg.CorsOrigin) || strings.HasSuffix(origin, "://"+rootSuffix)
			}
			return strings.HasSuffix(origin, cfg.CorsOrigin)
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	var authService *auth.AuthService = auth.NewAuthService(cfg)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	event.RegisterRoutes(protected, queries, pool, redisClient)
	medical_history.RegisterRoutes(protected, queries)
	permission.RegisterRoutes(protected, queries)
	role.RegisterRoutes(protected, queries, pool)
	setting.RegisterRoutes(protected, queries)
	user.RegisterRoutes(protected, queries, pool)

	// Mixed routes (public/protected)
	auth.RegisterRoutes(router, protected, queries, cfg)
	business.RegisterRoutes(router, protected, queries, pool, redisClient, cfg.AppDomain)

	// Public routes
	health.RegisterRoutes(router, pool)

	router.Run(":" + cfg.Port)

}
