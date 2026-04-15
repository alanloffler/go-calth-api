package event

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/alanloffler/go-calth-api/internal/professional_profile"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries, pool *pgxpool.Pool, queueClient *asynq.Client) {
	var repo *EventRepository = NewEventRepository(q)
	var profileRepo *professional_profile.ProfessionalProfileRepository = professional_profile.NewProfessionalProfileRepository(q)
	var handler *EventHandler = NewEventHandler(repo, pool, profileRepo, queueClient)
	var events *gin.RouterGroup = router.Group("/events")

	events.POST("", middleware.PermissionMiddleware(q, "events-create"), handler.Create)

	events.GET("/check-recurring", middleware.PermissionMiddleware(q, "events-view"), handler.CheckRecurring)
	events.GET("/business", middleware.PermissionMiddleware(q, "events-view"), handler.GetByBusinessID)
	events.GET("/filtered", middleware.PermissionMiddleware(q, "events-view"), handler.GetFiltered)
	events.GET("/professional/:id", middleware.PermissionMiddleware(q, "events-view"), handler.GetByProfessionalID)
	events.GET("/days-with-events/:id", middleware.PermissionMiddleware(q, "events-view"), handler.GetDaysWithEvents)
	// events.GET("/professional/:id/date-array/:day", middleware.PermissionMiddleware(q, "events-view"), handler.GetByProfessionalDay)
	events.GET("/patient/:patient_id", middleware.PermissionMiddleware(q, "events-view"), handler.GetByBusinessProfessionalPatient)
	events.GET("/professional/:id/date-array/:day", middleware.PermissionMiddleware(q, "events-view"), handler.GetByProfessionalDayArray)
	events.GET("/:id", middleware.PermissionMiddleware(q, "events-view"), handler.GetByID)

	events.PATCH("/:id/status", middleware.PermissionMiddleware(q, "events-update"), handler.UpdateStatus)
	events.PATCH("/:id", middleware.PermissionMiddleware(q, "events-update"), handler.Update)

	events.DELETE("/:id", middleware.PermissionMiddleware(q, "events-delete-hard"), handler.Delete)
}
