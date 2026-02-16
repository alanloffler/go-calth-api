package permission

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, q *sqlc.Queries) {
	var repo *PermissionRepository = NewPermissionRepository(q)
	var handler *PermissionHandler = NewPermissionHandler(repo)
	var permissions *gin.RouterGroup = router.Group("/permissions")

	permissions.POST("", handler.Create)
	permissions.GET("", handler.GetAll)
	permissions.GET("/soft", handler.GetAllWithSoftDeleted)
	permissions.GET("/:id", handler.GetOneByID)
	permissions.PATCH("/:id", handler.Update)
	permissions.PATCH("/:id/restore", handler.Restore)
	permissions.DELETE("/:id", handler.Delete)
	permissions.DELETE("/:id/soft", handler.SoftDelete)
}
