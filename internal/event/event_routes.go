package event

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries) {
	var repo *EventRepository = NewEventRepository(q)
	var handler *EventHandler = NewEventHandler(repo)
	var events *gin.RouterGroup = router.Group("/events")

	events.POST("", middleware.PermissionMiddleware(q, "events-create"), handler.Create)
}
