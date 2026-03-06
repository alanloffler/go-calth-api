package medical_history

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/alanloffler/go-calth-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, q *sqlc.Queries) {
	var repo *MedicalHistoryRepository = NewMedicalHistoryRepository(q)
	var handler *MedicalHistoryHandler = NewMedicalHistoryHandler(repo)
	var medical_histories *gin.RouterGroup = router.Group("/medical-history")

	medical_histories.POST("", middleware.PermissionMiddleware(q, "medical_history-create"), handler.CreateMedicalHistory)
	medical_histories.GET("/:id/patient/removed", middleware.PermissionMiddleware(q, "medical_history-view"), handler.GetAllByPatientIDWithSoftDeleted)
	medical_histories.GET("/:id/patient", middleware.PermissionMiddleware(q, "medical_history-view"), handler.GetAllByPatientIDWithSoftDeleted)
	medical_histories.PATCH("/:id/restore", middleware.PermissionMiddleware(q, "medical_history-restore"), handler.Restore)
	medical_histories.DELETE("/:id/soft", middleware.PermissionMiddleware(q, "medical_history-delete"), handler.SoftDelete)
}
