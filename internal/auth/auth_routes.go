package auth

import (
	"github.com/alanloffler/go-calth-api/internal/config"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, protected *gin.RouterGroup, q *sqlc.Queries, cfg *config.Config) *AuthHandler {
	var service *AuthService = NewAuthService(cfg)
	var repo *AuthRepository = NewAuthRepository(q)
	var handler *AuthHandler = NewAuthHandler(cfg, repo, service)

	public := router.Group("/auth")

	public.POST("/login", handler.Login)
	public.POST("/logout", handler.Logout)
	public.POST("/refresh", handler.Refresh)

	protected.GET("/auth/me", handler.GetMe)

	return handler
}
