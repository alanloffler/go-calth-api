package auth

import (
	"github.com/alanloffler/go-calth-api/internal/config"
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, q *sqlc.Queries, cfg *config.Config) {
	var service *AuthService = NewAuthService(cfg)
	var repo *AuthRepository = NewAuthRepository(q)
	var handler *AuthHandler = NewAuthHandler(service, repo)
	var auth *gin.RouterGroup = router.Group("/auth")

	auth.POST("/login", handler.Login)
	auth.POST("/logout", handler.Logout)
	auth.POST("/refresh", handler.Refresh)
}
