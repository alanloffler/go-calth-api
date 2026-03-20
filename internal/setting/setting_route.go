package setting

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries) {
	var repo *SettingRepository = NewSettingRepository(q)
	var handler *SettingHandler = NewSettingHandler(repo)
	var settings *gin.RouterGroup = router.Group("/settings")

	settings.PATCH("/:id", middleware.PermissionMiddleware(q, "settings-update"), handler.Update)

	settings.GET("", middleware.PermissionMiddleware(q, "settings-view"), handler.GetAll)
	settings.GET("/by-module/:module", middleware.PermissionMiddleware(q, "settings-view"), handler.GetByModule)
}
