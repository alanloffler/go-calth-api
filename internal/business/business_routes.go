package business

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/alanloffler/go-calth-api/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(public *gin.Engine, protected *gin.RouterGroup, q *sqlc.Queries, pool *pgxpool.Pool, queueClient *asynq.Client, appDomain string) {
	var repo *BusinessRepository = NewBusinessRepository(q)
	var userRepo *user.UserRepository = user.NewUserRepository(q)
	var handler *BusinessHandler = NewBusinessHandler(repo, userRepo, pool, queueClient, appDomain)

	public.POST("/businesses", handler.Create)
	public.GET("/businesses/availability/tax-id/:taxId", handler.CheckTaxIDAvailability)
	public.GET("/businesses/availability/slug/:slug", handler.CheckSlugAvailability)

	var businesses *gin.RouterGroup = protected.Group("/businesses")

	businesses.GET("", middleware.PermissionMiddleware(q, "business-view"), handler.GetAll)
	businesses.GET("/:id", middleware.PermissionMiddleware(q, "business-view"), handler.GetOneByID)

	businesses.PATCH("/:id", middleware.PermissionMiddleware(q, "business-update"), handler.Update)

	businesses.DELETE("/:id", middleware.PermissionMiddleware(q, "business-delete-hard"), handler.Delete)
}
