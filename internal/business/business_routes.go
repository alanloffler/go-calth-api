package business

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/alanloffler/go-calth-api/internal/user"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries) {
	var repo *BusinessRepository = NewBusinessRepository(q)
	var userRepo *user.UserRepository = user.NewUserRepository(q)
	var handler *BusinessHandler = NewBusinessHandler(repo, userRepo)
	var businesses *gin.RouterGroup = router.Group("/businesses")

	businesses.POST("", middleware.PermissionMiddleware(q, "business-create"), handler.Create)

	businesses.GET("", middleware.PermissionMiddleware(q, "business-view"), handler.GetAll)
	businesses.GET("/:id", middleware.PermissionMiddleware(q, "business-view"), handler.GetOneByID)

	businesses.PATCH("/:id", middleware.PermissionMiddleware(q, "business-update"), handler.Update)

	businesses.DELETE("/:id", middleware.PermissionMiddleware(q, "business-delete-hard"), handler.Delete)
}
