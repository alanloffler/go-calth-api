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
}
