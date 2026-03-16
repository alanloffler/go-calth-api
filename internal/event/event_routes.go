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

	events.GET("/professional/:id", middleware.PermissionMiddleware(q, "events-view"), handler.GetByProfessionalID)
	events.GET("/business", middleware.PermissionMiddleware(q, "events-view"), handler.GetByBusinessID)
	events.GET("/filtered", middleware.PermissionMiddleware(q, "events-view"), handler.GetFiltered)
	events.GET("/days-with-events/:id", middleware.PermissionMiddleware(q, "events-view"), handler.GetDaysWithEvents)
	// events.GET("/professional/:id/date-array/:day", middleware.PermissionMiddleware(q, "events-view"), handler.GetProfessionalEventsByDay)
	events.GET("/patient/:patient_id", middleware.PermissionMiddleware(q, "events-view"), handler.GetByBusinessProfessionalPatient)
	events.GET("/professional/:id/date-array/:day", middleware.PermissionMiddleware(q, "events-view"), handler.GetProfessionalEventsByDayArray)
	events.GET("/:id", middleware.PermissionMiddleware(q, "events-view"), handler.GetByID)

	events.PATCH("/:id/status", middleware.PermissionMiddleware(q, "events-update"), handler.UpdateStatus)
	events.PATCH("/:id", middleware.PermissionMiddleware(q, "events-update"), handler.Update)
}
