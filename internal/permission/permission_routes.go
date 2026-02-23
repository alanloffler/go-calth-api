package permission

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries) {
	var repo *PermissionRepository = NewPermissionRepository(q)
	var handler *PermissionHandler = NewPermissionHandler(repo)
	var permissions *gin.RouterGroup = router.Group("/permissions")

	permissions.POST("", middleware.PermissionMiddleware(q, "permissions-create"), handler.Create)
	permissions.GET("", middleware.PermissionMiddleware(q, "permissions-view"), handler.GetAll)
	permissions.GET("/soft", middleware.PermissionMiddleware(q, "permissions-view"), handler.GetAllWithSoftDeleted)
	permissions.GET("/category/:category", middleware.PermissionMiddleware(q, "permissions-view"), handler.GetAllByCategory)
	permissions.GET("/grouped", middleware.PermissionMiddleware(q, "permissions-view"), handler.GetAllGrouped)
	permissions.GET("/:id", middleware.PermissionMiddleware(q, "permissions-view"), handler.GetOneByID)
	permissions.GET("/:id/soft", middleware.PermissionMiddleware(q, "permissions-view"), handler.GetOneByIDWithSoftDeleted)
	permissions.PATCH("/:id", middleware.PermissionMiddleware(q, "permissions-update"), handler.Update)
	permissions.PATCH("/:id/restore", middleware.PermissionMiddleware(q, "permissions-restore"), handler.Restore)
	permissions.DELETE("/:id", middleware.PermissionMiddleware(q, "permissions-delete"), handler.Delete)
	permissions.DELETE("/:id/soft", middleware.PermissionMiddleware(q, "permissions-delete-hard"), handler.SoftDelete)
}
