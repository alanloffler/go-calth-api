package blocked_day

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries) {
	var repo *BlockedDayRepository = NewBlockedDayRepository(q)
	var handler *BlockedDayHandler = NewBlockedDayHandler(repo)
	var blocked_days *gin.RouterGroup = router.Group("blocked-days")

	blocked_days.POST("", middleware.PermissionMiddleware(q, "patient-update"), handler.Create)
}
